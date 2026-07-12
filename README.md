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
  register              Create a new account
  whoami                Show current account info

  deployments list      List deployments
  deployments create    Create a deployment
  deployments delete    Delete a deployment

  runs list             List runs
  runs show <id>        Show run details
  runs create           Trigger a run
  runs cancel <id>      Cancel a run

  credentials list      List stored credentials
  credentials add       Add a credential
  credentials delete    Delete a credential

  skills list           List available skills

Flags:
  --json                Output raw JSON
```

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
