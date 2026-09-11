# Installation

promptopt is a single static binary. It has no runtime dependencies — no
Python, Node, Go, Homebrew, or Docker.

## curl installer (recommended)

```sh
curl -fsSL https://promptopt.dev/install.sh | sh
```

The script:

1. Detects your OS (`macOS` or `Linux`) and architecture (`amd64` or `arm64`).
   It exits with a clear message on anything else.
2. Resolves the latest release tag from the GitHub API (or uses
   `PROMPTOPT_VERSION` if set).
3. Downloads `promptopt_<version>_<os>_<arch>.tar.gz` from GitHub Releases.
4. Downloads `checksums.txt` and verifies the archive's SHA256. If `sha256sum`
   and `shasum` are both missing, it warns and continues.
5. Extracts the binary and installs it to `~/.local/bin` (override with
   `PROMPTOPT_BIN_DIR`).
6. Prints the `PATH` line to add if `~/.local/bin` is not already on it.

### Installer environment variables

| Variable | Default | Purpose |
| --- | --- | --- |
| `PROMPTOPT_VERSION` | `latest` | Version to install, e.g. `0.2.0` or `v0.2.0`. |
| `PROMPTOPT_BIN_DIR` | `$HOME/.local/bin` | Install directory. |
| `PROMPTOPT_REPO` | `rakshit-gen/promptopt` | `owner/repo` to download from. |

### Inspect before running

If piping to `sh` makes you uneasy — reasonable — download and read it first:

```sh
curl -fsSL https://promptopt.dev/install.sh -o install.sh
less install.sh
sh install.sh
```

## Go

```sh
go install github.com/rakshit-gen/promptopt/cmd/promptopt@latest
```

Requires Go 1.24+. The binary lands in `$(go env GOPATH)/bin`. Version
metadata will read `dev` because `go install` does not pass build flags; use a
release binary or `make install` from a checkout for accurate `--version`
output.

## From source

```sh
git clone https://github.com/rakshit-gen/promptopt
cd promptopt
make install       # installs into $(go env GOPATH)/bin with version info
# or
make build         # produces ./bin/promptopt
```

## Homebrew

Once the tap is published:

```sh
brew install rakshit-gen/tap/promptopt
```

## Manual download

Every [release](https://github.com/rakshit-gen/promptopt/releases) attaches
`.tar.gz` archives for each platform and a `checksums.txt`. Verify and extract:

```sh
tar -tf promptopt_0.2.0_darwin_arm64.tar.gz
shasum -a 256 -c checksums.txt --ignore-missing
tar -xzf promptopt_0.2.0_darwin_arm64.tar.gz promptopt
sudo mv promptopt /usr/local/bin/
```

## Uninstall

```sh
rm ~/.local/bin/promptopt
rm -rf ~/.config/promptopt   # if you created a config file
```

## Shell completions

Cobra generates them:

```sh
promptopt completion zsh  > "${fpath[1]}/_promptopt"
promptopt completion bash > /etc/bash_completion.d/promptopt
promptopt completion fish > ~/.config/fish/completions/promptopt.fish
```
