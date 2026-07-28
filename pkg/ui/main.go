package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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

