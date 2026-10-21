```text
 _____       _____           _
| ____|_ __ |_   _|__   ___ | |
|  _| | '_ \  | |/ _ \ / _ \| |
| |___| | | | | | (_) | (_) | |
|_____|_| |_| |_|\___/ \___/|_|
              by christofrancis
```

# EnTool

**EnTool** is an open-source, offline command-line tool for encrypting short secrets into `.txt` files and decrypting them using a passphrase. It is written in **Go** and runs as a native executable without installing Go or Python.

## Downloads

| System | Download |
| --- | --- |
| Windows x64 | [EnToolWindows.exe](downloads/EnToolWindows.exe) |
| Linux x64 | [EnToolLinux.tar.gz](downloads/EnToolLinux.tar.gz) |
| macOS Intel x64 | [EnToolMacOS.tar.gz](downloads/EnToolMacOS.tar.gz) |
| macOS Apple Silicon arm64 | [EnToolMacOSAppleSilicon.tar.gz](downloads/EnToolMacOSAppleSilicon.tar.gz) |

Click your system's file in the table, then select **Download raw file** on GitHub. The programs are built from this repository by [GitHub Actions](.github/workflows/build.yml). SHA-256 file checksums are in [`downloads/SHA256SUMS`](downloads/SHA256SUMS).

**Windows:** place `EnToolWindows.exe` in a writable folder and open it. Follow the terminal menu.

**Linux:** extract `EnToolLinux.tar.gz`, then run `chmod +x EnToolLinux && ./EnToolLinux` in that folder.

**macOS Intel:** extract `EnToolMacOS.tar.gz`, then run `chmod +x EnToolMacOS && ./EnToolMacOS`.

**macOS Apple Silicon:** extract `EnToolMacOSAppleSilicon.tar.gz`, then run `chmod +x EnToolMacOSAppleSilicon && ./EnToolMacOSAppleSilicon`.

Keep each executable in a **writable folder**, because it creates and reads encrypted `.txt` files alongside itself. macOS binaries are not Apple-notarized and may trigger a Gatekeeper warning.

> **Security notice:** These executables are not digitally signed by a trusted publisher and have not undergone an independent security audit. **Antivirus/SmartScreen may still issue warnings**; writing software in Go cannot guarantee zero alerts. Do not disable protection or add antivirus exclusions. Report suspected false positives to your antivirus vendor for review. A successful build or SHA-256 verification is not a malware clearance.

## How to use

1. Launch EnTool.
2. Choose **1** to encrypt a secret into a new `.txt` file. Enter a file name, secret, encryption passphrase, confirmation, and desired length (for example 100 characters).
3. Send the encrypted `.txt` file separately from its passphrase.
4. The recipient runs EnTool and chooses **2** to decrypt the file with the same passphrase.
5. Choose **0** to exit.

The program reads and creates `.txt` files **next to the executable**, never overwrites existing files, and rejects paths outside that folder. Secrets and passphrases are entered without echoing them in an interactive terminal. Internet access is not required.

## Compatibility and encryption

- **Compatible with existing `ET1` files** created by the original Python EnTool.
- AES-256-GCM with random 16-byte salt, 12-byte nonce and 16-byte authentication tag.
- scrypt key derivation (`N=32768`, `r=8`, `p=1`).
- UTF-8 secrets and an exact user-selected output length (within the 1,000,000-character limit).
- A longer encrypted file does **not** make a weak passphrase stronger.
- Keep your key safe: it cannot be recovered if lost. Decrypted text is displayed in the terminal; be mindful of terminal logging and screen sharing.

## Technology

- **Go 1.23+** — native command-line executables for Windows, Linux, and macOS; no Python or PyInstaller.
- **AES-256-GCM** (Go standard library) for authenticated encryption, with **scrypt** (`golang.org/x/crypto`) to derive the encryption key from a passphrase.
- **`golang.org/x/term`** for hidden secret and passphrase input in the terminal.
- **GitHub Actions** for cross-platform builds and automated tests.

## Project files

| File or directory | Purpose |
| --- | --- |
| [`main.go`](main.go) | Application entry point, ASCII logo, menu and user input. |
| [`crypto.go`](crypto.go) | Encryption/decryption and compatibility with the `ET1` format. |
| [`files.go`](files.go) | Safe names, reading and writing encrypted `.txt` files. |
| [`entool_test.go`](entool_test.go) | Automated tests for encryption, file safety and the menu. |
| [`go.mod`](go.mod) | Go module and required dependencies. |
| [`go.sum`](go.sum) | Dependency verification checksums. |
| [`.github/workflows/build.yml`](.github/workflows/build.yml) | Automated builds and tests for the supported systems. |
| [`downloads/`](downloads/) | Ready-to-run programs and `SHA256SUMS` integrity checksums. |
| [`.gitignore`](.gitignore) | Keeps temporary files and local secrets out of Git. |
| [`.gitattributes`](.gitattributes) | Consistent source-file line endings across systems. |
| [`README.md`](README.md) | Project documentation and usage instructions. |
| [`LICENSE`](LICENSE) | MIT license. |

## Build from source

To build from source, install Go 1.23 or later:

```sh
go mod download
go test ./...
go build -o EnTool .
```

Windows PowerShell: `go build -o EnTool.exe .` (build on Windows). The same source works on Linux and macOS. See [`.github/workflows/build.yml`](.github/workflows/build.yml) for tested platform release builds.

## License

[MIT License](LICENSE) — copyright (c) 2026 christofrancis.
