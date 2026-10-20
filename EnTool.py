#!/usr/bin/env python3
"""EnTool: offline, portable encrypted text files for short secrets.

Build-time requirement: cryptography (pip install cryptography).
Each encrypted .txt file works with EnTool on Windows, macOS, and Linux.
The application never accesses the internet.
"""

import base64
import binascii
import os
import re
import secrets
import string
import sys
from getpass import getpass
from pathlib import Path

from cryptography.exceptions import InvalidTag
from cryptography.hazmat.primitives.ciphers.aead import AESGCM
from cryptography.hazmat.primitives.kdf.scrypt import Scrypt


# Use the folder containing EnTool.exe (or the source script).
APP_DIR = Path(
    sys.executable if getattr(sys, "frozen", False) else __file__
).resolve().parent

# Keep ET1 so files created by older EnTool versions can still be decrypted.
FORMAT_VERSION = "ET1"
FILLER_ALPHABET = string.ascii_letters + string.digits + "-_"
MAX_CHARACTERS = 1_000_000
WINDOWS_RESERVED = {
    "CON", "PRN", "AUX", "NUL",
    *(f"COM{n}" for n in range(1, 10)),
    *(f"LPT{n}" for n in range(1, 10)),
}


def file_path_from_name(user_input: str) -> Path:
    """Accept a portable file name, never a path outside EnTool's folder."""
    name = user_input.strip()
    if not name:
        raise ValueError("Enter a file name.")
    if "/" in name or "\\" in name or ":" in name:
        raise ValueError("Enter a file name only, not a folder or path.")
    if not name.lower().endswith(".txt"):
        if Path(name).suffix:
            raise ValueError("Only .txt files are supported.")
        name += ".txt"
    stem = name[:-4]
    if (
        not stem
        or len(name) > 100
        or not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9 ._-]*", stem)
        or stem.endswith((" ", "."))
    ):
        raise ValueError(
            "Use 1-100 characters: English letters, numbers, spaces, ., - or _."
        )
    if stem.split(".")[0].upper() in WINDOWS_RESERVED:
        raise ValueError("This name is reserved by Windows. Choose another name.")
    return APP_DIR / name


def derive_key(secret_key: str, salt: bytes) -> bytes:
    """Derive a 256-bit encryption key from the user's passphrase."""
    return Scrypt(salt=salt, length=32, n=2**15, r=8, p=1).derive(
        secret_key.encode("utf-8")
    )


def minimum_characters(password: str) -> int:
    """Return the minimum exact-length ET1 file for this UTF-8 secret."""
    required_bytes = 16 + 12 + 16 + 4 + len(password.encode("utf-8"))
    return len(FORMAT_VERSION) + 4 * ((required_bytes + 2) // 3)


def encrypt_password(password: str, secret_key: str, character_count: int) -> str:
    """Encrypt the secret into exactly character_count printable characters."""
    minimum = minimum_characters(password)
    if character_count < minimum:
        raise ValueError(
            f"The secret is too long for this file size. Choose at least {minimum} characters."
        )
    if character_count > MAX_CHARACTERS:
        raise ValueError(f"The maximum is {MAX_CHARACTERS:,} characters.")

    # The encrypted data contains the UTF-8 byte length and random padding.
    # The final 0-3 text characters are also random, for an exact file length.
    encoded_bytes = 3 * ((character_count - len(FORMAT_VERSION)) // 4)
    password_data = password.encode("utf-8")
    padding_size = encoded_bytes - 16 - 12 - 16 - 4 - len(password_data)
    plaintext = (
        len(password_data).to_bytes(4, "big")
        + password_data
        + secrets.token_bytes(padding_size)
    )
    salt = secrets.token_bytes(16)
    nonce = secrets.token_bytes(12)
    ciphertext = AESGCM(derive_key(secret_key, salt)).encrypt(
        nonce, plaintext, FORMAT_VERSION.encode("ascii")
    )
    result = FORMAT_VERSION + base64.urlsafe_b64encode(
        salt + nonce + ciphertext
    ).decode("ascii")
    result += "".join(
        secrets.choice(FILLER_ALPHABET)
        for _ in range(character_count - len(result))
    )
    return result


def decrypt_password(encrypted_text: str, secret_key: str) -> str:
    """Authenticate and decrypt an ET1 file (including older EnTool files)."""
    if not encrypted_text.startswith(FORMAT_VERSION):
        raise ValueError("Unknown EnTool file format.")
    encoded_length = 4 * ((len(encrypted_text) - len(FORMAT_VERSION)) // 4)
    encoded = encrypted_text[
        len(FORMAT_VERSION):len(FORMAT_VERSION) + encoded_length
    ]
    raw = base64.b64decode(encoded, altchars=b"-_", validate=True)
    if len(raw) < 48:
        raise ValueError("Encrypted file is too short or damaged.")

    salt, nonce, ciphertext = raw[:16], raw[16:28], raw[28:]
    plaintext = AESGCM(derive_key(secret_key, salt)).decrypt(
        nonce, ciphertext, FORMAT_VERSION.encode("ascii")
    )
    if len(plaintext) < 4:
        raise ValueError("Damaged encrypted contents.")
    secret_length = int.from_bytes(plaintext[:4], "big")
    if secret_length > len(plaintext) - 4:
        raise ValueError("Damaged encrypted contents.")
    return plaintext[4:4 + secret_length].decode("utf-8")


def encrypt_file() -> None:
    try:
        target = file_path_from_name(
            input("Name for the new text file (e.g. message.txt): ")
        )
    except ValueError as error:
        print(f"Error: {error}")
        return

    if target.exists() or target.is_symlink():
        print(f"File already exists: {target.name}. Choose another name.")
        return

    password = getpass("Secret to encrypt (hidden): ")
    if not password:
        print("The secret cannot be empty.")
        return

    secret_key = getpass("Create an encryption key (hidden): ")
    if not secret_key:
        print("The encryption key cannot be empty.")
        return
    if getpass("Confirm the encryption key: ") != secret_key:
        print("The encryption keys do not match. No file was created.")
        return

    try:
        count_text = input(
            f"How many characters should {target.name} contain? (e.g. 100): "
        ).strip()
        character_count = int(count_text)
        encrypted_text = encrypt_password(password, secret_key, character_count)
    except ValueError as error:
        print(f"Error: {error}")
        return

    try:
        # Exclusive creation: existing files are never overwritten.
        descriptor = os.open(str(target), os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
        with os.fdopen(descriptor, "w", encoding="ascii", newline="") as handle:
            handle.write(encrypted_text)
    except FileExistsError:
        print(f"File already exists: {target.name}. Nothing was overwritten.")
        return

    print(f"\nEncrypted file created: {target.name}")
    print(f"Location: {APP_DIR}")
    print(f"File length: {len(encrypted_text)} characters")
    print("Share the .txt file separately from the encryption key.")
    print("Keep the key safe: without it, the secret cannot be recovered.")


def decrypt_file() -> None:
    files = sorted(p.name for p in APP_DIR.glob("*.txt") if p.is_file())
    if files:
        print("Text files in the EnTool folder: " + ", ".join(files[:20]))
        if len(files) > 20:
            print(f"... and {len(files) - 20} more.")

    try:
        source = file_path_from_name(
            input("Name of the encrypted text file (e.g. message.txt): ")
        )
    except ValueError as error:
        print(f"Error: {error}")
        return

    if not source.is_file():
        print(f"File not found: {source.name}")
        print("The .txt file must be in the same folder as EnTool.")
        return
    if source.is_symlink() or source.resolve().parent != APP_DIR.resolve():
        print("For safety, only regular files inside the EnTool folder are supported.")
        return
    if source.stat().st_size > MAX_CHARACTERS + 2:
        print("This file is too large for EnTool.")
        return

    secret_key = getpass("Enter the decryption key (hidden): ")
    if not secret_key:
        print("The decryption key cannot be empty.")
        return
    try:
        encrypted_text = source.read_text(encoding="ascii").rstrip("\r\n")
        password = decrypt_password(encrypted_text, secret_key)
    except (InvalidTag, ValueError, UnicodeError, binascii.Error):
        print("Unable to decrypt: wrong key, unsupported format or damaged file.")
        return
    print(f"\nDecrypted secret: {password}")


def main() -> None:
    while True:
        print("\n========= EnTool =========")
        print("1. Encrypt a secret into a new text file")
        print("2. Decrypt a secret from a text file")
        print("0. Exit")
        choice = input("Select an option (1/2/0): ").strip()
        try:
            if choice == "1":
                encrypt_file()
            elif choice == "2":
                decrypt_file()
            elif choice == "0":
                print("Goodbye!")
                break
            else:
                print("Select 1, 2 or 0.")
        except OSError:
            print("File error. Check the folder and file permissions.")


if __name__ == "__main__":
    try:
        main()
    except (KeyboardInterrupt, EOFError):
        print("\nEnTool closed.")
