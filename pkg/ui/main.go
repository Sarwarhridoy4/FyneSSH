package ui

import (
	"context"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"golang.org/x/crypto/ssh"

	"github.com/Sarwarhridoy4/FyneSSH/internal/keygen"
	"github.com/Sarwarhridoy4/FyneSSH/internal/platform"
	"github.com/Sarwarhridoy4/FyneSSH/internal/sshclient"
)

const (
	titleLogin          = "Remote Login"
	labelUser           = "User:"
	labelHost           = "Host:"
	labelPort           = "Port:"
	placeholderHost     = "host or IP"
	placeholderPort     = "port"
	placeholderUser     = "username"
	placeholderPassword = "password"
	statusNotConnected  = "Not connected"

	titleKeys             = "SSH Keys"
	labelAlgorithm        = "Algorithm:"
	labelComment          = "Comment:"
	labelPrivPath         = "Private key path:"
	labelPubPath          = "Public key path:"
	placeholderComment    = "optional key comment"
	placeholderPassphrase = "optional passphrase"
	defaultPrivateName    = "id_ed25519"
	defaultPublicName     = "id_ed25519"
	defaultComment        = "FyneSSH generated key"
	labelTerminal         = "Terminal:"
	titleKeysSub          = "Generate an SSH key pair for authentication."
	labelPrivKey          = "Private key:"
	labelPubKey           = "Public key (authorized_keys format):"
	errGenFirst           = "Generate a key pair first"
	warningNoShare        = "Security warning: private key must never be shared."
	privateKeyPlaceholder = "Private key will appear here (never share this)"
	publicKeyPlaceholder  = "Public key will appear here"
	msgMissingFields      = "Missing fields"
	msgGenerated          = "Key pair generated. Save to disk or copy public key."
	msgCopied             = "Public key copied to clipboard."
	errGenFailed          = "Key generation failed: %v"
	errSaveFailed         = "Save failed: %v"
	errClipFailed         = "Copy failed: %v"
	errUploadFailed       = "Upload failed: %v"
	errUploadMissing      = "Public key path, user, host, port, and password are required"
	msgUploadSuccess      = "Public key uploaded successfully."
	labelUploadHost       = "Upload Host:"
	labelUploadUser       = "Upload User:"
	labelUploadPort       = "Port:"
	labelUploadPassword   = "Server Password:"
	placeholderUploadHost = "server IP or hostname"
	placeholderUploadUser = "server username"
	placeholderUploadPort = "22"
	msgSaved              = "Keys saved successfully to %s and %s"
	msgSaveFailed         = "Save failed: %v"
	msgConfigUpdated      = "SSH config updated for host %s"
	msgKnownHostsAdded    = "Host %s added to known_hosts"
	msgHostKeyMismatch    = "Host key mismatch for %s"
	msgTerminalOpened     = "Opening terminal with ssh %s..."

	labelKeyName           = "Key name:"
	placeholderKeyName     = "e.g. id_ed25519_personal"
	labelCfgHostAlias      = "Host alias:"
	labelCfgHostName       = "HostName:"
	labelCfgUser           = "User:"
	labelCfgPort           = "Port:"
	labelCfgAddKeysToAgent = "AddKeysToAgent:"
	labelCfgIdentitiesOnly = "IdentitiesOnly:"
	labelCfgAliveInterval  = "ServerAliveInterval:"
	labelCfgAliveCountMax  = "ServerAliveCountMax:"
	btnSaveConfig          = "Save Key & Update Config"
	btnAddKnownHost        = "Add to known_hosts"
	btnUploadAndConfig     = "Upload Key & Add Config"
	btnOpenTerminal        = "Open Terminal"

	titleInstructions = "Instructions"
	tabLogin          = "Login"
	tabKeys           = "Keys"
	tabInstructions   = "Instructions"
)

// App holds references to the main window and backend services.
type App struct {
	window fyne.Window
}

// NewApp builds the main FyneSSH window.
func NewApp() *App {
	// TODO: persist app-specific metadata with app.NewWithID.
	a := app.New()
	w := a.NewWindow("FyneSSH")
	w.Resize(fyne.NewSize(1000, 800))
	return &App{window: w}
}

func (a *App) buildLoginForm() *fyne.Container {
	hostEntry := widget.NewEntry()
	hostEntry.SetPlaceHolder(placeholderHost)

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder(placeholderPort)
	portEntry.SetText("22")

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder(placeholderUser)

	passEntry := widget.NewPasswordEntry()
	passEntry.SetPlaceHolder(placeholderPassword)

	status := widget.NewLabel(statusNotConnected)

	connectBtn := widget.NewButtonWithIcon("Open Terminal", theme.ComputerIcon(), func() {
		host := strings.TrimSpace(hostEntry.Text)
		port := strings.TrimSpace(portEntry.Text)
		user := strings.TrimSpace(userEntry.Text)

		if host == "" || port == "" || user == "" {
			status.SetText(msgMissingFields)
			return
		}

		target := fmt.Sprintf("%s@%s", user, host)
		if port != "22" {
			target = fmt.Sprintf("%s@%s -p %s", user, host, port)
		}

		if err := platform.OpenTerminal(fmt.Sprintf("ssh %s", target)); err != nil {
			status.SetText(fmt.Sprintf("Failed to open terminal: %v", err))
			return
		}

		status.SetText(fmt.Sprintf(msgTerminalOpened, target))
	})

	return container.NewVBox(
		widget.NewLabel(titleLogin),
		container.NewGridWithColumns(2,
			widget.NewLabel(labelUser), userEntry,
			widget.NewLabel(labelHost), hostEntry,
			widget.NewLabel(labelPort), portEntry,
			widget.NewLabel("Password:"), passEntry,
		),
		connectBtn,
		status,
	)
}

// KeysUI encapsulates the widgets used in the Keys tab.
type KeysUI struct {
	algoSelect        *widget.Select
	comment           *widget.Entry
	passEntry         *widget.Entry
	keyName           *widget.Entry
	privPath          *widget.Entry
	pubPath           *widget.Entry
	privDisplay       *widget.Entry
	pubDisplay        *widget.Entry
	uploadHost        *widget.Entry
	uploadUser        *widget.Entry
	uploadPort        *widget.Entry
	uploadPass        *widget.Entry
	cfgHostAlias      *widget.Entry
	cfgHostName       *widget.Entry
	cfgUser           *widget.Entry
	cfgPort           *widget.Entry
	cfgAddKeysToAgent *widget.Check
	cfgIdentitiesOnly *widget.Check
	cfgAliveInterval  *widget.Entry
	cfgAliveCountMax  *widget.Entry
	status            *widget.Label
}

func (a *App) buildKeysTab() *fyne.Container {
	ui := &KeysUI{}

	ui.algoSelect = widget.NewSelect([]string{string(keygen.AlgorithmEd25519), string(keygen.AlgorithmRSA)}, nil)
	ui.algoSelect.SetSelected(string(keygen.AlgorithmEd25519))

	ui.comment = widget.NewEntry()
	ui.comment.SetPlaceHolder(placeholderComment)

	ui.passEntry = widget.NewEntry()
	ui.passEntry.SetPlaceHolder(placeholderPassphrase)

	ui.keyName = widget.NewEntry()
	ui.keyName.SetText(defaultPrivateName)
	ui.keyName.SetPlaceHolder(placeholderKeyName)

	ui.privPath = widget.NewEntry()
	ui.privPath.SetText(platform.PrivateKeyPath(defaultPrivateName))
	ui.pubPath = widget.NewEntry()
	ui.pubPath.SetText(platform.PublicKeyPath(defaultPublicName))

	ui.keyName.OnChanged = func(name string) {
		name = strings.TrimSpace(name)
		if name == "" {
			name = defaultPrivateName
		}
		ui.privPath.SetText(platform.PrivateKeyPath(name))
		ui.pubPath.SetText(platform.PublicKeyPath(name))
		ui.cfgHostAlias.SetText(name)
		ui.cfgHostName.SetText("")
	}

	ui.uploadHost = widget.NewEntry()
	ui.uploadHost.SetPlaceHolder(placeholderUploadHost)

	ui.uploadUser = widget.NewEntry()
	ui.uploadUser.SetPlaceHolder(placeholderUploadUser)

	ui.uploadPort = widget.NewEntry()
	ui.uploadPort.SetPlaceHolder(placeholderUploadPort)
	ui.uploadPort.SetText("22")

	ui.uploadPass = widget.NewPasswordEntry()
	ui.uploadPass.SetPlaceHolder("server password")

	ui.cfgHostAlias = widget.NewEntry()
	ui.cfgHostAlias.SetPlaceHolder("host alias (e.g. github)")

	ui.cfgHostName = widget.NewEntry()
	ui.cfgHostName.SetPlaceHolder("hostname or IP")

	ui.cfgUser = widget.NewEntry()
	ui.cfgUser.SetPlaceHolder("username")

	ui.cfgPort = widget.NewEntry()
	ui.cfgPort.SetPlaceHolder("port")
	ui.cfgPort.SetText("22")

	ui.cfgAddKeysToAgent = widget.NewCheck("AddKeysToAgent", nil)
	ui.cfgAddKeysToAgent.SetChecked(true)

	ui.cfgIdentitiesOnly = widget.NewCheck("IdentitiesOnly", nil)
	ui.cfgIdentitiesOnly.SetChecked(true)

	ui.cfgAliveInterval = widget.NewEntry()
	ui.cfgAliveInterval.SetPlaceHolder("seconds (optional)")

	ui.cfgAliveCountMax = widget.NewEntry()
	ui.cfgAliveCountMax.SetPlaceHolder("max count (optional)")

	ui.privDisplay = widget.NewMultiLineEntry()
	ui.privDisplay.SetPlaceHolder(privateKeyPlaceholder)

	ui.pubDisplay = widget.NewMultiLineEntry()
	ui.pubDisplay.SetPlaceHolder(publicKeyPlaceholder)

	ui.status = widget.NewLabel("")

	generateBtn := widget.NewButtonWithIcon("Generate", theme.ContentAddIcon(), func() {
		opts := keygen.Options{
			Algorithm:  keygen.Algorithm(ui.algoSelect.Selected),
			Comment:    strings.TrimSpace(ui.comment.Text),
			Passphrase: strings.TrimSpace(ui.passEntry.Text),
		}

		kp, err := keygen.Generate(opts)
		if err != nil {
			ui.status.SetText(fmt.Sprintf(errGenFailed, err))
			return
		}

		ui.privDisplay.SetText(string(kp.PrivateKeyPEM))
		ui.pubDisplay.SetText(strings.TrimSpace(kp.PublicKeySSH))
		ui.status.SetText(msgGenerated)
	})

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		priv := strings.TrimSpace(ui.privDisplay.Text)
		pub := strings.TrimSpace(ui.pubDisplay.Text)
		if priv == "" || pub == "" {
			ui.status.SetText(errGenFirst)
			return
		}

		privPath := ui.privPath.Text
		pubPath := ui.pubPath.Text
		if privPath == "" || pubPath == "" {
			ui.status.SetText(msgMissingFields)
			return
		}

		kp := &keygen.KeyPair{
			PrivateKeyPEM: []byte(priv),
			PublicKeySSH:  pub + "\n",
			Algorithm:     keygen.Algorithm(ui.algoSelect.Selected),
			Comment:       strings.TrimSpace(ui.comment.Text),
		}

		if err := kp.Save(privPath, pubPath); err != nil {
			ui.status.SetText(fmt.Sprintf(msgSaveFailed, err))
			return
		}

		if _, err := os.Stat(privPath); err != nil {
			ui.status.SetText(fmt.Sprintf("Save reported success, but private key not found at %s: %v", privPath, err))
			return
		}

		ui.status.SetText(fmt.Sprintf(msgSaved, privPath, pubPath))
	})

	saveConfigBtn := widget.NewButtonWithIcon(btnSaveConfig, theme.DocumentSaveIcon(), func() {
		priv := strings.TrimSpace(ui.privDisplay.Text)
		pub := strings.TrimSpace(ui.pubDisplay.Text)
		if priv == "" || pub == "" {
			ui.status.SetText(errGenFirst)
			return
		}

		privPath := ui.privPath.Text
		pubPath := ui.pubPath.Text
		if privPath == "" || pubPath == "" {
			ui.status.SetText(msgMissingFields)
			return
		}

		kp := &keygen.KeyPair{
			PrivateKeyPEM: []byte(priv),
			PublicKeySSH:  pub + "\n",
			Algorithm:     keygen.Algorithm(ui.algoSelect.Selected),
			Comment:       strings.TrimSpace(ui.comment.Text),
		}

		if err := kp.Save(privPath, pubPath); err != nil {
			ui.status.SetText(fmt.Sprintf(msgSaveFailed, err))
			return
		}

		hostAlias := strings.TrimSpace(ui.cfgHostAlias.Text)
		hostName := strings.TrimSpace(ui.cfgHostName.Text)
		cfgUser := strings.TrimSpace(ui.cfgUser.Text)
		cfgPort := strings.TrimSpace(ui.cfgPort.Text)
		if hostAlias != "" && hostName != "" && cfgUser != "" {
			block := platform.HostBlock{
				Patterns:       []string{hostAlias},
				HostName:       hostName,
				User:           cfgUser,
				AddKeysToAgent: ui.cfgAddKeysToAgent.Checked,
				IdentitiesOnly: ui.cfgIdentitiesOnly.Checked,
				IdentityFile:   privPath,
			}
			if cfgPort != "" && cfgPort != "0" {
				var p int
				if _, err := fmt.Sscan(cfgPort, &p); err == nil {
					block.Port = p
				}
			}
			aliveInterval := strings.TrimSpace(ui.cfgAliveInterval.Text)
			if aliveInterval != "" {
				var v int
				if _, err := fmt.Sscan(aliveInterval, &v); err == nil {
					block.ServerAliveInterval = v
				}
			}
			aliveCountMax := strings.TrimSpace(ui.cfgAliveCountMax.Text)
			if aliveCountMax != "" {
				var v int
				if _, err := fmt.Sscan(aliveCountMax, &v); err == nil {
					block.ServerAliveCountMax = v
				}
			}
			if err := platform.UpdateOrAddHost(platform.ConfigPath(), block); err != nil {
				ui.status.SetText(fmt.Sprintf("Key saved, but config update failed: %v", err))
				return
			}
			ui.status.SetText(fmt.Sprintf("Key saved and SSH config updated for %s", hostAlias))
			return
		}

		ui.status.SetText(fmt.Sprintf(msgSaved, privPath, pubPath))
	})

	uploadBtn := widget.NewButtonWithIcon("Upload to Server", theme.UploadIcon(), func() {
		pubPath := strings.TrimSpace(ui.pubPath.Text)
		host := strings.TrimSpace(ui.uploadHost.Text)
		user := strings.TrimSpace(ui.uploadUser.Text)
		port := strings.TrimSpace(ui.uploadPort.Text)
		password := ui.uploadPass.Text

		if pubPath == "" || host == "" || user == "" || port == "" || password == "" {
			ui.status.SetText(errUploadMissing)
			return
		}

		if err := uploadPublicKey(pubPath, user, host, port, password); err != nil {
			ui.status.SetText(fmt.Sprintf(errUploadFailed, err))
			return
		}

		if err := addToKnownHostsUI(host); err != nil {
			ui.status.SetText(fmt.Sprintf("Uploaded, but failed to add to known_hosts: %v", err))
			return
		}

		ui.status.SetText(msgUploadSuccess)
	})

	openTerminalBtn := widget.NewButtonWithIcon(btnOpenTerminal, theme.ComputerIcon(), func() {
		hostAlias := strings.TrimSpace(ui.cfgHostAlias.Text)
		if hostAlias == "" {
			ui.status.SetText("Host alias is required")
			return
		}

		if err := platform.OpenTerminal(fmt.Sprintf("ssh %s", hostAlias)); err != nil {
			ui.status.SetText(fmt.Sprintf("Failed to open terminal: %v", err))
			return
		}

		ui.status.SetText(fmt.Sprintf(msgTerminalOpened, hostAlias))
	})

	copyBtn := widget.NewButtonWithIcon("Copy Public Key", theme.ContentCopyIcon(), func() {
		pub := strings.TrimSpace(ui.pubDisplay.Text)
		if pub == "" {
			ui.status.SetText(errGenFirst)
			return
		}

		a.window.Clipboard().SetContent(pub)
		ui.status.SetText(msgCopied)
	})

	addKnownHostBtn := widget.NewButtonWithIcon(btnAddKnownHost, theme.ComputerIcon(), func() {
		host := strings.TrimSpace(ui.uploadHost.Text)
		if host == "" {
			ui.status.SetText("Host is required")
			return
		}
		if err := addToKnownHostsUI(host); err != nil {
			ui.status.SetText(fmt.Sprintf("Failed to add to known_hosts: %v", err))
			return
		}
		ui.status.SetText(fmt.Sprintf(msgKnownHostsAdded, host))
	})

	warning := widget.NewLabel(warningNoShare)
	warning.TextStyle = fyne.TextStyle{Bold: true}

	var rootWarning fyne.CanvasObject
	if platform.IsRootUser() {
		rootWarning = widget.NewLabel("Warning: running as root. Keys will be saved to /root/.ssh/")
		rootWarning.(*widget.Label).TextStyle = fyne.TextStyle{Bold: true}
	}

	items := []fyne.CanvasObject{
		widget.NewLabel(titleKeys),
		widget.NewLabel(titleKeysSub),
		widget.NewLabel("Generate Key Pair:"),
		container.NewGridWithColumns(2,
			widget.NewLabel(labelAlgorithm), ui.algoSelect,
			widget.NewLabel(labelKeyName), ui.keyName,
			widget.NewLabel(labelComment), ui.comment,
			widget.NewLabel("Passphrase:"), ui.passEntry,
		),
		generateBtn,
		widget.NewLabel("Generated Keys:"),
		container.NewGridWithColumns(2,
			widget.NewLabel(labelPrivKey),
			widget.NewLabel(labelPubKey),
		),
		container.NewGridWithColumns(2, ui.privDisplay, ui.pubDisplay),
		widget.NewLabel("Save & Config:"),
		container.NewHBox(saveBtn, saveConfigBtn, copyBtn),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabel(labelPrivPath), ui.privPath),
			container.NewVBox(widget.NewLabel(labelPubPath), ui.pubPath),
		),
		widget.NewLabel("SSH Config Entry:"),
		container.NewGridWithColumns(2,
			widget.NewLabel(labelCfgHostAlias), ui.cfgHostAlias,
			widget.NewLabel(labelCfgHostName), ui.cfgHostName,
			widget.NewLabel(labelCfgUser), ui.cfgUser,
			widget.NewLabel(labelCfgPort), ui.cfgPort,
			widget.NewLabel(labelCfgAddKeysToAgent), ui.cfgAddKeysToAgent,
			widget.NewLabel(labelCfgIdentitiesOnly), ui.cfgIdentitiesOnly,
			widget.NewLabel(labelCfgAliveInterval), ui.cfgAliveInterval,
			widget.NewLabel(labelCfgAliveCountMax), ui.cfgAliveCountMax,
		),
		openTerminalBtn,
		widget.NewLabel("Upload to Server:"),
		container.NewGridWithColumns(2,
			widget.NewLabel(labelUploadUser), ui.uploadUser,
			widget.NewLabel(labelUploadHost), ui.uploadHost,
			widget.NewLabel(labelUploadPort), ui.uploadPort,
			widget.NewLabel(labelUploadPassword), ui.uploadPass,
		),
		container.NewHBox(uploadBtn, addKnownHostBtn),
		ui.status,
		warning,
	}
	if rootWarning != nil {
		items = append(items, rootWarning)
	}

	return container.NewVBox(items...)
}

// Run starts the application.
func (a *App) Run() {
	tabs := container.NewAppTabs(
		container.NewTabItem(tabLogin, container.NewScroll(a.buildLoginForm())),
		container.NewTabItem(tabKeys, container.NewScroll(a.buildKeysTab())),
		container.NewTabItem(tabInstructions, container.NewScroll(a.buildInstructionsTab())),
	)
	// TODO: add File Manager, UFW, Port Config tabs.
	a.window.SetContent(tabs)
	a.window.ShowAndRun()
}

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
`
	instructions := strings.ReplaceAll(rawInstructions, "«", "`")
	instructions = strings.ReplaceAll(instructions, "»", "`")

	richText := widget.NewRichTextFromMarkdown(instructions)
	richText.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(richText)
	scroll.SetMinSize(fyne.NewSize(900, 500))

	return container.NewVBox(
		widget.NewLabel(titleInstructions),
		widget.NewLabel("How to use FyneSSH:"),
		scroll,
	)
}
func uploadPublicKey(pubKeyPath, user, host, port, password string) error {
	portNum := parsePort(port)

	authMethods := []ssh.AuthMethod{ssh.Password(password)}

	client, err := sshclient.Dial(context.Background(), user, host, portNum, authMethods)
	if err != nil {
		return fmt.Errorf("connect to %s@%s:%s: %w", user, host, port, err)
	}
	defer client.Close()

	session, err := client.Session()
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	defer session.Close()

	pubKeyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return fmt.Errorf("read public key file: %w", err)
	}

	pubKeyContent := strings.TrimSpace(string(pubKeyBytes))
	escapedKey := strings.ReplaceAll(pubKeyContent, "'", "'\\''")

	shellCmd := fmt.Sprintf("mkdir -p ~/.ssh && chmod 700 ~/.ssh && echo '%s' >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys", escapedKey)

	if err := session.Run(shellCmd); err != nil {
		return fmt.Errorf("append public key: %w", err)
	}

	return nil
}

func addToKnownHostsUI(host string) error {
	portNum := 22
	authMethods := []ssh.AuthMethod{}

	hostKeyCallback, err := platform.HostKeyCallback()
	if err != nil {
		return fmt.Errorf("create host key callback: %w", err)
	}

	client, err := sshclient.DialWithHostKey(context.Background(), "root", host, portNum, authMethods, hostKeyCallback)
	if err != nil {
		return fmt.Errorf("connect to %s: %w", host, err)
	}
	defer client.Close()

	return nil
}

func parsePort(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 22
	}
	var p int
	_, err := fmt.Sscan(s, &p)
	if err != nil {
		return 22
	}
	return p
}
