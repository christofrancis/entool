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
python -m PyInstaller --noconfirm --clean --onefile --console --name EnToolWindows EnTool.py
Copy-Item dist\EnToolWindows.exe downloads\EnToolWindows.exe
```

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

Rebuild all three downloads whenever `EnTool.py` changes. Existing encrypted
`ET1` files remain compatible; adding a startup banner does not change the
file format.
