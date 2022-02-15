package client

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Aserose/ModusOperandi/pkg/client/handler"
	"os"
)

type validator struct {
	notificationLabel *widget.Label
	n                 notification
	handler           *handler.Handler
}

func NewValidator(handler *handler.Handler, n notification) *validator {
	notific := widget.NewLabel("")
	notific.TextStyle.Bold = true

	return &validator{
		notificationLabel: notific,
		n:                 n,
		handler:           handler,
	}
}

func (v validator) removePathQuotes(path []rune) string {
	return string(path[1 : len(path)-1])
}

func (v validator) pathDoesNotExist(path string) bool {
	if _, err := os.Stat(v.removePathQuotes([]rune(path))); os.IsNotExist(err) {
		return true
	}
	return false
}

func (v validator) instructionNameIsBlank(instructionName string) bool {
	if instructionName == "" {
		v.notificationLabel.SetText("the name of the instruction must not be empty")
		v.n.createNotification(container.NewCenter(v.notificationLabel))

		return true
	}
	return false
}

func (v validator) instructionNameExist(instructionName string) bool {
	if v.handler.InstructionsExist(instructionName) {
		v.notificationLabel.SetText("this name is occupied by another instruction")
		v.n.createNotification(container.NewCenter(v.notificationLabel))

		return true
	}
	return false
}

func (v validator) pathAlreadyExist(paths []string, input string) bool {
	for _, path := range paths {
		if path == input {
			v.notificationLabel.SetText("this path is already in the instruction")
			v.n.createNotification(container.NewCenter(v.notificationLabel))

			return true
		}
	}
	return false
}

func (v validator) pathAlreadyAdded(paths []string, input string) bool {
	for _, path := range paths {
		if path == input {
			v.notificationLabel.SetText("this path is already in the addlist")
			v.n.createNotification(container.NewCenter(v.notificationLabel))

			return true
		}
	}
	return false
}

func (v validator) pathDoesNotValid(path string) bool {
	if _, err := os.Stat(v.removePathQuotes([]rune(path))); os.IsNotExist(err) {
		v.notificationLabel.SetText("file not found")
		v.n.createNotification(container.NewCenter(v.notificationLabel))

		return true
	}
	return false
}
