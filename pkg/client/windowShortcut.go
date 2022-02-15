package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Aserose/ModusOperandi/pkg/client/handler"
	"github.com/Aserose/ModusOperandi/pkg/config"
)

type windowShortcut struct {
	handler *handler.Handler
	validation
	notification
	cfg config.CfgClient
}

func NewWindowShortcut(handler *handler.Handler, v validation, n notification, cfg config.CfgClient) *windowShortcut {
	return &windowShortcut{
		handler:      handler,
		validation:   v,
		notification: n,
		cfg:          cfg,
	}
}

func (ws windowShortcut) getWindowShortcut() fyne.Window {
	windowShortcut := fyne.CurrentApp().NewWindow(ws.cfg.Shortcut.Name)
	windowShortcut.Resize(fyne.NewSize(ws.cfg.Shortcut.WindowWidth, ws.cfg.Shortcut.WindowHeight))

	return windowShortcut
}

func (ws windowShortcut) openWindowShortcut(instructionPaths []string) {
	w := ws.getWindowShortcut()

	entryForShortcutName := widget.NewEntry()
	entryForPath := widget.NewEntry()

	confirmButton := widget.NewButton("confirm", func() {
		ws.handler.CreateShortcut(instructionPaths, entryForShortcutName.Text, entryForPath.Text)
		w.Close()
	})
	cancelButton := widget.NewButton("cancel", func() { w.Close() })

	w.SetContent(container.NewGridWithRows(2,
		widget.NewForm(
			widget.NewFormItem("path", container.NewHSplit(
				entryForPath,
				widget.NewButtonWithIcon("", theme.FolderNewIcon(), func() { ws.openFileDialog(w, entryForPath) }))),
			widget.NewFormItem("shortcut name",
				entryForShortcutName)),

		container.NewGridWithColumns(2,
			confirmButton,
			cancelButton)))

	w.Show()
}

func (ws windowShortcut) openFileDialog(w fyne.Window, h *widget.Entry) {
	dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
		var path string
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if dir != nil {
			path = dir.Path()
		}
		h.SetText(path)
	}, w)
}
