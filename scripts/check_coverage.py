#!/usr/bin/env python3
"""Fail when a Go coverage profile falls below the configured percentage."""

import argparse
from pathlib import Path
import re
import subprocess


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("profile", type=Path)
    parser.add_argument("minimum", type=float)
    args = parser.parse_args()

    output = subprocess.check_output(
        ["go", "tool", "cover", f"-func={args.profile}"], text=True
    )
    match = re.search(r"^total:\s+\(statements\)\s+([0-9.]+)%$", output, re.MULTILINE)
    if match is None:
        raise SystemExit("coverage total not found")
    coverage = float(match.group(1))
    print(f"coverage: {coverage:.1f}% (minimum: {args.minimum:.1f}%)")
    if coverage < args.minimum:
        raise SystemExit("coverage is below the required minimum")


if __name__ == "__main__":
    main()
