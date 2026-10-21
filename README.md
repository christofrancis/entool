```text
 _____       _____           _
| ____|_ __ |_   _|__   ___ | |
|  _| | '_ \  | |/ _ \ / _ \| |
| |___| | | | | | (_) | (_) | |
|_____|_| |_| |_|\___/ \___/|_|
              by christofrancis
```
```text
+====================================================================+
|                  EnTool - OFFLINE ENCRYPTION TOOL                  |
+====================================================================+
```
**EnTool** is a small, open-source command-line application for
encrypting private text into `.txt` files and decrypting them with
a passphrase. It runs entirely offline on Windows, Linux and macOS.

```text
+--------------------------------------------------------------------+
| IMPORTANT NOTICE                                                   |
+--------------------------------------------------------------------+
```
Encrypted files can be shared, but the passphrase should be sent
separately and kept private. A lost passphrase cannot be recovered.

Use a strong, unique passphrase. Keep the encrypted files and
decrypted output away from public terminals, logs and screen sharing.

```text
+--------------------------------------------------------------------+
| PROJECT FILES                                                      |
+--------------------------------------------------------------------+
```
**Application source:**

- [`main.go`](main.go) — startup logo, interactive menu and user input.
- [`crypto.go`](crypto.go) — encryption, decryption and `ET1` compatibility.
- [`files.go`](files.go) — safe file names and local `.txt` file operations.
- [`entool_test.go`](entool_test.go) — automated tests.

**Go configuration:**

- [`go.mod`](go.mod) — module definition and dependencies.
- [`go.sum`](go.sum) — dependency integrity checksums.
- [`.gitattributes`](.gitattributes) — consistent line endings.
- [`.gitignore`](.gitignore) — excludes local and temporary files.

**Release and documentation:**

- [`downloads/`](downloads/) — Windows, Linux and macOS programs.
- [`.github/workflows/build.yml`](.github/workflows/build.yml) — automated builds and tests.
- [`README.md`](README.md) — project information and instructions.
- [`LICENSE`](LICENSE) — MIT license.

```text
+--------------------------------------------------------------------+
| TECHNOLOGY                                                         |
+--------------------------------------------------------------------+
```
- **Go 1.23+** — native command-line programs without a Python runtime.
- **AES-256-GCM** — authenticated encryption (Go standard library).
- **scrypt** — passphrase-based key derivation (`golang.org/x/crypto`).
- **`golang.org/x/term`** — hidden secret and passphrase entry.
- **GitHub Actions** — cross-platform builds and automated checks.

```text
+--------------------------------------------------------------------+
| HOW THE PROJECT WORKS                                              |
+--------------------------------------------------------------------+
```
1. Start EnTool and select **1** to encrypt a secret.
2. Enter a new `.txt` file name, secret, passphrase and confirmation.
3. Choose the output file length (up to 1,000,000 characters).
4. EnTool creates an encrypted `.txt` file next to the executable.
5. Share the encrypted file separately from its passphrase.
6. Select **2** to open that file and decrypt it with the passphrase.
7. Select **0** to exit the program.

Existing files are never overwritten. Only local `.txt` file names
are accepted. Secrets and passphrases are hidden during interactive
entry; decrypted text is displayed on screen.

```text
+--------------------------------------------------------------------+
| DOWNLOADS                                                          |
+--------------------------------------------------------------------+
```
Choose your platform and open the file on GitHub. If its file page
appears, select **Download raw file**.

- **Windows x64:** [EnToolWindows.exe](downloads/EnToolWindows.exe)
- **Linux x64:** [EnToolLinux.tar.gz](downloads/EnToolLinux.tar.gz)
- **macOS Intel x64:** [EnToolMacOS.tar.gz](downloads/EnToolMacOS.tar.gz)
- **macOS Apple Silicon arm64:** [EnToolMacOSAppleSilicon.tar.gz](downloads/EnToolMacOSAppleSilicon.tar.gz)

File integrity hashes: [`downloads/SHA256SUMS`](downloads/SHA256SUMS).

```text
+--------------------------------------------------------------------+
| PLATFORM INSTRUCTIONS                                              |
+--------------------------------------------------------------------+
```
**Windows:** Put `EnToolWindows.exe` in a writable folder, then open it.

**Linux:** Extract `EnToolLinux.tar.gz` and run:

```sh
chmod +x EnToolLinux
./EnToolLinux
```

**macOS Intel:** Extract `EnToolMacOS.tar.gz` and run:

```sh
chmod +x EnToolMacOS
./EnToolMacOS
```

**macOS Apple Silicon:** Extract `EnToolMacOSAppleSilicon.tar.gz` and run:

```sh
chmod +x EnToolMacOSAppleSilicon
./EnToolMacOSAppleSilicon
```

All encrypted `.txt` files are saved next to the executable, so use
a folder where the program can create files.

```text
+--------------------------------------------------------------------+
| COMPATIBILITY AND SECURITY                                         |
+--------------------------------------------------------------------+
```
- Compatible with existing `ET1` files from the Python version.
- AES-256-GCM uses a random 16-byte salt and 12-byte nonce, with
  a 16-byte authentication tag.
- scrypt parameters: `N=32768`, `r=8`, `p=1`.
- Text is UTF-8; output length is user-selected.
- Increasing the encrypted file length does not strengthen a
  weak passphrase.

**Security notice:** EnTool has not undergone an independent security
audit. Release binaries are not signed by a trusted publisher;
antivirus software or SmartScreen may show warnings. Do not disable
your security protection or add exclusions to force a download.
SHA-256 checksums confirm integrity, not malware-free status.
macOS builds are not Apple-notarized and may show a Gatekeeper warning.

```text
+--------------------------------------------------------------------+
| BUILD FROM SOURCE                                                  |
+--------------------------------------------------------------------+
```
Requires **Go 1.23 or later**.

```sh
go mod download
go test ./...
go build -o EnTool .
```

On Windows, use `go build -o EnTool.exe .` in PowerShell.
The [build workflow](.github/workflows/build.yml) contains the
platform-specific release build commands.

```text
+--------------------------------------------------------------------+
| LICENSE                                                            |
+--------------------------------------------------------------------+
```
[MIT License](LICENSE) — copyright (c) 2026 christofrancis.
