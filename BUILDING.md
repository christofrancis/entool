# Building EnTool executables

The startup logo lives in `EnTool.py` (`BANNER`), so it is included in every
platform build. Build **on the target operating system**, because PyInstaller
does not cross-compile its executables.

Install Python 3.12 and run from the repository root:

```sh
python -m pip install -r requirements.txt
python -m unittest discover -s tests -v
```

Then run the command for your operating system:

**Windows (x64), for local investigation only (not a verified release):**

```powershell
python -m PyInstaller --noconfirm --clean --onedir --console --noupx --name EnToolWindows EnTool.py
# Local testing only: do not publish this build until independently checked.
```

Windows downloadable binaries are suspended after Microsoft Defender detections
against Windows builds. Even the onedir format cannot be considered safe merely
because its smoke test passes. Do not publish Windows binaries until independent
verification, Defender scanning on a normal Windows machine, and appropriate
publisher signing/reputation checks have been completed.

**Linux (x64):**

```sh
python -m PyInstaller --noconfirm --clean --onefile --console --name EnToolLinux EnTool.py
tar -czf downloads/EnToolLinux.tar.gz -C dist EnToolLinux
```

**macOS (Intel x64):**

```sh
python -m PyInstaller --noconfirm --clean --onefile --console --name EnToolMacOS EnTool.py
# If necessary, ad-hoc sign this trusted local build before packaging:
# codesign --force --sign - dist/EnToolMacOS
tar -czf downloads/EnToolMacOS.tar.gz -C dist EnToolMacOS
```

For the macOS Intel x64 release, use an Intel machine or x86_64 Python under
Rosetta. A new macOS executable may require local ad-hoc code signing:
`codesign --force --sign - dist/EnToolMacOS`. This is *not* notarization;
Gatekeeper may still warn about downloaded unsigned/unnotarized software.

The GitHub Actions workflow currently builds and publishes Linux and macOS
packages only, runs startup tests, and writes `downloads/SHA256SUMS`. Windows
publishing has been deliberately disabled because the hosted Windows runner's
Defender was inactive during the previous attempt and users reported detections.
For any future Windows release, scan and inspect the actual build on Windows
with up-to-date Defender and investigate or submit suspected false positives
through https://www.microsoft.com/en-us/wdsi/filesubmission.

To strengthen trust in a future Windows release, Authenticode-sign the EXE with
a trusted code-signing certificate belonging to the publisher; do not store
private signing keys in this repository. Signing helps establish publisher
identity but still cannot guarantee that antivirus will accept every release.

Rebuild the published Linux/macOS packages whenever `EnTool.py` changes. Existing encrypted
`ET1` files remain compatible; adding a startup banner does not change the
file format.
