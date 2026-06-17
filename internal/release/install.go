package release

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultRepo = "syugeeeeeeeeeei/MemADR"
	installName = "memadr-install.sh"
)

func InstallName() string {
	return installName
}

func InstallScript(repo string) string {
	if strings.TrimSpace(repo) == "" {
		repo = DefaultRepo
	}

	return fmt.Sprintf(`#!/bin/sh
set -eu

bin="memadr"
asset="memadr_linux_amd64"
repo="%s"
version="${MEMADR_VERSION:-latest}"
bin_dir="${MEMADR_INSTALL_DIR:-}"

die() {
  echo "$*" >&2
  exit 1
}

fetch() {
  url="$1"
  out="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$out"
    return
  fi
  if command -v wget >/dev/null 2>&1; then
    wget -qO "$out" "$url"
    return
  fi

  die "curl or wget is required"
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    -b|--bin-dir)
      [ "$#" -ge 2 ] || die "missing value for $1"
      bin_dir="$2"
      shift 2
      ;;
    -v|--version)
      [ "$#" -ge 2 ] || die "missing value for $1"
      version="$2"
      shift 2
      ;;
    -r|--repo)
      [ "$#" -ge 2 ] || die "missing value for $1"
      repo="$2"
      shift 2
      ;;
    -h|--help)
      cat <<'EOF'
Usage: sh memadr-install.sh [options]

Options:
  -b | --bin-dir DIR   install destination
  -v | --version TAG   release tag such as v0.1.0
  -r | --repo REPO     GitHub repo in owner/name form
EOF
      exit 0
      ;;
    *)
      die "unknown arg: $1"
      ;;
  esac
done

os_name="$(uname -s 2>/dev/null || echo unknown)"
arch="$(uname -m 2>/dev/null || echo unknown)"

[ "$os_name" = "Linux" ] || die "this installer only supports Linux"

case "$arch" in
  x86_64|amd64)
    ;;
  *)
    die "unsupported arch: $arch (need x86_64/amd64)"
    ;;
esac

if [ -z "$bin_dir" ]; then
  if [ -w /usr/local/bin ]; then
    bin_dir=/usr/local/bin
  elif [ -n "${HOME:-}" ]; then
    bin_dir="$HOME/.local/bin"
  else
    bin_dir=.
  fi
fi

mkdir -p "$bin_dir"

base="https://github.com/$repo/releases"
if [ "$version" = "latest" ]; then
  url="$base/latest/download/$asset"
else
  url="$base/download/$version/$asset"
fi

tmp="$(mktemp "${TMPDIR:-/tmp}/memadr.XXXXXX")"
trap 'rm -f "$tmp"' EXIT HUP INT TERM

fetch "$url" "$tmp"
chmod +x "$tmp"
mv "$tmp" "$bin_dir/$bin"

trap - EXIT HUP INT TERM

echo "installed $bin to $bin_dir/$bin"

case ":$PATH:" in
  *:"$bin_dir":*)
    ;;
  *)
    echo "add $bin_dir to PATH if needed"
    ;;
esac
`, repo)
}

func WriteInstallScript(dir string, repo string) (string, error) {
	path := filepath.Join(dir, InstallName())
	if err := os.WriteFile(path, []byte(InstallScript(repo)), 0o755); err != nil {
		return "", err
	}
	if err := os.Chmod(path, 0o755); err != nil {
		return "", err
	}
	return path, nil
}
