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

Prebuilt Linux and macOS packages are included **directly in this repository**. On Windows, **prebuilt executable downloads are temporarily suspended** after Microsoft Defender detections. Use the auditable Python source until a new Windows build has been independently checked.

```text
+--------------------------------------------------------------------+
| DOWNLOAD                                                           |
+--------------------------------------------------------------------+
```

| Operating system | Download |
|---|---|
| Windows (x64) | **Temporarily unavailable** — [run from Python source](#windows) |
| Linux (x64, glibc) | [EnToolLinux.tar.gz](downloads/EnToolLinux.tar.gz) |
| macOS (Intel x64) | [EnToolMacOS.tar.gz](downloads/EnToolMacOS.tar.gz) |

Open your platform's file link and use **Download raw file** on GitHub to save it.

```text
+--------------------------------------------------------------------+
| WINDOWS                                                            |
+--------------------------------------------------------------------+
```

**Windows downloads paused:** Microsoft Defender has reported detections against earlier EnTool Windows binaries (including `Trojan:Win32/Ulthar.A!ml`). Another serious detection was reported after the folder-bundle release, but its affected file path has not yet been verified. We withdrew the Windows ZIP as a precaution. Do not restore it from quarantine or turn off Microsoft Defender to run it.

To run the reviewed-in-repository source instead, install Python 3.12 from [python.org](https://www.python.org/downloads/windows/), clone or download the repository source, and open PowerShell in its root folder:

```powershell
py -3.12 -m pip install -r requirements.txt
py -3.12 EnTool.py
```

Encrypted `.txt` files are written next to `EnTool.py` when run from source. The Python code has not been independently audited; running from source is not a guarantee of safety. Verify anything you install and keep Microsoft Defender enabled. Windows binary distribution will resume only after further review and testing. Suspected false positives can be submitted to [Microsoft Security Intelligence](https://www.microsoft.com/en-us/wdsi/filesubmission).

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
3. EnTool creates the encrypted `.txt` file next to the executable (or next to `EnTool.py` when run from source).
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
- Windows binaries are currently **withdrawn** pending investigation and independent scanning. SHA-256 checksums in [`downloads/SHA256SUMS`](downloads/SHA256SUMS) verify only the integrity of the remaining Linux/macOS packages; checksums do **not** constitute malware clearance.
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

Windows prebuilt executable downloads are suspended. Run the Python source as documented in the Windows section. Do not bypass Windows security warnings to force an executable to run.

For integrity verification of the published Linux/macOS archives, check their SHA-256 hashes against [`downloads/SHA256SUMS`](downloads/SHA256SUMS).

```text
+--------------------------------------------------------------------+
| LICENSE                                                            |
+--------------------------------------------------------------------+
```

[MIT License](LICENSE) — copyright (c) 2026 christofrancis.
