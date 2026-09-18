#!/usr/bin/env python3
"""Collect exact license notices from the locked Go module cache."""
import json
from pathlib import Path
import subprocess

ROOT = Path(__file__).resolve().parent.parent

def main():
    paths = subprocess.check_output(
        ["go", "list", "-deps", "-f", "{{if .Module}}{{.Module.Path}}{{end}}",
         "./cmd/scrcpy-tui"], cwd=ROOT, text=True
    ).splitlines()
    paths = sorted(set(path for path in paths if path))
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
        if not licenses:
            raise RuntimeError(f"No license notice found for {module['Path']}")
        chunks.append(f"\n## {module['Path']} {module['Version']}\n")
        for file in licenses:
            text = file.read_text(encoding="utf-8")
            chunks.append(f"\n{file.name}\n\n```text\n{text.rstrip()}\n```\n")
    (ROOT / "THIRD_PARTY_NOTICES.md").write_text("".join(chunks), encoding="utf-8")

if __name__ == "__main__":
    main()
