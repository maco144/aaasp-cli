# aaasp-cli

CLI client for [AAASP](https://aaasp.ai) — a multi-tenant Agent-as-a-Service platform. Talks to the `/v1` JSON API.

---

## Install

```sh
curl -fsSL https://raw.githubusercontent.com/maco144/aaasp-cli/main/scripts/install.sh | sh
```

Or with Go:

```sh
go install github.com/maco144/aaasp-cli/cmd/aaasp@latest
```

Or download a prebuilt binary from the [releases page](https://github.com/maco144/aaasp-cli/releases).

---

## Configuration

```sh
export AAASP_API_KEY=your_key_here
export AAASP_BASE_URL=https://aaasp.ai   # optional, this is the default
```

No API key yet? Run `aaasp register` to create an account.

---

## Usage

```
aaasp <command> [subcommand] [flags]

Commands:
  register                            Create a new account
  whoami                              Show current account info

  deployments list                    List deployments (READY column + why any can't run)
  deployments create <agent_def_id>   Create a deployment
  deployments show <id>               Show a deployment, including runnable/reason
  deployments delete <id>             Delete a deployment

  runs list                                        List runs
  runs create <deployment_id> [prompt...] [--sync]  Trigger a run
  runs show <id>                                   Show run details
  runs cancel <id>                                 Cancel a run

  credentials list      List stored credentials
  credentials add       Add a credential
  credentials delete    Delete a credential

  skills list            List available skills

  mcp serve              Run an MCP server over stdio (dispatch_run, get_run,
                          list_deployments, list_agents) for MCP clients like
                          Claude Desktop or Claude Code

Flags:
  --json                 Output raw JSON
```

`runs create` submits directly against a known deployment; `mcp serve`'s
`dispatch_run` tool (or `POST /v1/dispatch`) is the alternative when you want
AAASP to auto-route the goal instead of naming a deployment yourself.

---

## MCP server

`aaasp mcp serve` runs an MCP (Model Context Protocol) server over stdio, exposing
the AAASP core loop as tools for any MCP client:

| Tool | Description |
|------|-------------|
| `dispatch_run` | Submit a goal — AAASP classifies it and auto-routes to the best agent |
| `get_run` | Check a run's status and result by ID |
| `list_deployments` | List agent deployments, each with `runnable` and `readiness_error` (why it can't run) |
| `list_agents` | List agent definitions for the authenticated tenant |

Example Claude Desktop / Claude Code config:

```json
{
  "mcpServers": {
    "aaasp": {
      "command": "aaasp",
      "args": ["mcp", "serve"],
      "env": { "AAASP_API_KEY": "your_key_here" }
    }
  }
}
```

Authentication reuses the same `AAASP_API_KEY` / `AAASP_BASE_URL` config as every
other command — one tenant, one API key, same scoping rules as the JSON API.

---

## Development

This is the extracted, publicly distributed counterpart of the `cli/` directory
that used to live in the AAASP monorepo — same split pattern as
[aaasp-ex](https://github.com/maco144/aaasp-ex). Develop and release from here;
the AAASP monorepo just links to this repo for CLI documentation.

```sh
make build     # build ./bin/aaasp for your platform
make install   # go install into $GOPATH/bin
make release   # cut a real release via goreleaser (needs GITHUB_TOKEN + a pushed tag)
```

---

## License

[Rising Sun License v1.0](./LICENSE) — free for personal/research use; commercial
deployments integrate with the Nous network.
