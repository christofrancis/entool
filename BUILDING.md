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

**Windows (x64), in PowerShell:**

```powershell
python -m PyInstaller --noconfirm --clean --onedir --console --noupx --name EnToolWindows EnTool.py
Compress-Archive -Path dist\EnToolWindows -DestinationPath downloads\EnToolWindows.zip -Force
```

The ZIP contains `EnToolWindows/EnToolWindows.exe` and its `_internal/` runtime
folder. Do not move the EXE out of the extracted folder. The old standalone
`downloads/EnToolWindows.exe` was removed following a Microsoft Defender
detection (`Trojan:Win32/Ulthar.A!ml`) against that build. The replacement is
not guaranteed to be recognized as clean by all antivirus products.

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

The GitHub Actions workflow rebuilds the packages after changes to the source
or build configuration, tests each executable, and writes
`downloads/SHA256SUMS`. Its Windows job attempts a Microsoft Defender scan **only
when Defender is enabled on the hosted runner**; an unavailable scan is not a
clean bill of health. For a public Windows release, test the ZIP on a normal
Windows installation with up-to-date Defender and submit suspected false
positives through https://www.microsoft.com/en-us/wdsi/filesubmission.

To strengthen trust in a future Windows release, Authenticode-sign the EXE with
a trusted code-signing certificate belonging to the publisher; do not store
private signing keys in this repository. Signing helps establish publisher
identity but still cannot guarantee that antivirus will accept every release.

Rebuild all three downloads whenever `EnTool.py` changes. Existing encrypted
`ET1` files remain compatible; adding a startup banner does not change the
file format.
