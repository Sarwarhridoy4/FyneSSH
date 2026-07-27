package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.fcni.fyne-ssh")
	w := a.NewWindow("FyneSSH")
	w.Resize(fyne.NewSize(900, 600))

	entry := widget.NewEntry()
	entry.SetPlaceHolder("host or IP")
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("port")
	connectBtn := widget.NewButton("Connect", func() {})
	terminal := widget.NewMultiLineEntry()
	terminal.SetPlaceHolder("Terminal output will stream here...")

	c := container.NewVBox(
		widget.NewLabel("FyneSSH — placeholder window"),
		container.NewHBox(
			widget.NewLabel("Host:"),
			entry,
			widget.NewLabel("Port:"),
			portEntry,
			connectBtn,
		),
		terminal,
	)

	w.SetContent(c)
	w.ShowAndRun()
}
