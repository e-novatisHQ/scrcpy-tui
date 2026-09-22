#!/usr/bin/env python3
"""Build deterministic Linux archives, Debian packages and SHA-256 hashes."""
import argparse
import gzip
import hashlib
import os
from pathlib import Path
import re
import shutil
import subprocess
import tarfile
import tempfile

ROOT = Path(__file__).resolve().parent.parent


def add_checksum(path, checksums):
    with path.open("rb") as content:
        digest = hashlib.file_digest(content, "sha256").hexdigest()
    checksums.append(f"{digest}  {path.name}\n")


def build_deb(binary, version, arch, output, temporary):
    package_root = temporary / f"deb_{arch}"
    control = package_root / "DEBIAN"
    executable_directory = package_root / "usr/bin"
    documentation = package_root / "usr/share/doc/scrcpy-tui"
    control.mkdir(parents=True)
    executable_directory.mkdir(parents=True)
    documentation.mkdir(parents=True)

    debian_version = version.replace("-", "~", 1)
    (control / "control").write_text(
        "\n".join([
            "Package: scrcpy-tui",
            f"Version: {debian_version}",
            "Section: utils",
            "Priority: optional",
            f"Architecture: {arch}",
            "Maintainer: e-novatisHQ <contact@e-novatis.com>",
            "Depends: adb, scrcpy",
            "Homepage: https://github.com/e-novatisHQ/scrcpy-tui",
            "Description: terminal launcher for scrcpy",
            " Select an ADB device and a preset, inspect the command, then launch scrcpy.",
            "",
        ]),
        encoding="utf-8",
    )
    shutil.copyfile(binary, executable_directory / "scrcpy-tui")
    shutil.copyfile(ROOT / "LICENSE", documentation / "copyright")
    shutil.copyfile(ROOT / "README.md", documentation / "README.md")
    shutil.copyfile(ROOT / "THIRD_PARTY_NOTICES.md", documentation / "THIRD_PARTY_NOTICES.md")

    for path in package_root.rglob("*"):
        if path.is_file():
            path.chmod(0o755 if path.name == "scrcpy-tui" else 0o644)
        elif path.is_dir():
            path.chmod(0o755)
        os.utime(path, (0, 0), follow_symlinks=False)
    package_root.chmod(0o755)
    os.utime(package_root, (0, 0))

    package = output / f"scrcpy-tui_{version}_linux_{arch}.deb"
    env = {**os.environ, "SOURCE_DATE_EPOCH": "0"}
    subprocess.run(
        ["dpkg-deb", "--root-owner-group", "--build", str(package_root), str(package)],
        env=env,
        check=True,
    )
    return package


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--version", required=True)
    parser.add_argument("--output", type=Path, default=ROOT / "dist")
    args = parser.parse_args()
    if not re.fullmatch(r"\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?", args.version):
        parser.error("version must be a SemVer version without a v prefix")
    if args.version != (ROOT / "VERSION").read_text().strip():
        parser.error("version must match the committed VERSION file")
    if shutil.which("dpkg-deb") is None:
        parser.error("dpkg-deb is required to build Debian packages")
    expected = {"linux_amd64", "linux_arm64"}
    args.output.mkdir(parents=True, exist_ok=True)
    notices = ROOT / "THIRD_PARTY_NOTICES.md"
    if not notices.exists():
        parser.error("run make notices before building archives")
    checksums = []
    with tempfile.TemporaryDirectory(prefix=".release-", dir=args.output) as temporary:
        for target in sorted(expected):
            stage = Path(temporary) / target
            stage.mkdir()
            arch = target.split("_")[1]
            env = {**os.environ, "GOOS": "linux", "GOARCH": arch, "CGO_ENABLED": "0"}
            subprocess.run(["go", "build", "-buildvcs=false", "-trimpath", "-ldflags",
                            f"-s -w -X main.version={args.version}", "-o", str(stage / "scrcpy-tui"),
                            "./cmd/scrcpy-tui"], cwd=ROOT, env=env, check=True)
            for name in ("LICENSE", "README.md", "THIRD_PARTY_NOTICES.md"):
                shutil.copyfile(ROOT / name, stage / name)
            archive = args.output / f"scrcpy-tui_{args.version}_{target}.tar.gz"
            with archive.open("wb") as raw:
                with gzip.GzipFile(filename="", mode="wb", fileobj=raw, mtime=0) as compressed:
                    with tarfile.open(fileobj=compressed, mode="w") as tar:
                        for file in sorted(stage.iterdir()):
                            info = tar.gettarinfo(str(file), arcname=file.name)
                            info.uid = info.gid = 0
                            info.uname = info.gname = ""
                            info.mtime = 0
                            info.mode = 0o755 if file.name == "scrcpy-tui" else 0o644
                            with file.open("rb") as content:
                                tar.addfile(info, content)
            add_checksum(archive, checksums)
            package = build_deb(stage / "scrcpy-tui", args.version, arch,
                                args.output, Path(temporary))
            add_checksum(package, checksums)
    (args.output / "SHA256SUMS").write_text("".join(checksums), encoding="ascii")
    print(f"Built {len(checksums)} release artifacts in {args.output}")

if __name__ == "__main__":
    main()
