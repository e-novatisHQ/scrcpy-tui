#!/usr/bin/env python3
"""Compare archive bytes and report only public payload hashes on mismatch."""
import hashlib
from pathlib import Path
import sys
import zipfile

left, right = map(Path, sys.argv[1:])
if left.read_bytes() == right.read_bytes():
    print('Cross-host ZIP reproducibility: PASS')
    raise SystemExit(0)
for path in (left, right):
    print(path.name, path.parent.name)
    with zipfile.ZipFile(path) as archive:
        for info in archive.infolist():
            print(info.filename, info.file_size,
                  hashlib.sha256(archive.read(info.filename)).hexdigest(),
                  info.date_time, info.compress_type, info.create_system,
                  info.external_attr, info.flag_bits)
raise SystemExit('Linux/Windows ZIP builds differ')
