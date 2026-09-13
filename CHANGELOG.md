# Changelog

## 0.4.1 (2026-09-13)

- Fixed: `aaasp --version` printed a hard-coded `0.3.0`. `version` was a Go
  `const`, which `-ldflags -X` cannot set, so the release version was never
  injected. v0.4.0 binaries report `0.3.0`, but their code is 0.4.0

## 0.4.0 (2026-09-13)

- `deployments list` adds a READY column and lists why each non-runnable
  deployment can't run (from the API's `runnable` / `readiness_error`)
- `deployments show` shows `runnable`, `reason`, `credential` and `created`.
  The `schedule` row was removed (deployments have no such field)
- `deployments create` warns on stderr when the new deployment can't run yet
- MCP `list_deployments` description tells agents about `runnable` /
  `readiness_error`, so they check them before dispatching
- API errors show the server's `message` alongside the error code
  (e.g. `not_runnable: This deployment has no anthropic key ...`)
- Fixed: `credentials add` sent the key as `vault_key`, which the API ignores,
  so every credential added through the CLI was stored without a key. It now
  sends `encrypted_key` and defaults the required label to the provider name
- Fixed: `runs list` / `runs show` read `inserted_at`, but the API returns
  `created_at`, so CREATED was always blank. `runs show` now also prints
  `error` and `result`
- First unit tests (`go test ./...`)

## 0.3.0 (2026-07-12)

- `aaasp deployments create <agent_def_id>` — `POST /v1/deployments`
- `aaasp runs create <deployment_id> [prompt...] [--sync]` — `POST /v1/runs`;
  `--sync` waits for a terminal result instead of returning immediately
- Both verified end-to-end against a live AAASP dev server

## 0.2.0 (2026-07-12)

- `aaasp mcp serve` — MCP (Model Context Protocol) server over stdio, exposing
  `dispatch_run`, `get_run`, `list_deployments`, and `list_agents` as tools for
  MCP clients (Claude Desktop, Claude Code, etc.)
- Fixed: `deployments create` and `runs create` were documented in `--help` but
  never implemented — removed from help/README pending an actual implementation

## 0.1.0 (2026-07-12)

Initial public release — extracted from the AAASP monorepo `cli/` directory.

- `aaasp register` / `whoami` — account bootstrap
- `aaasp deployments {list,delete}`
- `aaasp runs {list,show,cancel}`
- `aaasp credentials {list,add,delete}`
- `aaasp skills list`
- `--json` global flag for machine-readable output
- Cross-platform builds (linux/darwin/windows, amd64/arm64) via goreleaser
