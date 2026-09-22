#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
test_directory=$(mktemp -d "${TMPDIR:-/tmp}/scrcpy-tui-install-test.XXXXXX")
trap 'rm -rf "$test_directory"' EXIT HUP INT TERM

fixtures="$test_directory/fixtures/v9.8.7"
fake_bin="$test_directory/fake-bin"
prefix="$test_directory/prefix"
mkdir -p "$fixtures/archive" "$fake_bin"

cat > "$fixtures/archive/scrcpy-tui" <<'EOF'
#!/bin/sh
echo 9.8.7
EOF
chmod +x "$fixtures/archive/scrcpy-tui"
tar -czf "$fixtures/scrcpy-tui_9.8.7_linux_amd64.tar.gz" -C "$fixtures/archive" scrcpy-tui
(cd "$fixtures" && sha256sum scrcpy-tui_9.8.7_linux_amd64.tar.gz > SHA256SUMS)

cat > "$fake_bin/curl" <<'EOF'
#!/bin/sh
set -eu
destination=""
url=""
while [ "$#" -gt 0 ]; do
    case "$1" in
        -o) destination=$2; shift 2 ;;
        -*) shift ;;
        *) url=$1; shift ;;
    esac
done
cp "${FAKE_RELEASE_ROOT}/${url##*/}" "$destination"
EOF
chmod +x "$fake_bin/curl"

output=$(PATH="$fake_bin:/usr/bin:/bin" SHELL=/bin/bash FAKE_RELEASE_ROOT="$fixtures" \
    SCRCPY_TUI_RELEASE_BASE_URL=https://invalid.example \
    sh "$root/install.sh" --version 9.8.7 --prefix "$prefix")

test -x "$prefix/bin/scrcpy-tui"
test "$("$prefix/bin/scrcpy-tui")" = "9.8.7"
printf '%s\n' "$output" | grep -F "$prefix/bin is not in PATH" >/dev/null
printf '%s\n' "$output" | grep -F 'For this terminal only:' >/dev/null

printf 'invalid checksum\n' > "$fixtures/SHA256SUMS"
if PATH="$fake_bin:/usr/bin:/bin" FAKE_RELEASE_ROOT="$fixtures" \
    SCRCPY_TUI_RELEASE_BASE_URL=https://invalid.example \
    sh "$root/install.sh" --version 9.8.7 --prefix "$test_directory/rejected" \
    >"$test_directory/rejected.out" 2>&1; then
    echo "install_test: invalid checksum was accepted" >&2
    exit 1
fi
test ! -e "$test_directory/rejected/bin/scrcpy-tui"

echo "install_test: OK"
