# FyneSSH — Agent Specification

## 1. Project Identity

- **Name:** FyneSSH
- **Purpose:** Desktop SSH client and server management tool built with the Fyne UI toolkit.
- **Scope:** Cover all functionalities described in `instruction.md`.
- **Platform:** Primary development on Linux; secondary support for Windows client workflows.

## 2. Functional Requirements

### 2.1 Remote Login
- Support SSH login to a remote server using `ssh user@host -p port`.
- Accept hostname or IP address, custom port, and username.
- Maintain persistent terminal session in the UI.

### 2.2 Command Execution
- Provide a terminal-like interface to execute remote commands over the established SSH connection.
- Display stdout and stderr clearly.

### 2.3 File Transfer
- Integrate SCP and SFTP for secure file transfer.
- UI actions: upload, download, browse remote directory, list local files.

### 2.4 Tunneling / Port Forwarding
- Support local and remote port forwarding configuration.
- Allow users to define forwarding rules via the UI.
- Manage tunnel lifecycle (start/stop).

### 2.5 Git Authentication via SSH
- Auto-detect Git remote URLs over SSH.
- Use configured SSH identities for Git operations.
- Provide ability to test Git authentication in-app.

### 2.6 SSH Key-Based Authentication Management
- Generate SSH key pairs: Ed25519 (primary), RSA 4096 (fallback).
- Set optional passphrase for private key.
- Display generated keys in UI with safety warnings.
- Copy public key to clipboard for manual server deployment.
- Upload public key to remote server using `ssh-copy-id` equivalent.

### 2.7 SSH Config Alias Management
- Create and edit `~/.ssh/config` entries via UI.
- Fields per entry: Host (alias), HostName (IP or FQDN), User, Port, IdentityFile, IdentitiesOnly, ServerAliveInterval, ServerAliveCountMax.
- Validate config before saving.
- Test alias connectivity.

### 2.8 Custom SSH Port Configuration (Server-side)
- Edit `/etc/ssh/sshd_config` on connected Linux servers.
- Change default port, validate config with `sshd -t`.
- Manage UFW firewall rules for the new SSH port.
- Enable/disable UFW via UI.
- Allow or deny specific ports and IPs via UFW.
- Remove old SSH port rule from UFW.
- Restart `ssh.service` and `ssh.socket` after changes.
- Guide users through login on custom port.

### 2.9 UFW Firewall Management
- Display active rules, numbered rules.
- Add/delete allow/deny rules.
- Toggle UFW active/inactive.
- Show status in real time.

### 2.10 SSH Service Management
- Start, stop, restart, enable, disable, and check status of `ssh` or `sshd` service depending on distro family.
- Detect distro family (Debian-based vs Red Hat-based) automatically.
- Warn users before performing service restart.

### 2.11 Cross-platform File System Access
- On Linux/macOS: respect `~/.ssh` permissions (0700 for `.ssh`, 0600 for private keys).
- On Windows: support `%USERPROFILE%\.ssh` pathing.

## 3. Non-Functional Requirements

- **Security:** Never expose private key contents in logs or UI error messages. Mask sensitive fields.
- **Performance:** Terminal output must stream in real time without blocking UI (use goroutines).
- **Reliability:** Gracefully handle SSH connection drops, timeouts, and auth failures.
- **Usability:** Keyboard shortcuts for common actions. Dark mode support via Fyne preferences.
- **Extensibility:** Backend SSH logic decoupled from UI widgets for future reuse.

## 4. Architecture Guidelines

- **Language:** Go (1.21+).
- **UI Framework:** Fyne v2.4+.
- **Backend Packages:** `golang.org/x/crypto/ssh` for SSH client; `github.com/pkg/sftp` for SFTP; `github.com/c-balko/ssh-copy-id/go-ssh-copy-id` or internal equivalent for key upload.
- **Concurrency:** Use channels and goroutines for terminal streaming and file transfers. Avoid blocking the UI thread.
- **Error Handling:** Wrap errors with context. Display user-friendly messages.
- **Configuration:** Persist app settings and saved sessions using Fyne storage or local JSON.

## 5. Coding Standards

### 5.1 Go Standards
- Format code with `go fmt`.
- Lint with `golangci-lint` (enable `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`).
- Use `context.Context` for cancellable operations.
- Avoid global mutable state; prefer dependency injection.
- Keep functions small and single-responsibility.

### 5.2 Naming Conventions
- Go: standard library conventions.
- Fyne widgets: descriptive names (`sessionList`, `terminalOutput`, `sshOptionsTab`).
- Internal packages: short lowercase names (`sshclient`, `keygen`, `ufw`).

### 5.3 Security Rules
- Never log private keys or passphrases.
- Validate all user inputs before passing to SSH or system commands.
- Use `exec.CommandContext` with controlled arguments; avoid shell injection.
- Restrict file permissions when writing key files: 0600 for private, 0644 for public.

## 6. UI / UX Standards

- Follow Fyne container and widget conventions.
- Use icons (`ssh`, `key`, `folder`, `terminal`) consistently.
- Provide loading indicators for long-running operations.
- Disable actions during in-progress transfers or connections.
- Group related controls in tabs: Terminal, File Manager, Keys, UFW, Port Config.

## 7. Error Handling and Logging

- Use structured logging (`slog` or `zap`).
- Show actionable error messages in dialog or status bar.
- Retry transient network errors with backoff.

## 8. Testing Requirements

- Unit tests for backend logic: `sshclient`, `keygen`, `ufw`, `configparser`.
- Use Go test with table-driven tests.
- Mock SSH server using `github.com/gliderlabs/ssh` for integration tests.
- UI tests using Fyne test harness for critical workflows.
- Minimum coverage target: 70%.

## 9. Documentation Requirements

- README: setup, build, run.
- Code comments for exported functions and types.
- User guide embedded in-app or as markdown.

## 11. Branching Strategy

- Use a separate branch for each feature or fix.
- Branch naming convention: `feature/<short-description>` or `fix/<short-description>`.
- Do not commit directly to `main`.
- Keep branches focused and small; merge via pull request after review.

---

## 13. Codebase Structure and Documentation Standards

### 13.1 Split Codebase
- Organize the project into clearly separated packages:
  - `internal/sshclient` — SSH session, channel management, command execution
  - `internal/sftpclient` — SFTP file transfer operations
  - `internal/keygen` — SSH key pair generation and management
  - `internal/ufw` — UFW firewall rule management
  - `internal/configparser` — SSH config file read/write/validation
  - `internal/platform` — OS/distro detection and path resolution
  - `pkg/ui` — Fyne screens, widgets, and views
  - `pkg/store` — Persisted app settings and saved sessions
- Each package must have a single responsibility and expose a clean, typed API.
- Avoid coupling UI and backend logic; pass interfaces instead of concrete types.

### 13.2 Documentation Requirements
- Every exported package, variable, function, method, and type must have a doc comment (Go convention).
- Doc comments start with the name being documented and explain its purpose.
- Function comments describe behavior, parameters, return values, and errors.
- Complex algorithms and non-trivial logic require inline comments explaining intent.
- UI components and widgets must describe layout and user interaction.
- Internal documented utilities that are not exported still require clear comments.
- No undocumented public API; build must fail if documentation is missing.

### 13.3 Go and Fyne Standards
- Follow standard Go project layout; top-level `cmd/` contains the application entrypoint.
- All code passes `go fmt`, `go vet`, and `golangci-lint` without error.
- Use Go modules with minimal and well-justified dependencies.
- Fyne UI code must use widgets and containers idiomatically.
- Package names are short, lowercase, and stable.
- Group related declarations; avoid long files (> 500 lines). Split by responsibility when necessary.
- No hardcoded strings for user-facing text; prepare for localization.

---

## 14. Review Checklist

Before merging or releasing:
- [ ] `go fmt ./...` passes with no changes.
- [ ] `golangci-lint run` passes cleanly.
- [ ] All tests pass: `go test ./...`.
- [ ] No private key or passphrase appears in logs or UI.
- [ ] UI actions remain responsive during SSH operations.
- [ ] File permissions enforced when writing key files.
- [ ] UFW and SSH service commands run with confirmation dialogs.
- [ ] `sshd -t` validation step implemented before applying port changes.
- [ ] Config alias support includes IdentitiesOnly and keepalive parameters.
- [ ] README and AGENT.md are up to date.
