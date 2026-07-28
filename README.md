# FyneSSH

A desktop SSH client and server management tool built with [Fyne](https://fyne.io/).

> **Note:** This project is under active development. Several features from the roadmap are not yet implemented. See `progress.md` for the current status.

## Currently Implemented

- **SSH Key Management** — Generate Ed25519 or RSA 4096 key pairs, save to `~/.ssh/`, copy public key, upload to remote `authorized_keys`
- **Key Alias Support** — Use custom key filenames like `id_ed25519_personal` to avoid confusion
- **SSH Config Management** — Create and update `~/.ssh/config` Host entries with IdentityFile, IdentitiesOnly, AddKeysToAgent, ServerAliveInterval, ServerAliveCountMax, and more
- **known_hosts Management** — Host key verification and automatic known_hosts updates on first connection
- **System Terminal Launch** — Open your default system terminal with `ssh user@host -p port` or `ssh <config-alias>`, reading all options from `~/.ssh/config`
- **Instructions Tab** — In-app usage guide

## Roadmap (Not Yet Implemented)

- In-app terminal emulator with real-time command execution
- SFTP/SCP file transfer
- Local and remote port forwarding
- Git authentication over SSH
- Server-side SSH port configuration with UFW integration
- SSH service management (start/stop/restart)
- Passphrase encryption for generated keys
- Structured logging and retry logic
- Unit tests and CI

## Tech Stack

- [Go](https://go.dev/) 1.22+
- [Fyne](https://fyne.io/) v2.4+
- `golang.org/x/crypto/ssh`

## Prerequisites

- Go 1.22+
- Fyne dependencies (see [Fyne setup](https://developer.fyne.io/started/))

## Build

```bash
go build -o fyne-ssh ./cmd/fyne-ssh
```

## Run

```bash
./fyne-ssh
```

## Project Structure

```
cmd/fyne-ssh/           Application entrypoint
internal/
  keygen/               SSH key pair generation and saving
  platform/             Path resolution, shell detection, SSH config, known_hosts
  sshclient/            SSH dial, session, and command execution
pkg/ui/                 Fyne UI: tabs for Login, Keys, and Instructions
```

## Documentation

- `progress.md` — Detailed implementation progress against `AGENT.md`
- `AGENT.md` — Full feature specification and coding standards
- In-app **Instructions** tab — Usage guide

## License

See [LICENSE](./LICENSE)
