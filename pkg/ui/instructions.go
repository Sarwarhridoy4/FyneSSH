package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) buildInstructionsTab() *fyne.Container {
	rawInstructions := `# FyneSSH - SSH Client

**FyneSSH** helps you manage SSH keys, organize ~/.ssh/config, and launch terminal sessions from one place.

---

## Tabs

- **Login** — Launch a terminal SSH session
- **Keys** — Generate keys, update SSH config, upload public keys
- **Instructions** — This help page

---

## Login Tab

Use this tab to quickly open a terminal and connect via SSH.

1. Enter **Host** — hostname or IP address
2. Enter **Port** — default is 22
3. Enter **User** — remote username
4. Click **Open Terminal**

FyneSSH launches your system terminal with:

    ssh user@host -p port

The terminal window handles the actual session. Password prompts and interactive shells run there, not inside the app.

---

## Keys Tab

### 1. Generate Key Pair

- **Algorithm** — ed25519 (recommended) or rsa
- **Key name** — Filename stem, e.g. id_ed25519_personal
- **Comment** — Optional key comment
- **Passphrase** — Optional encryption passphrase

Click **Generate** to create the key pair in memory.

### 2. Review Generated Keys

The **Private key** and **Public key** panels show the generated material.

- **Never share** the private key.
- The public key is in OpenSSH authorized_keys format.

### 3. Save Keys

Click **Save** to write files to ~/.ssh/:

- **Private key** — ~/.ssh/{key_name} (permissions: 0600)
- **Public key** — ~/.ssh/{key_name}.pub (permissions: 0644)

### 4. Save Key & Update Config

This saves the keys **and** writes a Host block into ~/.ssh/config.

- **Host alias** — Short config name, e.g. github, homelab
- **HostName** — Real hostname or IP
- **User** — Remote username
- **Port** — SSH port, default 22
- **AddKeysToAgent** — Auto-load key into ssh-agent
- **IdentitiesOnly** — Use only the specified key
- **ServerAliveInterval** — Keep-alive seconds
- **ServerAliveCountMax** — Max keep-alive probes

If the **Host alias** already exists in ~/.ssh/config, it will be **updated**.

#### Example config block

    Host github
        HostName github.com
        User git
        IdentityFile ~/.ssh/id_ed25519_personal
        AddKeysToAgent yes
        IdentitiesOnly yes

    Host homelab.com
        HostName 192.168.122.6
        User sarwar
        Port 9255
        IdentityFile ~/.ssh/id_ed25519
        IdentitiesOnly yes
        ServerAliveInterval 20
        ServerAliveCountMax 2

### 5. Open Terminal

After filling the **SSH Config Entry** section, click **Open Terminal**.

FyneSSH reads ~/.ssh/config and launches your system terminal with:

    ssh <host_alias>

All configured options — IdentityFile, port, user, keep-alive — are handled by your system SSH client.

### 6. Upload Public Key

Use this to copy the public key to a remote server.

1. Enter remote **Host**, **User**, **Port**, and **Password**
2. Click **Upload to Server**

The public key is appended to ~/.ssh/authorized_keys on the remote host.

### 7. Add to known_hosts

1. Enter the target hostname in **Upload Host**
2. Click **Add to known_hosts**

FyneSSH connects and records the host key in ~/.ssh/known_hosts to prevent man-in-the-middle attacks.

### 8. Copy Public Key

Click **Copy Public Key** to copy the public key to clipboard. Paste it manually into authorized_keys if you prefer.

---

## Key Naming Tips

Use descriptive names to avoid confusion:

- id_ed25519_personal — personal GitHub
- id_ed25519_work — work servers
- id_ed25519_homelab — homelab machines
- id_ed25519_backup — backup servers

---

## File Locations

All paths are under ~/.ssh/:

- **Private key** — ~/.ssh/{key_name}
- **Public key** — ~/.ssh/{key_name}.pub
- **SSH config** — ~/.ssh/config
- **Known host keys** — ~/.ssh/known_hosts
- **Remote authorized keys** — ~/.ssh/authorized_keys

---

## Security Notes

- **Never share** your private key
- Use **Ed25519** over RSA when possible
- Add a **passphrase** for sensitive keys
- known_hosts protects against man-in-the-middle attacks
- Use ssh-add to load keys into ssh-agent

---

## Troubleshooting

- **Permission denied** — Ensure private key is 0600
- **Connection refused** — Verify host, port, firewall
- **Host key mismatch** — Remove old entry from ~/.ssh/known_hosts
- **Config not updating** — Check write permissions on ~/.ssh/config
- **Terminal won't open** — Ensure a terminal emulator is installed and in PATH

---

## Dark Mode

Use the **Dark Mode** checkbox above to switch between light and dark themes.

`
	instructions := strings.ReplaceAll(rawInstructions, "«", "`")
	instructions = strings.ReplaceAll(instructions, "»", "`")

	richText := widget.NewRichTextFromMarkdown(instructions)
	richText.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(richText)
	scroll.SetMinSize(fyne.NewSize(900, 500))

	darkModeCheck := widget.NewCheck("Dark Mode", func(checked bool) {
		a.setDarkMode(checked)
	})

	return container.NewVBox(
		widget.NewLabel(titleInstructions),
		darkModeCheck,
		widget.NewLabel("How to use FyneSSH:"),
		scroll,
	)
}
