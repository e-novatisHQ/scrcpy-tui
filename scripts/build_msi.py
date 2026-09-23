#!/usr/bin/env python3
"""Build a per-user MSI using the explicitly provisioned, pinned wixl container."""
import argparse
import hashlib
from pathlib import Path
import re
import shutil
import subprocess
import tempfile
import uuid
from xml.sax.saxutils import escape

ROOT = Path(__file__).resolve().parent.parent
IMAGE = "scrcpy-tui-msi-toolchain:0.106"
UPGRADE_CODE = "26D5D8DF-48CB-568E-972F-C7D0320C89EC"
NAMESPACE = uuid.UUID(UPGRADE_CODE)


def build(stage, version, output):
    if not re.fullmatch(r"\d+\.\d+\.\d+", version):
        raise ValueError("MSI requires a stable major.minor.patch version")
    if any(n > limit for n, limit in zip(map(int, version.split('.')), (255, 255, 65535))):
        raise ValueError("version exceeds Windows Installer limits")
    names = ("scrcpy-tui.exe", "LICENSE", "README.md", "THIRD_PARTY_NOTICES.md")
    content = hashlib.sha256()
    for name in names:
        content.update(name.encode())
        content.update((stage / name).read_bytes())
    product = str(uuid.uuid5(NAMESPACE, "product:"+version)).upper()
    package = str(uuid.uuid5(NAMESPACE, "package:"+version+":"+content.hexdigest())).upper()
    component = str(uuid.uuid5(NAMESPACE, "per-user-amd64-files")).upper()
    with tempfile.TemporaryDirectory(prefix="scrcpy-msi-") as temporary:
        work = Path(temporary)
        for name in names:
            shutil.copyfile(stage / name, work / name)
        files = '\n'.join(f'<File Id="File{i}" Name="{escape(name)}" Source="{escape(name)}" />' for i, name in enumerate(names))
        (work / "product.wxs").write_text(f'''<?xml version="1.0" encoding="utf-8"?>
<Wix xmlns="http://schemas.microsoft.com/wix/2006/wi">
 <Product Id="{product}" UpgradeCode="{UPGRADE_CODE}" Name="scrcpy-tui" Version="{version}" Manufacturer="e-novatisHQ" Language="1033">
  <Package Id="{package}" InstallerVersion="500" Compressed="yes" InstallScope="perUser" />
  <MajorUpgrade DowngradeErrorMessage="A newer scrcpy-tui version is already installed." />
  <Media Id="1" Cabinet="payload.cab" EmbedCab="yes" />
  <Property Id="ARPNOMODIFY" Value="1" />
  <Directory Id="TARGETDIR" Name="SourceDir">
   <Directory Id="LocalAppDataFolder">
    <Directory Id="ProgramsFolder" Name="Programs">
     <Directory Id="INSTALLFOLDER" Name="scrcpy-tui">
      <Component Id="AppFiles" Guid="{component}" Win64="yes">
       {files}
       <RegistryValue Id="InstalledVersion" Root="HKCU" Key="Software\\e-novatisHQ\\scrcpy-tui" Name="InstalledVersion" Type="string" Value="{version}" KeyPath="yes" />
       <Environment Id="UserPath" Name="PATH" Action="set" Part="last" System="no" Permanent="no" Value="[INSTALLFOLDER]" />
       <RemoveFolder Id="RemoveInstallFolder" On="uninstall" />
      </Component>
     </Directory>
    </Directory>
   </Directory>
  </Directory>
  <Feature Id="Main" Title="scrcpy-tui" Level="1"><ComponentRef Id="AppFiles" /></Feature>
 </Product>
</Wix>
''', encoding="utf-8")
        subprocess.run(["docker", "run", "--rm", "--network", "none", "--mount",
                        f"type=bind,src={work},dst=/build", IMAGE,
                        "-a", "x64", "-o", "product.msi", "product.wxs"], check=True)
        # wixl 0.106 omits the uninstall marker despite Permanent="no" and
        # generates a random Environment row key. Set the exact MSI contract.
        (work / "Environment.idt").write_text(
            "Environment\tName\tValue\tComponent_\n"
            "s72\tl64\tL255\ts72\nEnvironment\tEnvironment\n"
            "UserPath\t=-PATH\t[~];[INSTALLFOLDER]\tAppFiles\n", encoding="utf-8")
        subprocess.run(["docker", "run", "--rm", "--network", "none", "--mount",
                        f"type=bind,src={work},dst=/build", "--entrypoint", "msibuild", IMAGE,
                        "product.msi", "-q", "DELETE FROM `Environment`", "-i", "Environment.idt"], check=True)
        output.parent.mkdir(parents=True, exist_ok=True)
        shutil.copyfile(work / "product.msi", output)
    print(f"Built {output.name} (product {product})")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--stage", type=Path, required=True)
    parser.add_argument("--version", required=True)
    parser.add_argument("--output", type=Path, required=True)
    parser.add_argument("--test-fixture", action="store_true")
    args = parser.parse_args()
    if not args.test_fixture and args.version != (ROOT / "VERSION").read_text().strip():
        parser.error("version must match VERSION; --test-fixture is only for isolated upgrade tests")
    build(args.stage, args.version, args.output)


if __name__ == "__main__":
    main()
