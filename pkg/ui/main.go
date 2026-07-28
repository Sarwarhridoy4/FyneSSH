package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	app    fyne.App
}

// NewApp builds the main FyneSSH window.
func NewApp() *App {
	// TODO: persist app-specific metadata with app.NewWithID.
	a := app.New()
	w := a.NewWindow("FyneSSH")
	w.Resize(fyne.NewSize(1000, 800))
	return &App{window: w, app: a}
}

func (a *App) setDarkMode(enabled bool) {
	if enabled {
		a.app.Settings().SetTheme(theme.DarkTheme())
	} else {
		a.app.Settings().SetTheme(theme.LightTheme())
	}
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

