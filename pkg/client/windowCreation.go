package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/Aserose/ModusOperandi/pkg/config"
)

type windowInstructionCreation struct {
	cfg config.CfgClient
}

func NewInstructionCreation(cfg config.CfgClient) *windowInstructionCreation {
	return &windowInstructionCreation{
		cfg: cfg,
	}
}

func (wc windowInstructionCreation) getWindowCreation() fyne.Window {
	creationWindow := fyne.CurrentApp().NewWindow(wc.cfg.Creation.Name)
	creationWindow.SetIcon(theme.ContentAddIcon())
	creationWindow.Resize(fyne.NewSize(wc.cfg.Creation.WindowWidth, wc.cfg.Creation.WindowHeight))

	return creationWindow
}
