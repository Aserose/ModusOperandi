package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/Aserose/ModusOperandi/pkg/client/handler"
	"github.com/Aserose/ModusOperandi/pkg/config"
)

type Client struct {
	a fyne.App
	w fyne.Window
	*windowMain
}

func NewClient(handler *handler.Handler, cfg config.CfgClient) *Client {
	a := app.New()
	w := a.NewWindow(cfg.Name)
	w.Resize(fyne.NewSize(cfg.Main.WindowHeight, cfg.Main.WindowWidth))
	w.SetIcon(theme.FolderNewIcon())
	w.CenterOnScreen()

	return &Client{
		a:          a,
		w:          w,
		windowMain: NewMainWindow(w, handler, cfg),
	}
}

func (c Client) Start() {
	c.showInstructions()
	c.w.ShowAndRun()
}
