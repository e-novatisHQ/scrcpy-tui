#!/bin/sh
set -eu

repo="e-novatisHQ/scrcpy-tui"
prefix=${PREFIX:-"$HOME/.local"}
version=""

usage() {
    cat <<'EOF'
Usage: install.sh [--version VERSION] [--prefix DIRECTORY]

Install the latest scrcpy-tui release under ~/.local/bin by default.
EOF
}

while [ "$#" -gt 0 ]; do
    case "$1" in
        --version)
            [ "$#" -ge 2 ] || { echo "install.sh: --version requires a value" >&2; exit 2; }
            version=$2
            shift 2
            ;;
        --prefix)
            [ "$#" -ge 2 ] || { echo "install.sh: --prefix requires a value" >&2; exit 2; }
            prefix=$2
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "install.sh: unknown option: $1" >&2
            usage >&2
            exit 2
            ;;
    esac
done

for command_name in curl tar sha256sum install mktemp uname grep; do
    command -v "$command_name" >/dev/null 2>&1 || {
        echo "install.sh: required command not found: $command_name" >&2
        exit 1
    }
done

case $(uname -s) in
    Linux) ;;
    *) echo "install.sh: only Linux is supported" >&2; exit 1 ;;
esac

case $(uname -m) in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) echo "install.sh: unsupported CPU architecture: $(uname -m)" >&2; exit 1 ;;
esac

latest_url=${SCRCPY_TUI_LATEST_URL:-"https://github.com/$repo/releases/latest"}
release_base_url=${SCRCPY_TUI_RELEASE_BASE_URL:-"https://github.com/$repo/releases/download"}

if [ -z "$version" ]; then
    effective_url=$(curl -fsSL -o /dev/null -w '%{url_effective}' "$latest_url")
    version=${effective_url##*/}
fi
case "$version" in
    v*) tag=$version; release_version=${version#v} ;;
    *) tag=v$version; release_version=$version ;;
esac
printf '%s\n' "$tag" | grep -Eq '^v[0-9][0-9A-Za-z.+-]*$' || {
    echo "install.sh: invalid release version: $version" >&2
    exit 1
}

archive="scrcpy-tui_${release_version}_linux_${arch}.tar.gz"
temporary_directory=$(mktemp -d "${TMPDIR:-/tmp}/scrcpy-tui-install.XXXXXX")
trap 'rm -rf "$temporary_directory"' EXIT HUP INT TERM

echo "Downloading scrcpy-tui $release_version for linux/$arch..."
curl -fsSL -o "$temporary_directory/$archive" "$release_base_url/$tag/$archive"
curl -fsSL -o "$temporary_directory/SHA256SUMS" "$release_base_url/$tag/SHA256SUMS"

checksum_line=$(grep -E "^[[:xdigit:]]{64}  ${archive}$" "$temporary_directory/SHA256SUMS" || true)
[ -n "$checksum_line" ] || {
    echo "install.sh: $archive is absent from SHA256SUMS" >&2
    exit 1
}
printf '%s\n' "$checksum_line" | (cd "$temporary_directory" && sha256sum -c -)

mkdir_path="$prefix/bin"
install -d "$mkdir_path"
tar -xzf "$temporary_directory/$archive" -C "$temporary_directory" scrcpy-tui
install -m 0755 "$temporary_directory/scrcpy-tui" "$mkdir_path/scrcpy-tui"

echo "Installed $mkdir_path/scrcpy-tui"
case ":${PATH:-}:" in
    *":$mkdir_path:"*)
        echo "Run: scrcpy-tui"
        ;;
    *)
        echo ""
        echo "$mkdir_path is not in PATH. Add it, then open a new terminal:"
        case ${SHELL##*/} in
            zsh)
                echo "Add this line to ~/.zshrc:"
                echo "  export PATH=\"$mkdir_path:\$PATH\""
                ;;
            fish)
                echo "  fish_add_path \"$mkdir_path\""
                ;;
            *)
                echo "Add this line to ~/.bashrc:"
                echo "  export PATH=\"$mkdir_path:\$PATH\""
                ;;
        esac
        echo "For this terminal only:"
        echo "  export PATH=\"$mkdir_path:\$PATH\""
        ;;
esac
