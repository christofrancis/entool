```text
 _____       _____           _
| ____|_ __ |_   _|__   ___ | |
|  _| | '_ \  | |/ _ \ / _ \| |
| |___| | | | | | (_) | (_) | |
|_____|_| |_| |_|\___/ \___/|_|
              by christofrancis

+====================================================================+
|                  ENTOOL - OFFLINE ENCRYPTION TOOL                  |
+====================================================================+
```

**EnTool** is a free, open-source, offline command-line tool for encrypting short text secrets into `.txt` files and decrypting them with a shared passphrase.

Prebuilt packages are included **directly in this repository**. You do not need Python or an internet connection to run them. The Windows package is a ZIP folder bundle rather than a self-extracting single-file EXE.

```text
+--------------------------------------------------------------------+
| DOWNLOAD                                                           |
+--------------------------------------------------------------------+
```

| Operating system | Download |
|---|---|
| Windows (x64) | [EnToolWindows.zip](downloads/EnToolWindows.zip) |
| Linux (x64, glibc) | [EnToolLinux.tar.gz](downloads/EnToolLinux.tar.gz) |
| macOS (Intel x64) | [EnToolMacOS.tar.gz](downloads/EnToolMacOS.tar.gz) |

Open your platform's file link and use **Download raw file** on GitHub to save it.

```text
+--------------------------------------------------------------------+
| WINDOWS                                                            |
+--------------------------------------------------------------------+
```

Download `EnToolWindows.zip`, right-click it and choose **Extract All**, then open the extracted `EnToolWindows` folder and launch `EnToolWindows.exe`. **Keep the `_internal` folder together with the EXE**; it contains runtime dependencies. Encrypted `.txt` files will appear next to `EnToolWindows.exe`.

This packaging avoids the temporary self-extraction used by PyInstaller `--onefile`. It **does not guarantee** that Microsoft Defender, SmartScreen, or other security software will trust the binary. If Windows blocks it, do **not** turn off protection or add a blind exclusion. Use the Python source version below or submit the detection to [Microsoft Security Intelligence](https://www.microsoft.com/en-us/wdsi/filesubmission) for review.

```text
+--------------------------------------------------------------------+
| LINUX                                                              |
+--------------------------------------------------------------------+
```

Download `EnToolLinux.tar.gz`, extract it in a folder, then run:

```bash
tar -xzf EnToolLinux.tar.gz
chmod +x EnToolLinux
./EnToolLinux
```

```text
+--------------------------------------------------------------------+
| MACOS                                                              |
+--------------------------------------------------------------------+
```

Download `EnToolMacOS.tar.gz`, extract it in a folder, then run:

```bash
tar -xzf EnToolMacOS.tar.gz
chmod +x EnToolMacOS
./EnToolMacOS
```

The macOS build targets **Intel x64**. It may run on Apple Silicon Macs through Rosetta 2. Because the executable is not Apple-signed or notarized, macOS may prevent launching it until the user reviews the security warning and explicitly authorizes it. Do not bypass security warnings for files from untrusted sources.

```text
+--------------------------------------------------------------------+
| HOW TO USE                                                         |
+--------------------------------------------------------------------+
```

1. Launch the executable and select `1` to **encrypt** a secret.
2. Enter a **file name** (such as `message.txt`), the secret text, an encryption passphrase, and the desired output length in characters.
3. EnTool creates the encrypted `.txt` file next to the executable (inside the extracted `EnToolWindows` folder on Windows).
4. Send the `.txt` file to another person **separately from the passphrase**.
5. The recipient places that `.txt` file next to their EnTool executable, selects `2`, enters the file name and correct passphrase, and sees the decrypted text.
6. Select `0` to quit.

The same encrypted text file can be decrypted on Windows, Linux, or macOS. Different files can use different passphrases. Neither the secret text nor the passphrase is stored as plaintext in the encrypted file.

```text
+--------------------------------------------------------------------+
| SECURITY AND LIMITATIONS                                           |
+--------------------------------------------------------------------+
```

- Uses **AES-256-GCM** authenticated encryption with a key derived from the passphrase using **scrypt**.
- Encrypted files have an `ET1` format identifier and printable random-looking contents.
- The output length is customizable, but **longer ciphertext does not make weak passphrases stronger**.
- Choose a long, unique passphrase. **If you lose it, EnTool cannot recover your secret.**
- Only regular `.txt` files in the same folder as the executable can be decrypted.
- The decrypted secret is printed on-screen. Be careful with terminal recordings and screen sharing.
- The program operates offline and does not contact any network service.
- **EnTool has not undergone an independent security audit.**
- Windows binaries are not Authenticode-signed. A ZIP archive and SHA-256 checksums do **not** constitute a malware clearance or trusted publisher identity. The checksums in [`downloads/SHA256SUMS`](downloads/SHA256SUMS) are for file integrity verification only.
- Compatibility: Windows x64; Linux x64 on compatible glibc-based distributions (built on Ubuntu 22.04); macOS Intel x64 (built on macOS 15). Other architectures and older systems are not guaranteed.

```text
+--------------------------------------------------------------------+
| SOURCE CODE                                                        |
+--------------------------------------------------------------------+
```

The implementation is in [`EnTool.py`](EnTool.py). The README logo and `by christofrancis` line are displayed once at startup. For platform-specific executable builds, see [`BUILDING.md`](BUILDING.md). Developers may install dependencies from `requirements.txt` and run the included tests:

```bash
python -m pip install -r requirements.txt
python -m unittest discover -s tests -v
python EnTool.py
```

If Windows blocks the prebuilt executable, you can run the readable Python source with Python 3.12 installed and its dependencies from `requirements.txt`. Do not bypass Windows security warnings to force the downloaded executable to run.

For integrity verification of the published archives, check their SHA-256 hashes against [`downloads/SHA256SUMS`](downloads/SHA256SUMS). On Windows, PowerShell can calculate the checksum with:

```powershell
Get-FileHash .\EnToolWindows.zip -Algorithm SHA256
```

```text
+--------------------------------------------------------------------+
| LICENSE                                                            |
+--------------------------------------------------------------------+
```

[MIT License](LICENSE) — copyright (c) 2026 christofrancis.
