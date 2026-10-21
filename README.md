```text
 _____       _____           _
| ____|_ __ |_   _|__   ___ | |
|  _| | '_ \  | |/ _ \ / _ \| |
| |___| | | | | | (_) | (_) | |
|_____|_| |_| |_|\___/ \___/|_|
              by christofrancis
```

# EnTool — Native Go edition

**EnTool** is a free, open-source, entirely offline command-line tool for encrypting a short text secret into a `.txt` file and decrypting it with a passphrase.

This is a **complete rewrite in Go**, not Python or PyInstaller. It produces native executables without shipping a Python interpreter or unpacking one at startup.

## Downloads

| System | Download |
| --- | --- |
| Windows x64 | [EnToolWindows.exe](downloads/EnToolWindows.exe) |
| Linux x64 | [EnToolLinux.tar.gz](downloads/EnToolLinux.tar.gz) |
| macOS Intel x64 | [EnToolMacOS.tar.gz](downloads/EnToolMacOS.tar.gz) |
| macOS Apple Silicon arm64 | [EnToolMacOSAppleSilicon.tar.gz](downloads/EnToolMacOSAppleSilicon.tar.gz) |

Click a download above, then **Download raw file** on GitHub. Downloaded files are built from the Go source on GitHub Actions, with automatic startup tests. SHA-256 checksums are listed in [`downloads/SHA256SUMS`](downloads/SHA256SUMS).

**Windows:** put `EnToolWindows.exe` in a writable folder, open it, and use its terminal menu.

**Linux:** extract `EnToolLinux.tar.gz` to a writable folder, then `chmod +x EnToolLinux && ./EnToolLinux`.

**macOS:** extract the appropriate `.tar.gz` archive to a writable folder, then `chmod +x EnToolMacOS && ./EnToolMacOS` (or `EnToolMacOSAppleSilicon` on ARM64). The macOS builds are **not Apple-notarized** and may show a Gatekeeper warning.

> **Security notice:** These executables are not digitally signed by a trusted publisher and have not undergone an independent security audit. **Antivirus/SmartScreen may still issue warnings**; writing software in Go cannot guarantee zero alerts. Do not disable protection or add antivirus exclusions. Report suspected false positives to your antivirus vendor for review. A successful build or SHA-256 verification is not a malware clearance.

## How to use

1. Launch EnTool. The ASCII logo and `by christofrancis` appear above the menu.
2. Choose **1** to encrypt a secret into a new `.txt` file. Enter a file name, secret, encryption passphrase, confirmation, and desired length (for example 100 characters).
3. Send the encrypted `.txt` file separately from its passphrase.
4. The recipient runs EnTool and chooses **2** to decrypt the file with the same passphrase.
5. Choose **0** to exit.

The program reads and creates `.txt` files **next to the executable**. It will not overwrite existing files and rejects paths outside that folder. Secret and passphrase input is hidden when running in a real terminal. It does not need an internet connection.

## Compatibility and encryption

- **Compatible with existing `ET1` files** created by the original Python EnTool.
- AES-256-GCM with random 16-byte salt, 12-byte nonce and 16-byte authentication tag.
- scrypt key derivation (`N=32768`, `r=8`, `p=1`).
- UTF-8 secrets and an exact user-selected output length (within the 1,000,000-character limit).
- A longer encrypted file does **not** make a weak passphrase stronger.
- Keep your key safe: it cannot be recovered if lost. Decrypted text is displayed in the terminal; be mindful of terminal logging and screen sharing.

## Build from source

Requires Go 1.23 or later:

```sh
go mod download
go test ./...
go build -o EnTool .
```

Windows PowerShell: `go build -o EnTool.exe .` (build on Windows). The same source works on Linux and macOS. See [`.github/workflows/build.yml`](.github/workflows/build.yml) for tested platform release builds.

## License

[MIT License](LICENSE) — copyright (c) 2026 christofrancis.
