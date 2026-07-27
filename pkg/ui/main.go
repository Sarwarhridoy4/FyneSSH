package ui

import (
	"context"
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Sarwarhridoy4/FyneSSH/internal/keygen"
	"github.com/Sarwarhridoy4/FyneSSH/internal/sshclient"
)

const (
	titleLogin         = "Remote Login"
	labelUser          = "User:"
	labelHost          = "Host:"
	labelPort          = "Port:"
	placeholderHost    = "host or IP"
	placeholderPort    = "port"
	placeholderUser    = "username"
	statusNotConnected = "Not connected"

	titleKeys             = "SSH Keys"
	labelAlgorithm        = "Algorithm:"
	labelComment          = "Comment:"
	labelPrivPath         = "Private key path:"
	labelPubPath          = "Public key path:"
	placeholderComment    = "optional key comment"
	defaultPrivateName    = "id_ed25519"
	defaultPublicName     = "id_ed25519.pub"
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
	msgSaved              = "Keys saved successfully."
	msgCopied             = "Public key copied to clipboard."
	errGenFailed          = "Key generation failed: %v"
	errSaveFailed         = "Save failed: %v"
	errClipFailed         = "Copy failed: %v"

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

	status := widget.NewLabel(statusNotConnected)

	terminal := widget.NewMultiLineEntry()
	terminal.Disable()

	connectBtn := widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
		host := strings.TrimSpace(hostEntry.Text)
		port := strings.TrimSpace(portEntry.Text)
		user := strings.TrimSpace(userEntry.Text)

		if host == "" || port == "" || user == "" {
			status.SetText(msgMissingFields)
			return
		}

		client, err := sshclient.Dial(context.Background(), user, host, defaultPort(port), nil)
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
	privPath    *widget.Entry
	pubPath     *widget.Entry
	privDisplay *widget.Entry
	pubDisplay  *widget.Entry
	status      *widget.Label
}

func (a *App) buildKeysTab() *fyne.Container {
	ui := &KeysUI{}

	ui.algoSelect = widget.NewSelect([]string{string(keygen.AlgorithmEd25519), string(keygen.AlgorithmRSA)}, nil)
	ui.algoSelect.SetSelected(string(keygen.AlgorithmEd25519))

	ui.comment = widget.NewEntry()
	ui.comment.SetPlaceHolder(placeholderComment)

	ui.privPath = widget.NewEntry()
	ui.privPath.SetText(defaultPrivateName)
	ui.pubPath = widget.NewEntry()
	ui.pubPath.SetText(defaultPublicName)

	ui.privDisplay = widget.NewMultiLineEntry()
	ui.privDisplay.Disable()
	ui.privDisplay.SetPlaceHolder(privateKeyPlaceholder)

	ui.pubDisplay = widget.NewMultiLineEntry()
	ui.pubDisplay.Disable()
	ui.pubDisplay.SetPlaceHolder(publicKeyPlaceholder)

	ui.status = widget.NewLabel("")

	generateBtn := widget.NewButtonWithIcon("Generate", theme.ContentAddIcon(), func() {
		opts := keygen.Options{
			Algorithm: keygen.Algorithm(ui.algoSelect.Selected),
			Comment:   strings.TrimSpace(ui.comment.Text),
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
			ui.status.SetText(fmt.Sprintf(errSaveFailed, err))
			return
		}

		ui.status.SetText(msgSaved)
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

	warning := widget.NewLabel(warningNoShare)
	warning.TextStyle = fyne.TextStyle{Bold: true}

	return container.NewVBox(
		widget.NewLabel(titleKeys),
		widget.NewLabel(titleKeysSub),
		container.NewGridWithColumns(2,
			widget.NewLabel(labelAlgorithm), ui.algoSelect,
			widget.NewLabel(labelComment), ui.comment,
		),
		container.NewHBox(generateBtn, saveBtn, copyBtn),
		ui.status,
		warning,
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
}

// Run starts the application.
func (a *App) Run() {
	tabs := container.NewAppTabs(
		container.NewTabItem(tabLogin, a.buildLoginForm()),
		container.NewTabItem(tabKeys, a.buildKeysTab()),
	)
	// TODO: add File Manager, UFW, Port Config tabs.
	a.window.SetContent(tabs)
	a.window.ShowAndRun()
}

func defaultPort(s string) int {
	switch s {
	case "", "22":
		return 22
	}
	var p int
	_, err := fmt.Sscan(s, &p)
	if err != nil {
		return 22
	}
	return p
}
