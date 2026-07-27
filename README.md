# FyneSSH

A desktop SSH client and server management tool built with [Fyne](https://fyne.io/).

## Features

- Remote SSH login with username, host/IP, and custom port
- Terminal-like command execution with real-time output
- Secure file transfer via SCP/SFTP
- Tunneling / port forwarding management
- Git authentication over SSH
- SSH key pair generation (Ed25519, RSA 4096)
- SSH config alias management (`~/.ssh/config`)
- Custom SSH port configuration for Ubuntu/Debian servers
- UFW firewall management
- SSH service management (start/stop/restart/status)
- Cross-platform support (Linux/macOS/Windows)

## Tech Stack

- [Go](https://go.dev/) 1.21+
- [Fyne](https://fyne.io/) v2.4+
- `golang.org/x/crypto/ssh`
- `github.com/pkg/sftp`

## Prerequisites

- Go 1.21+
- Fyne dependencies (See [Fyne setup](https://developer.fyne.io/started/))

## Build

```bash
go build -o fyne-ssh ./cmd/fyne-ssh
```

## Run

```bash
./fyne-ssh
```

## Usage

See `instruction.md` for detailed SSH concepts and server setup instructions.

## Run

Build the binary:

```bash
go build -o fyne-ssh ./cmd/fyne-ssh
```

Run the app:

```bash
./fyne-ssh
```

## License

See [LICENSE](./LICENSE)
