package ui

import (
	"context"
	"fmt"
	"io"
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
	titleLogin             = "Remote Login"
	labelUser              = "User:"
	labelHost              = "Host:"
	labelPort              = "Port:"
	placeholderHost        = "host or IP"
	placeholderPort        = "port"
	placeholderUser        = "username"
	placeholderPassword    = "password"
	statusNotConnected     = "Not connected"

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

	labelKeyName        = "Key name:"
	placeholderKeyName  = "e.g. id_ed25519_personal"
	labelCfgHostAlias   = "Host alias:"
	labelCfgHostName    = "HostName:"
	labelCfgUser        = "User:"
	labelCfgPort        = "Port:"
	labelCfgAddKeysToAgent   = "AddKeysToAgent:"
	labelCfgIdentitiesOnly   = "IdentitiesOnly:"
	labelCfgAliveInterval    = "ServerAliveInterval:"
	labelCfgAliveCountMax    = "ServerAliveCountMax:"
	btnSaveConfig      = "Save Key & Update Config"
	btnAddKnownHost    = "Add to known_hosts"
	btnUploadAndConfig = "Upload Key & Add Config"

	titleInstructions    = "Instructions"
	tabLogin             = "Login"
	tabKeys              = "Keys"
	tabInstructions      = "Instructions"
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
	w.Resize(fyne.NewSize(900, 700))
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

	terminal := widget.NewMultiLineEntry()
	terminal.Disable()

	connectBtn := widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
		host := strings.TrimSpace(hostEntry.Text)
		port := strings.TrimSpace(portEntry.Text)
		user := strings.TrimSpace(userEntry.Text)
		password := passEntry.Text

		if host == "" || port == "" || user == "" {
			status.SetText(msgMissingFields)
			return
		}

		var authMethods []ssh.AuthMethod
		if password != "" {
			authMethods = append(authMethods, ssh.Password(password))
		}

		hostKeyCallback, err := platform.HostKeyCallback()
		if err != nil {
			status.SetText(fmt.Sprintf("Host key callback error: %v", err))
			return
		}

		client, err := sshclient.DialWithHostKey(context.Background(), user, host, parsePort(port), authMethods, hostKeyCallback)
		if err != nil {
			status.SetText(err.Error())
			return
		}
		defer client.Close()

		status.SetText(fmt.Sprintf("Connected to %s@%s", user, host))

		terminal.SetText(func() string {
			var b strings.Builder
			_, _ = io.WriteString(&b, fmt.Sprintf("Connected to %s@%s\n", user, host))
			return b.String()
		}())
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
		widget.NewLabel(labelTerminal),
		terminal,
	)
}

// KeysUI encapsulates the widgets used in the Keys tab.
type KeysUI struct {
	algoSelect         *widget.Select
	comment            *widget.Entry
	passEntry          *widget.Entry
	keyName            *widget.Entry
	privPath           *widget.Entry
	pubPath            *widget.Entry
	privDisplay        *widget.Entry
	pubDisplay         *widget.Entry
	uploadHost         *widget.Entry
	uploadUser         *widget.Entry
	uploadPort         *widget.Entry
	uploadPass         *widget.Entry
	cfgHostAlias       *widget.Entry
	cfgHostName        *widget.Entry
	cfgUser            *widget.Entry
	cfgPort            *widget.Entry
	cfgAddKeysToAgent  *widget.Check
	cfgIdentitiesOnly  *widget.Check
	cfgAliveInterval   *widget.Entry
	cfgAliveCountMax   *widget.Entry
	status             *widget.Label
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
				Patterns:           []string{hostAlias},
				HostName:           hostName,
				User:               cfgUser,
				AddKeysToAgent:     ui.cfgAddKeysToAgent.Checked,
				IdentitiesOnly:     ui.cfgIdentitiesOnly.Checked,
				IdentityFile:       privPath,
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
		container.NewGridWithColumns(2,
			widget.NewLabel(labelAlgorithm), ui.algoSelect,
			widget.NewLabel(labelKeyName), ui.keyName,
			widget.NewLabel(labelComment), ui.comment,
			widget.NewLabel("Passphrase:"), ui.passEntry,
		),
		container.NewHBox(generateBtn, saveBtn, saveConfigBtn, copyBtn),
		ui.status,
		warning,
	}
	if rootWarning != nil {
		items = append(items, rootWarning)
	}
	items = append(items,
		container.NewGridWithColumns(2,
			widget.NewLabel(labelUploadUser), ui.uploadUser,
			widget.NewLabel(labelUploadHost), ui.uploadHost,
			widget.NewLabel(labelUploadPort), ui.uploadPort,
			widget.NewLabel(labelUploadPassword), ui.uploadPass,
		),
		container.NewHBox(uploadBtn, addKnownHostBtn),
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
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabel(labelPrivPath), ui.privPath),
			container.NewVBox(widget.NewLabel(labelPubPath), ui.pubPath),
		),
		container.NewGridWithColumns(2,
			widget.NewLabel(labelPrivKey),
			widget.NewLabel(labelPubKey),
		),
		container.NewGridWithColumns(2, ui.privDisplay, ui.pubDisplay),
	)

	return container.NewVBox(items...)
}

// Run starts the application.
func (a *App) Run() {
	tabs := container.NewAppTabs(
		container.NewTabItem(tabLogin, a.buildLoginForm()),
		container.NewTabItem(tabKeys, a.buildKeysTab()),
		container.NewTabItem(tabInstructions, a.buildInstructionsTab()),
	)
	// TODO: add File Manager, UFW, Port Config tabs.
	scroll := container.NewScroll(tabs)
	a.window.SetContent(scroll)
	a.window.ShowAndRun()
}

func (a *App) buildInstructionsTab() *fyne.Container {
	instructions := `FyneSSH - SSH Client

OVERVIEW
FyneSSH is a desktop SSH client that helps you manage SSH keys, connect to remote servers, and organize your SSH configuration.

TABS
1. Login - Connect to remote SSH servers
2. Keys - Generate and manage SSH key pairs
3. Instructions - This help page

LOGIN TAB
- Enter the remote host (hostname or IP address)
- Enter the port (default: 22)
- Enter your username
- Enter your password (optional if using key-based auth)
- Click "Connect" to establish an SSH connection
- The terminal area will show connection status

SSH KEYS TAB

Generate a Key Pair:
1. Select algorithm: Ed25519 (recommended) or RSA
2. Enter a key name/alias (e.g. id_ed25519_personal, id_ed25519_work)
   - This determines the filename: ~/.ssh/{key_name}
3. Optional: Add a comment
4. Optional: Add a passphrase for encryption
5. Click "Generate" to create the key pair
6. The private and public keys will appear in the display areas

Save Keys:
- Click "Save" to write keys to disk
- Default location: ~/.ssh/{key_name} and ~/.ssh/{key_name}.pub
- Private key permissions: 0600 (read/write for owner only)
- Public key permissions: 0644

Save Key & Update Config:
- Saves the key pair AND adds/updates an entry in ~/.ssh/config
- Fill in the SSH Config Entry section:
  * Host alias: Short name for the host (e.g. github, homelab)
  * HostName: Actual hostname or IP address
  * User: Remote username
  * Port: SSH port (default 22)
  * AddKeysToAgent: Automatically add key to ssh-agent (yes/no)
  * IdentitiesOnly: Only use specified identity files (yes/no)
  * ServerAliveInterval: Keep-alive interval in seconds
  * ServerAliveCountMax: Maximum keep-alive count
- If the Host alias already exists in config, it will be updated

Example SSH config entry:
Host github
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_personal
    AddKeysToAgent yes
    IdentitiesOnly yes

Upload Public Key:
- Enter the remote server details (host, user, port, password)
- Click "Upload to Server" to append your public key to ~/.ssh/authorized_keys on the remote server
- This enables passwordless SSH login

Add to known_hosts:
- Enter a hostname in the Upload Host field
- Click "Add to known_hosts" to record the server's host key
- This prevents man-in-the-middle attacks

Copy Public Key:
- Click "Copy Public Key" to copy the public key to clipboard
- Paste it into remote server's authorized_keys file manually

KEY NAMES / ALIASES
Use descriptive key names to avoid confusion:
- id_ed25519_personal - Personal GitHub account
- id_ed25519_work - Work servers
- id_ed25519_homelab - Homelab machines
- id_ed25519_backup - Backup servers

KEY FILE LOCATIONS
All keys are saved to: ~/.ssh/
- Private key: ~/.ssh/{key_name}
- Public key: ~/.ssh/{key_name}.pub
- SSH config: ~/.ssh/config
- Known hosts: ~/.ssh/known_hosts

SECURITY NOTES
- Never share your private key
- Use strong passphrases for sensitive keys
- Ed25519 is recommended over RSA for better security and performance
- known_hosts prevents man-in-the-middle attacks
- Run "ssh-add" to add keys to ssh-agent for passwordless auth

TROUBLESHOOTING
- Permission denied: Check key file permissions (private key should be 600)
- Connection refused: Verify host, port, and firewall settings
- Host key mismatch: Remove old entry from ~/.ssh/known_hosts
- Config not updating: Check write permissions on ~/.ssh/config
`

	richText := widget.NewRichTextFromMarkdown(instructions)
	richText.Wrapping = fyne.TextWrapWord

	scroll := container.NewScroll(richText)
	scroll.SetMinSize(fyne.NewSize(800, 500))

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
