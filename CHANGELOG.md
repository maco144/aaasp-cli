# Changelog

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
