#!/usr/bin/env python3
"""Verify MSI payload, user PATH contract and logical reproducibility offline."""
import argparse
import csv
import io
from pathlib import Path
import shutil
import subprocess
import tempfile

from build_msi import IMAGE, UPGRADE_CODE


def container(work, tool, *args):
    return subprocess.check_output([
        "docker", "run", "--rm", "--network", "none", "--mount",
        f"type=bind,src={work},dst=/build", "--entrypoint", tool, IMAGE, *args,
    ], text=True)


def table(work, msi, name):
    text = container(work, "msiinfo", "export", msi, name)
    lines = list(csv.reader(io.StringIO(text), delimiter='\t'))
    return [dict(zip(lines[0], row)) for row in lines[3:]]


def verify(msi, stage, compare):
    with tempfile.TemporaryDirectory(prefix="scrcpy-msi-verify-") as temporary:
        work = Path(temporary)
        shutil.copyfile(msi, work / "product.msi")
        container(work, "msiextract", "-C", "extracted", "product.msi")
        files = [p for p in (work / "extracted").rglob('*') if p.is_file()]
        expected = {"scrcpy-tui.exe", "LICENSE", "README.md", "THIRD_PARTY_NOTICES.md"}
        assert len(files) == len(expected) and {p.name for p in files} == expected, "unexpected MSI payload"
        for file in files:
            assert file.read_bytes() == (stage / file.name).read_bytes(), f"payload mismatch: {file.name}"
        properties = {r['Property']: r['Value'] for r in table(work, 'product.msi', 'Property')}
        assert not properties.get('ALLUSERS'), 'MSI is not per-user'
        assert properties['UpgradeCode'].strip('{}').upper() == UPGRADE_CODE
        assert table(work, 'product.msi', 'Environment') == [{
            'Environment': 'UserPath', 'Name': '=-PATH',
            'Value': '[~];[INSTALLFOLDER]', 'Component_': 'AppFiles',
        }], 'unsafe user PATH contract'
        registry = table(work, 'product.msi', 'Registry')
        assert all(r['Root'] == '1' for r in registry), 'registry write outside HKCU'
        actions = {r['Action'] for r in table(work, 'product.msi', 'InstallExecuteSequence')}
        assert {'WriteEnvironmentStrings', 'RemoveEnvironmentStrings', 'RemoveExistingProducts'} <= actions
        if compare:
            shutil.copyfile(compare, work / 'second.msi')
            # Summary timestamps/package code are documented as non-deterministic;
            # all installer decisions and payload hashes must remain identical.
            for name in ('Property','Directory','Component','File','Registry','Environment',
                         'Upgrade','InstallExecuteSequence','RemoveFile','Feature','FeatureComponents','Media'):
                first = table(work, 'product.msi', name)
                second = table(work, 'second.msi', name)
                assert first == second, f'non-deterministic MSI table: {name}'
    print('MSI payload, per-user registration, PATH and logical tables: PASS')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--msi', type=Path, required=True)
    parser.add_argument('--stage', type=Path, required=True)
    parser.add_argument('--compare', type=Path)
    args = parser.parse_args()
    verify(args.msi, args.stage, args.compare)


if __name__ == '__main__':
    main()
