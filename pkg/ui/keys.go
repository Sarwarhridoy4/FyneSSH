package ui

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Sarwarhridoy4/FyneSSH/internal/keygen"
	"github.com/Sarwarhridoy4/FyneSSH/internal/platform"
)

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

	generateBtn     *widget.Button
	saveBtn         *widget.Button
	saveConfigBtn   *widget.Button
	copyBtn         *widget.Button
	uploadBtn       *widget.Button
	openTerminalBtn *widget.Button
	addKnownHostBtn *widget.Button
}

func (ui *KeysUI) setEnabled(enabled bool) {
	buttons := []*widget.Button{
		ui.generateBtn,
		ui.saveBtn,
		ui.saveConfigBtn,
		ui.copyBtn,
		ui.uploadBtn,
		ui.openTerminalBtn,
		ui.addKnownHostBtn,
	}
	for _, btn := range buttons {
		if btn != nil {
			if enabled {
				btn.Enable()
			} else {
				btn.Disable()
			}
		}
	}
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
		ui.setEnabled(false)
		go func() {
			opts := keygen.Options{
				Algorithm:  keygen.Algorithm(ui.algoSelect.Selected),
				Comment:    strings.TrimSpace(ui.comment.Text),
				Passphrase: strings.TrimSpace(ui.passEntry.Text),
			}

			kp, err := keygen.Generate(opts)
			if err != nil {
				ui.status.SetText(fmt.Sprintf(errGenFailed, err))
				ui.setEnabled(true)
				return
			}

			ui.privDisplay.SetText(string(kp.PrivateKeyPEM))
			ui.pubDisplay.SetText(strings.TrimSpace(kp.PublicKeySSH))
			ui.status.SetText(msgGenerated)
			ui.setEnabled(true)
		}()
	})
	ui.generateBtn = generateBtn

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		ui.setEnabled(false)
		go func() {
			priv := strings.TrimSpace(ui.privDisplay.Text)
			pub := strings.TrimSpace(ui.pubDisplay.Text)
			if priv == "" || pub == "" {
				ui.status.SetText(errGenFirst)
				ui.setEnabled(true)
				return
			}

			privPath := ui.privPath.Text
			pubPath := ui.pubPath.Text
			if privPath == "" || pubPath == "" {
				ui.status.SetText(msgMissingFields)
				ui.setEnabled(true)
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
				ui.setEnabled(true)
				return
			}

			if _, err := os.Stat(privPath); err != nil {
				ui.status.SetText(fmt.Sprintf("Save reported success, but private key not found at %s: %v", privPath, err))
				ui.setEnabled(true)
				return
			}

			ui.status.SetText(fmt.Sprintf(msgSaved, privPath, pubPath))
			ui.setEnabled(true)
		}()
	})
	ui.saveBtn = saveBtn

	saveConfigBtn := widget.NewButtonWithIcon(btnSaveConfig, theme.DocumentSaveIcon(), func() {
		ui.setEnabled(false)
		go func() {
			priv := strings.TrimSpace(ui.privDisplay.Text)
			pub := strings.TrimSpace(ui.pubDisplay.Text)
			if priv == "" || pub == "" {
				ui.status.SetText(errGenFirst)
				ui.setEnabled(true)
				return
			}

			privPath := ui.privPath.Text
			pubPath := ui.pubPath.Text
			if privPath == "" || pubPath == "" {
				ui.status.SetText(msgMissingFields)
				ui.setEnabled(true)
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
				ui.setEnabled(true)
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
					ui.setEnabled(true)
					return
				}
				ui.status.SetText(fmt.Sprintf("Key saved and SSH config updated for %s", hostAlias))
				ui.setEnabled(true)
				return
			}

			ui.status.SetText(fmt.Sprintf(msgSaved, privPath, pubPath))
			ui.setEnabled(true)
		}()
	})
	ui.saveConfigBtn = saveConfigBtn

	uploadBtn := widget.NewButtonWithIcon("Upload to Server", theme.UploadIcon(), func() {
		ui.setEnabled(false)
		go func() {
			pubPath := strings.TrimSpace(ui.pubPath.Text)
			host := strings.TrimSpace(ui.uploadHost.Text)
			user := strings.TrimSpace(ui.uploadUser.Text)
			port := strings.TrimSpace(ui.uploadPort.Text)
			password := ui.uploadPass.Text

			if pubPath == "" || host == "" || user == "" || port == "" || password == "" {
				ui.status.SetText(errUploadMissing)
				ui.setEnabled(true)
				return
			}

			if err := uploadPublicKey(pubPath, user, host, port, password); err != nil {
				ui.status.SetText(fmt.Sprintf(errUploadFailed, err))
				ui.setEnabled(true)
				return
			}

			if err := addToKnownHostsUI(host); err != nil {
				ui.status.SetText(fmt.Sprintf("Uploaded, but failed to add to known_hosts: %v", err))
				ui.setEnabled(true)
				return
			}

			ui.status.SetText(msgUploadSuccess)
			ui.setEnabled(true)
		}()
	})
	ui.uploadBtn = uploadBtn

	openTerminalBtn := widget.NewButtonWithIcon(btnOpenTerminal, theme.ComputerIcon(), func() {
		ui.setEnabled(false)
		go func() {
			hostAlias := strings.TrimSpace(ui.cfgHostAlias.Text)
			if hostAlias == "" {
				ui.status.SetText("Host alias is required")
				ui.setEnabled(true)
				return
			}

			if err := platform.OpenTerminal(fmt.Sprintf("ssh %s", hostAlias)); err != nil {
				ui.status.SetText(fmt.Sprintf("Failed to open terminal: %v", err))
				ui.setEnabled(true)
				return
			}

			ui.status.SetText(fmt.Sprintf(msgTerminalOpened, hostAlias))
			ui.setEnabled(true)
		}()
	})
	ui.openTerminalBtn = openTerminalBtn

	copyBtn := widget.NewButtonWithIcon("Copy Public Key", theme.ContentCopyIcon(), func() {
		pub := strings.TrimSpace(ui.pubDisplay.Text)
		if pub == "" {
			ui.status.SetText(errGenFirst)
			return
		}

		a.window.Clipboard().SetContent(pub)
		ui.status.SetText(msgCopied)
	})
	ui.copyBtn = copyBtn

	addKnownHostBtn := widget.NewButtonWithIcon(btnAddKnownHost, theme.ComputerIcon(), func() {
		ui.setEnabled(false)
		go func() {
			host := strings.TrimSpace(ui.uploadHost.Text)
			if host == "" {
				ui.status.SetText("Host is required")
				ui.setEnabled(true)
				return
			}
			if err := addToKnownHostsUI(host); err != nil {
				ui.status.SetText(fmt.Sprintf("Failed to add to known_hosts: %v", err))
				ui.setEnabled(true)
				return
			}
			ui.status.SetText(fmt.Sprintf(msgKnownHostsAdded, host))
			ui.setEnabled(true)
		}()
	})
	ui.addKnownHostBtn = addKnownHostBtn

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

