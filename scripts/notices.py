#!/usr/bin/env python3
"""Collect exact license notices from the locked Go module cache."""
import json
import os
from pathlib import Path
import subprocess

ROOT = Path(__file__).resolve().parent.parent

def main():
    paths = set()
    for system, arch in (("linux", "amd64"), ("linux", "arm64"), ("windows", "amd64")):
        output = subprocess.check_output(
            ["go", "list", "-deps", "-f", "{{if .Module}}{{.Module.Path}}{{end}}",
             "./cmd/scrcpy-tui"], cwd=ROOT, text=True,
            env={**os.environ, "GOOS": system, "GOARCH": arch, "CGO_ENABLED": "0"},
        )
        paths.update(path for path in output.splitlines() if path)
    paths = sorted(paths)
    raw = subprocess.check_output(
        ["go", "list", "-m", "-json", *paths], cwd=ROOT, text=True
    )
    decoder = json.JSONDecoder()
    modules = []
    while raw.strip():
        value, offset = decoder.raw_decode(raw.lstrip())
        raw = raw.lstrip()[offset:]
        if not value.get("Main"):
            modules.append(value)
    chunks = ["# Third-party notices\n\nGenerated from go.mod/go.sum by `make notices`. "
              "adb and scrcpy are external tools and are not bundled.\n"]
    for module in sorted(modules, key=lambda item: item["Path"]):
        directory = module.get("Dir")
        if not directory:
            raise RuntimeError(f"Missing module cache directory: {module['Path']}; run go mod download")
        path = Path(directory)
        licenses = sorted(file for file in path.iterdir() if file.is_file() and
                          (file.name.lower().startswith("license") or file.name.lower().startswith("copying")))
        # This exact locked Windows dependency declares MIT and its author only
        # in README.md. Preserve the upstream declaration verbatim; do not invent
        # a copyright year or silently apply this exception to future versions.
        if not licenses and (module["Path"], module["Version"]) == ("github.com/mattn/go-localereader", "v0.0.1"):
            declaration = path / "README.md"
            if "## License\n\nMIT\n" in declaration.read_text(encoding="utf-8"):
                licenses = [declaration]
        if not licenses:
            raise RuntimeError(f"No license notice found for {module['Path']}")
        chunks.append(f"\n## {module['Path']} {module['Version']}\n")
        for file in licenses:
            text = file.read_text(encoding="utf-8")
            chunks.append(f"\n{file.name}\n\n```text\n{text.rstrip()}\n```\n")
    (ROOT / "THIRD_PARTY_NOTICES.md").write_text("".join(chunks), encoding="utf-8")

if __name__ == "__main__":
    main()
