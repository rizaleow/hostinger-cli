# hostinger-cli

Go CLI for the [Hostinger API](https://developers.hostinger.com) — broader scope than the
official `hapi`, designed to be equally useful for humans and AI agents.

```
$ hostinger-cli vps vm restart 12345 --wait
$ hostinger-cli domains portfolio list -o json --jq '.[].domain'
$ hostinger-cli describe | jq '[.. | objects | select(.path) | .path]'
```

## Why this exists

Hostinger ships an official CLI, [`hostinger/api-cli`](https://github.com/hostinger/api-cli)
(`hapi`). It's a solid Go tool — but it only covers the **VPS** subset of the
public API, and is geared at humans at a terminal.

`hostinger-cli` is an alternative that:

- Covers the **entire Hostinger public API** — 119 commands across Billing, DNS,
  Domains, Domain-Access-Verifier, Hosting, Reach (email marketing), and VPS
  (incl. the experimental Docker Compose Manager).
- Is **scriptable and agent-friendly**: `--jq` filtering, structured errors
  under `-o json`, and a `describe` subcommand that emits the full command
  tree as JSON so an LLM can learn the surface without parsing `--help`.
- Handles **async ops** properly: `--wait` polls the VPS action endpoint
  until it reaches a terminal state, with `--wait-timeout`.
- Supports **named profiles** and **OS-keychain** token storage (opt-in).
- Ships through **Homebrew**, **GHCR Docker**, **`go install`**, and **GitHub
  release archives**.

The client layer is generated from Hostinger's `api.json` via
[oapi-codegen](https://github.com/oapi-codegen/oapi-codegen), so it tracks the
upstream spec instead of being hand-maintained.

### How it compares to `hapi`

|  | [`hapi`](https://github.com/hostinger/api-cli) (official) | `hostinger-cli` (this) |
|---|---|---|
| API scope | VPS only | Billing · DNS · Domains · Hosting · Reach · VPS (119 cmds) |
| Output formats | table · JSON · tree | table · JSON · YAML · Go-template (TTY-detected default) |
| Filter expression | — | `--jq <expr>` (pure-Go gojq) |
| Profiles | single token | named profiles |
| Token storage | plaintext `~/.hapi.json` | file or OS keychain (`--keyring`) |
| Async actions | — | `--wait` / `--wait-timeout` |
| Agent discovery | — | `describe` (JSON command tree) |
| Distribution | GitHub release tarballs | Homebrew tap · GHCR · `go install` · release archives |
| Auth env | `HAPI_API_TOKEN` | `HOSTINGER_API_TOKEN` |

Both are Go. This is **not** a rewrite for the sake of it — it exists for the
scope and UX rows above.

## Install

```bash
# Homebrew (macOS / Linux)
brew install rizaleow/tap/hostinger-cli

# Go
go install github.com/rizaleow/hostinger-cli/cmd/hostinger-cli@latest

# Docker
docker run --rm -e HOSTINGER_API_TOKEN ghcr.io/rizaleow/hostinger-cli:latest vps vm list

# Or grab a binary from GitHub Releases
# https://github.com/rizaleow/hostinger-cli/releases
```

## Authenticate

Get a token from the [Hostinger panel API page](https://hpanel.hostinger.com/profile/api).

```bash
# One-off
export HOSTINGER_API_TOKEN=hpe_...

# Or save it (writes ~/.config/hostinger-cli/config.yaml, chmod 0600)
hostinger-cli auth login

# Or store in the OS keychain instead of a file
hostinger-cli auth login --keyring

# Multiple accounts
hostinger-cli auth login --as-profile work
hostinger-cli --profile work vps vm list

hostinger-cli auth status
hostinger-cli auth whoami
```

Token resolution order: `--token` flag → `HOSTINGER_API_TOKEN` env →
OS keychain (if enabled) → config file.

## Examples

```bash
# VPS
hostinger-cli vps vm list
hostinger-cli vps vm get 12345
hostinger-cli vps vm restart 12345 --wait --wait-timeout 3m
hostinger-cli vps vm metrics 12345 --from 2026-01-01T00:00:00Z --to 2026-01-02T00:00:00Z

# Docker Compose Manager (experimental)
hostinger-cli vps docker project list 12345
hostinger-cli vps docker project create 12345 --from-file compose.json --wait
hostinger-cli vps docker project logs 12345 my-app

# DNS
hostinger-cli dns zone get example.com
hostinger-cli dns zone update example.com --from-file zone.json
hostinger-cli dns snapshot list example.com

# Domains
hostinger-cli domains availability check --domain mysite --tld com --tld dev
hostinger-cli domains portfolio list -o json --jq '.[] | .domain'
hostinger-cli domains portfolio lock example.com

# Billing
hostinger-cli billing subscriptions list
hostinger-cli billing subscriptions disable-auto-renewal sub_abc123

# Hosting
hostinger-cli hosting websites list --enabled
hostinger-cli hosting websites create --from-file site.json
```

## Output

Stdout is **JSON when piped or in CI**, **a table on a TTY** — same idea as
`gh`, `kubectl`, `doctl`. Override with `-o`:

```bash
hostinger-cli vps data-centers list                    # table on TTY
hostinger-cli vps data-centers list | jq '.[0]'        # auto JSON
hostinger-cli vps data-centers list -o yaml
hostinger-cli vps data-centers list -o json --jq '.[].location'
hostinger-cli vps data-centers list -o template --template '{{range .}}{{.name}}{{"\n"}}{{end}}'
```

`NO_COLOR` and `--no-color` disable ANSI styling. `--compact` produces
one-line JSON.

## For AI agents and scripts

```bash
# Discover the full command surface (119 commands, with flags and arg types)
hostinger-cli describe | jq '[.. | objects | select(.path) | .path]'

# Errors are structured under -o json:
#   {"error":{"code":"...","message":"...","correlation_id":"..."}}
# Exit code is non-zero on any non-2xx response.
```

## Develop

```bash
make generate    # regenerate internal/api/zz_generated.go from api.json
make test        # go test -race -count=1 ./...
make build       # bin/hostinger-cli
make snapshot    # goreleaser snapshot (no publish)
```

The OpenAPI spec lives at `api.json` at the repo root and is the
source of truth for the generated client.

## Releasing (maintainer)

One-time setup so the Homebrew tap actually works:

```bash
gh repo create rizaleow/homebrew-tap --public \
  --description "Homebrew tap for rizaleow's tools"
# In rizaleow/hostinger-cli → Settings → Secrets → Actions, add:
#   HOMEBREW_TAP_GITHUB_TOKEN  (PAT with `contents: write` on homebrew-tap)
```

Cut a release:

```bash
git tag v0.1.0
git push origin v0.1.0
# .github/workflows/release.yml runs goreleaser:
#   - cross-builds darwin/linux/windows × amd64/arm64
#   - publishes to GitHub Releases
#   - pushes the formula to rizaleow/homebrew-tap
#   - publishes ghcr.io/rizaleow/hostinger-cli multi-arch image
```

---

*Unofficial. Not affiliated with Hostinger International, UAB. `hapi` is © Hostinger and lives at [hostinger/api-cli](https://github.com/hostinger/api-cli).*
