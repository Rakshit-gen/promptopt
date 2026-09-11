#!/bin/sh
# promptopt installer.
#
#   curl -fsSL https://promptopt.dev/install.sh | sh
#
# Downloads the release binary for your OS/architecture from GitHub Releases,
# verifies its SHA256 checksum, and installs it into a user-local bin
# directory. No Python, Node, Go, Homebrew, or Docker required.
#
# Environment overrides:
#   PROMPTOPT_VERSION   version to install (default: latest)
#   PROMPTOPT_BIN_DIR   install directory (default: $HOME/.local/bin)
#   PROMPTOPT_REPO      owner/repo (default: rakshit-gen/promptopt)

set -eu

REPO="${PROMPTOPT_REPO:-rakshit-gen/promptopt}"
VERSION="${PROMPTOPT_VERSION:-latest}"
BIN_DIR="${PROMPTOPT_BIN_DIR:-$HOME/.local/bin}"
BINARY="promptopt"

info() { printf '  %s\n' "$*"; }
err()  { printf 'error: %s\n' "$*" >&2; }
die()  { err "$*"; exit 1; }

need() { command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"; }

main() {
	need uname
	need mkdir
	need mv
	need chmod

	DOWNLOADER=""
	if command -v curl >/dev/null 2>&1; then
		DOWNLOADER="curl"
	elif command -v wget >/dev/null 2>&1; then
		DOWNLOADER="wget"
	else
		die "need curl or wget to download the release"
	fi

	os="$(uname -s)"
	arch="$(uname -m)"

	case "$os" in
		Linux)  os="linux" ;;
		Darwin) os="darwin" ;;
		*) die "unsupported operating system: $os (promptopt ships for macOS and Linux)" ;;
	esac

	case "$arch" in
		x86_64|amd64)  arch="amd64" ;;
		arm64|aarch64) arch="arm64" ;;
		*) die "unsupported architecture: $arch (promptopt ships for amd64 and arm64)" ;;
	esac

	if [ "$VERSION" = "latest" ]; then
		VERSION="$(resolve_latest)"
		[ -n "$VERSION" ] || die "could not determine the latest version; set PROMPTOPT_VERSION"
	fi
	# Normalize: accept both "1.2.3" and "v1.2.3".
	tag="$VERSION"
	case "$tag" in v*) ;; *) tag="v$tag" ;; esac
	num="${tag#v}"

	archive="${BINARY}_${num}_${os}_${arch}.tar.gz"
	base="https://github.com/${REPO}/releases/download/${tag}"

	tmp="$(mktemp -d "${TMPDIR:-/tmp}/promptopt.XXXXXX")"
	trap 'rm -rf "$tmp"' EXIT

	info "downloading ${archive} (${tag})"
	fetch "${base}/${archive}" "${tmp}/${archive}" || die "download failed: ${base}/${archive}"

	if fetch "${base}/checksums.txt" "${tmp}/checksums.txt" 2>/dev/null; then
		verify_checksum "$tmp" "$archive"
	else
		info "checksums.txt not found for this release; skipping verification"
	fi

	need tar
	tar -xzf "${tmp}/${archive}" -C "$tmp"
	[ -f "${tmp}/${BINARY}" ] || die "archive did not contain a ${BINARY} binary"

	mkdir -p "$BIN_DIR"
	mv "${tmp}/${BINARY}" "${BIN_DIR}/${BINARY}"
	chmod +x "${BIN_DIR}/${BINARY}"

	info "installed ${BINARY} ${tag} to ${BIN_DIR}/${BINARY}"
	post_install
}

resolve_latest() {
	url="https://api.github.com/repos/${REPO}/releases/latest"
	if [ "$DOWNLOADER" = "curl" ]; then
		curl -fsSL "$url"
	else
		wget -qO- "$url"
	fi | sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1
}

fetch() {
	# fetch URL DEST
	if [ "$DOWNLOADER" = "curl" ]; then
		curl -fsSL "$1" -o "$2"
	else
		wget -q "$1" -O "$2"
	fi
}

verify_checksum() {
	dir="$1"; file="$2"
	expected="$(grep " ${file}\$" "${dir}/checksums.txt" | awk '{print $1}')"
	[ -n "$expected" ] || { info "no checksum entry for ${file}; skipping verification"; return 0; }

	if command -v sha256sum >/dev/null 2>&1; then
		actual="$(sha256sum "${dir}/${file}" | awk '{print $1}')"
	elif command -v shasum >/dev/null 2>&1; then
		actual="$(shasum -a 256 "${dir}/${file}" | awk '{print $1}')"
	else
		info "no sha256 tool found; skipping verification"
		return 0
	fi

	[ "$expected" = "$actual" ] || die "checksum mismatch for ${file} (expected ${expected}, got ${actual})"
	info "checksum verified"
}

post_install() {
	case ":${PATH}:" in
		*":${BIN_DIR}:"*)
			info "run: promptopt --help"
			;;
		*)
			printf '\n'
			info "${BIN_DIR} is not on your PATH. Add it:"
			printf '\n'
			shell_name="$(basename "${SHELL:-sh}")"
			case "$shell_name" in
				zsh)  rc="~/.zshrc" ;;
				bash) rc="~/.bashrc" ;;
				fish) rc="~/.config/fish/config.fish" ;;
				*)    rc="your shell profile" ;;
			esac
			if [ "$shell_name" = "fish" ]; then
				printf '    fish_add_path %s\n\n' "$BIN_DIR"
			else
				printf '    echo '\''export PATH="%s:$PATH"'\'' >> %s\n\n' "$BIN_DIR" "$rc"
			fi
			info "then restart your shell and run: promptopt --help"
			;;
	esac
	printf '\n'
	info "set your Groq API key:"
	printf '\n'
	printf '    export GROQ_API_KEY="your-key"   # from https://console.groq.com/keys\n\n'
}

main "$@"
