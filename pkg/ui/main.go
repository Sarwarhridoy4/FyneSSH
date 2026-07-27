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

	"github.com/Sarwarhridoy4/FyneSSH/internal/sshclient"
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
	hostEntry.SetPlaceHolder("host or IP")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("port")
	portEntry.SetText("22")

	userEntry := widget.NewEntry()
	userEntry.SetPlaceHolder("username")

	status := widget.NewLabel("Not connected")

	terminal := widget.NewMultiLineEntry()
	terminal.Disable()

	connectBtn := widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
		host := strings.TrimSpace(hostEntry.Text)
		port := strings.TrimSpace(portEntry.Text)
		user := strings.TrimSpace(userEntry.Text)

		if host == "" || port == "" || user == "" {
			status.SetText("Missing fields")
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
		widget.NewLabel("Remote Login"),
		container.NewHBox(
			widget.NewLabel("User:"), userEntry,
			widget.NewLabel("Host:"), hostEntry,
			widget.NewLabel("Port:"), portEntry,
			connectBtn,
		),
		status,
		widget.NewLabel("Terminal:"),
		terminal,
	)
}

// Run starts the application.
func (a *App) Run() {
	tabs := container.NewAppTabs(
		container.NewTabItem("Login", a.buildLoginForm()),
	)
	// TODO: add File Manager, Keys, UFW, Port Config tabs.
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
