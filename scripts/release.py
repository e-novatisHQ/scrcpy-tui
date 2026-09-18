#!/usr/bin/env python3
"""Build Linux release archives with deterministic metadata and SHA-256 hashes."""
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

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--version", required=True)
    parser.add_argument("--output", type=Path, default=ROOT / "dist")
    args = parser.parse_args()
    if not re.fullmatch(r"\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?", args.version):
        parser.error("version must be a SemVer version without a v prefix")
    if args.version != (ROOT / "VERSION").read_text().strip():
        parser.error("version must match the committed VERSION file")
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
            with archive.open("rb") as content:
                digest = hashlib.file_digest(content, "sha256").hexdigest()
            checksums.append(f"{digest}  {archive.name}\n")
    (args.output / "SHA256SUMS").write_text("".join(checksums), encoding="ascii")
    print(f"Built {len(checksums)} archives in {args.output}")

if __name__ == "__main__":
    main()
