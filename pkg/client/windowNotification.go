package client

import (
	"fyne.io/fyne/v2"
	"github.com/Aserose/ModusOperandi/pkg/config"
)

type windowNotification struct {
	cfg config.CfgClient
}

func NewNotificationWindow(cfg config.CfgClient) *windowNotification {
	return &windowNotification{
		cfg: cfg,
	}
}

func (n windowNotification) createNotification(content fyne.CanvasObject) {
	w := fyne.CurrentApp().NewWindow(n.cfg.Notification.Name)
	w.CenterOnScreen()
	w.SetContent(content)
	w.Show()
}
