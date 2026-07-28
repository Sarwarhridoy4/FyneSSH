# FyneSSH — Implementation Progress Report

**Overall completion: ~25%**

---

## 2. Functional Requirements

| # | Requirement | Status | Notes |
|---|-------------|--------|-------|
| 2.1 | Remote Login | **Partial** | Login tab exists but launches external system terminal (`ssh user@host -p port`). Does **not** maintain a persistent terminal session inside the app as specified. |
| 2.2 | Command Execution | **Not Started** | No in-app terminal emulator or remote command execution interface. |
| 2.3 | File Transfer (SCP/SFTP) | **Not Started** | No SFTP/SCP integration. |
| 2.4 | Tunneling / Port Forwarding | **Not Started** | No local/remote port forwarding UI or logic. |
| 2.5 | Git Authentication via SSH | **Not Started** | No Git remote detection or Git-specific auth testing. |
| 2.6 | SSH Key-Based Auth Management | **Partial (70%)** | Key generation (Ed25519/RSA), display, save to `~/.ssh/`, copy to clipboard, and upload to remote `authorized_keys` are implemented. **Passphrase support is not implemented** (UI field exists but `keygen.Generate` ignores it). `ssh-copy-id` was replaced with a manual upload flow. |
| 2.7 | SSH Config Alias Management | **Partial (80%)** | Create/read/update `~/.ssh/config` Host blocks with all specified fields (Host, HostName, User, Port, IdentityFile, IdentitiesOnly, AddKeysToAgent, ServerAliveInterval, ServerAliveCountMax). **Missing:** config validation and alias connectivity testing. |
| 2.8 | Custom SSH Port Configuration (Server-side) | **Not Started** | No `sshd_config` editing, `sshd -t` validation, UFW integration, or service restart. |
| 2.9 | UFW Firewall Management | **Not Started** | No UFW rule management. |
| 2.10 | SSH Service Management | **Not Started** | No start/stop/restart/enable/disable for `ssh`/`sshd`. Distro detection exists in `internal/platform/distro.go` but is unused. |
| 2.11 | Cross-platform File System Access | **Partial (50%)** | Linux `~/.ssh/` paths work. Windows `%USERPROFILE%\.ssh` pathing is not explicitly handled. |

---

## 3. Non-Functional Requirements

| Requirement | Status | Notes |
|-------------|--------|-------|
| Security | **Partial** | Private keys are not logged, and file permissions are set correctly (0600/0644). However, passphrase encryption is not implemented, and the in-app Instructions tab displays the full private key in a read-only text field (by design, but still a risk if screenshared). |
| Performance | **Unknown** | No stress testing; terminal streaming is delegated to external terminal. |
| Reliability | **Basic** | Errors are wrapped and shown to user. No retry/backoff logic. |
| Usability | **Basic** | Tabbed UI with Instructions. No keyboard shortcuts or dark mode toggle. |
| Extensibility | **Partial** | Backend packages (`keygen`, `platform`, `sshclient`) are separated from UI, but interfaces are not used for dependency injection. |

---

## 4. Architecture Guidelines

| Requirement | Status | Notes |
|-------------|--------|-------|
| Go 1.21+ | **Done** | Using Go 1.22+. |
| Fyne v2.4+ | **Done** | Using Fyne v2.4.0. |
| `golang.org/x/crypto/ssh` | **Done** | Used for SSH client. |
| `github.com/pkg/sftp` | **Not Started** | SFTP not implemented. |
| `ssh-copy-id` equivalent | **Partial** | Custom upload implemented instead. |
| Concurrency (goroutines/channels) | **Partial** | Used in some places, but no streaming terminal UI. |
| `context.Context` | **Partial** | Used in `sshclient.Dial`, but not consistently across all operations. |
| Config persistence | **Not Started** | No JSON config for app settings/saved sessions. |

---

## 5. Coding Standards

| Requirement | Status | Notes |
|-------------|--------|-------|
| `go fmt` | **Passes** | `go fmt ./...` clean. |
| `go vet` | **Passes** | Clean. |
| `golangci-lint` | **Not Checked** | Not run. |
| Unit tests | **Not Started** | Zero test files in project. |
| No global mutable state | **Mostly** | No globals found. |
| Doc comments | **Partial** | Some packages/functions lack doc comments. |
| File size < 500 lines | **Partial** | `pkg/ui/main.go` is ~780 lines. |

---

## 6. UI / UX Standards

| Requirement | Status | Notes |
|-------------|--------|-------|
| Tabs: Terminal, File Manager, Keys, UFW, Port Config | **Partial** | Only Login, Keys, and Instructions tabs exist. |
| Icons | **Partial** | Some Fyne theme icons used. |
| Loading indicators | **Not Started** | No progress spinners. |
| Disable actions during ops | **Not Started** | Buttons not disabled during upload/key generation. |
| Dark mode | **Not Started** | Uses Fyne default. |

---

## 7. Error Handling and Logging

| Requirement | Status | Notes |
|-------------|--------|-------|
| Structured logging (`slog`/`zap`) | **Not Started** | No structured logging. |
| Actionable error messages | **Partial** | Status labels show errors, but no dialogs. |
| Retry with backoff | **Not Started** | No retry logic. |

---

## 8. Testing Requirements

| Requirement | Status | Notes |
|-------------|--------|-------|
| Unit tests (`sshclient`, `keygen`, etc.) | **Not Started** | No tests. |
| Table-driven tests | **Not Started** | — |
| Mock SSH server | **Not Started** | — |
| Fyne test harness | **Not Started** | — |
| 70% coverage | **Not Started** | — |

---

## 9. Documentation Requirements

| Requirement | Status | Notes |
|-------------|--------|-------|
| README | **Unknown** | Not reviewed in this session. |
| Code comments | **Partial** | Some exported types/functions lack comments. |
| In-app user guide | **Done** | Instructions tab with markdown content. |

---

## 11. Branching Strategy

| Requirement | Status | Notes |
|-------------|--------|-------|
| Feature branches | **Done** | Work done on `feature/ssh-config-and-instructions` and `feature/organize-keys-tab`. |
| No direct push to main | **Done** | Changes pushed via branches. |

---

## 13. Codebase Structure

| Requirement | Status | Notes |
|-------------|--------|-------|
| `internal/sshclient` | **Done** | Basic dial/session/run. |
| `internal/sftpclient` | **Not Started** | Package does not exist. |
| `internal/keygen` | **Done** | Ed25519/RSA generation, save. Passphrase pending. |
| `internal/ufw` | **Not Started** | Package does not exist. |
| `internal/configparser` | **Partial** | `internal/platform/sshconfig.go` handles SSH config, but not a dedicated `configparser` package. |
| `internal/platform` | **Done** | Paths, shell detection, SSH config, known_hosts. Distro detection unused. |
| `pkg/ui` | **Partial** | Main window exists but is ~780 lines. Missing Terminal, File Manager, UFW, Port Config views. |
| `pkg/store` | **Not Started** | No persisted app settings. |
| `cmd/` entrypoint | **Done** | `cmd/fyne-ssh/main.go` exists. |

---

## 14. Review Checklist

| Item | Status |
|------|--------|
| `go fmt ./...` passes | ✅ |
| `golangci-lint run` passes | ❌ Not run |
| All tests pass | ❌ No tests |
| No private key/passphrase in logs | ✅ (passphrase not implemented yet) |
| UI responsive during SSH ops | ❌ Operations block; terminal is external |
| File permissions enforced | ✅ 0600/0644 |
| UFW/SSH service commands with confirmation | ❌ Not implemented |
| `sshd -t` validation | ❌ Not implemented |
| Config alias with IdentitiesOnly/keepalive | ✅ |
| README and AGENT.md up to date | ❓ README not reviewed |

---

## Summary

**Completed (~25%):**
- Project scaffolding, Fyne UI shell, tab structure
- SSH key generation (Ed25519/RSA), save to `~/.ssh/`
- SSH config read/write/update with all requested fields
- known_hosts management and host key callback
- System terminal launch for SSH sessions
- Instructions tab with usage docs
- Basic error handling and status messaging

**Not Started (~75%):**
- In-app terminal / command execution
- SFTP/SCP file transfer
- Port forwarding
- Git SSH integration
- Server-side SSH port configuration and UFW
- SSH service management
- Passphrase encryption for keys
- Tests, structured logging, config persistence
- Distro-specific service commands
- README updates

**Biggest gaps vs AGENT.md:** no terminal emulator, no file transfer, no tunneling, no UFW/service management, no tests, and no passphrase support.
