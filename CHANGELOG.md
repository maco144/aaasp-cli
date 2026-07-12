# Changelog

## 0.1.0 (2026-07-12)

Initial public release — extracted from the AAASP monorepo `cli/` directory.

- `aaasp register` / `whoami` — account bootstrap
- `aaasp deployments {list,create,delete}`
- `aaasp runs {list,show,create,cancel}`
- `aaasp credentials {list,add,delete}`
- `aaasp skills list`
- `--json` global flag for machine-readable output
- Cross-platform builds (linux/darwin/windows, amd64/arm64) via goreleaser
