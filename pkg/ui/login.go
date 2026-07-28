package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/Sarwarhridoy4/FyneSSH/internal/platform"
)

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

