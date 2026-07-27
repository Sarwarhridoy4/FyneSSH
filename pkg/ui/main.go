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

	tabLogin = "Login"
	tabKeys  = "Keys"
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
	w.Resize(fyne.NewSize(900, 600))
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

		client, err := sshclient.Dial(context.Background(), user, host, parsePort(port), authMethods)
		if err != nil {
			status.SetText(err.Error())
			return
		}

		status.SetText(fmt.Sprintf("Connected to %s@%s", user, host))

		terminal.SetText(func() string {
			var b strings.Builder
			_, _ = io.WriteString(&b, fmt.Sprintf("Connected to %s@%s\n", user, host))
			return b.String()
		}())

		_ = client
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
	algoSelect  *widget.Select
	comment     *widget.Entry
	passEntry   *widget.Entry
	privPath    *widget.Entry
	pubPath     *widget.Entry
	privDisplay *widget.Entry
	pubDisplay  *widget.Entry
	uploadHost  *widget.Entry
	uploadUser  *widget.Entry
	uploadPort  *widget.Entry
	uploadPass  *widget.Entry
	status      *widget.Label
}

func (a *App) buildKeysTab() *fyne.Container {
	ui := &KeysUI{}

	ui.algoSelect = widget.NewSelect([]string{string(keygen.AlgorithmEd25519), string(keygen.AlgorithmRSA)}, nil)
	ui.algoSelect.SetSelected(string(keygen.AlgorithmEd25519))

	ui.comment = widget.NewEntry()
	ui.comment.SetPlaceHolder(placeholderComment)

	ui.passEntry = widget.NewEntry()
	ui.passEntry.SetPlaceHolder(placeholderPassphrase)

	ui.privPath = widget.NewEntry()
	ui.privPath.SetText(platform.PrivateKeyPath(defaultPrivateName))
	ui.pubPath = widget.NewEntry()
	ui.pubPath.SetText(platform.PublicKeyPath(defaultPublicName))

	ui.uploadHost = widget.NewEntry()
	ui.uploadHost.SetPlaceHolder(placeholderUploadHost)

	ui.uploadUser = widget.NewEntry()
	ui.uploadUser.SetPlaceHolder(placeholderUploadUser)

	ui.uploadPort = widget.NewEntry()
	ui.uploadPort.SetPlaceHolder(placeholderUploadPort)
	ui.uploadPort.SetText("22")

	ui.uploadPass = widget.NewPasswordEntry()
	ui.uploadPass.SetPlaceHolder("server password")

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

	copyBtn := widget.NewButtonWithIcon("Copy Public Key", theme.ContentCopyIcon(), func() {
		pub := strings.TrimSpace(ui.pubDisplay.Text)
		if pub == "" {
			ui.status.SetText(errGenFirst)
			return
		}

		a.window.Clipboard().SetContent(pub)
		ui.status.SetText(msgCopied)
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

		ui.status.SetText(msgUploadSuccess)
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
			widget.NewLabel(labelComment), ui.comment,
			widget.NewLabel("Passphrase:"), ui.passEntry,
		),
		container.NewHBox(generateBtn, saveBtn, copyBtn, uploadBtn),
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
	)
	// TODO: add File Manager, UFW, Port Config tabs.
	scroll := container.NewScroll(tabs)
	a.window.SetContent(scroll)
	a.window.ShowAndRun()
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
