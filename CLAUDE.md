# hostinger-cli

Go CLI for the Hostinger API. Client is generated from `api.json` by
oapi-codegen; everything in `internal/api/zz_generated.go` is rebuilt — do
not hand-edit it.

## Commands
- `make build` / `make test` / `make generate` / `make update-spec`
- `./bin/hostinger-cli describe` — emit full command tree as JSON
- `make snapshot` — GoReleaser dry-run (no publish)

## Architecture
- `internal/api/` — generated client + transport (auth/retry/UA)
- `internal/clictx/` — bridge so subcommand packages can read root state
  without importing the `cli` package (avoids import cycle)
- `internal/output/` — TTY-detected table/json/yaml/template + gojq filter
- `internal/poll/` — generic `--wait` poller for VPS action resources
- `internal/cli/<tag>/` — one package per API tag (billing/dns/domains/…)

## Gotchas
- `dprotaso/go-yit` is silently pinned in `go.mod`. Bumping it pulls the
  `yaml/v4` RC and breaks `vmware-labs/yaml-jsonpath` → `make generate`
  fails. Re-pin: `v0.0.0-20220510233725-9ba8df137936`.
- `internal/api/tools.go` is `//go:build tools` — keeps `oapi-codegen` as
  a tool dep without leaking into the binary. Don't remove.
- Upstream serves the spec as JSON only at
  `https://developers.hostinger.com/openapi/openapi.json`. No YAML endpoint.
