package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Aserose/ModusOperandi/pkg/client/customButtons"
	"github.com/Aserose/ModusOperandi/pkg/client/handler"
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
)

type button struct {
	handler       *handler.Handler
	customButtons *customButtons.Button
	validation
	shortcutCreation
	manageInstruction
	bufferize
}

func NewButton(handler *handler.Handler, customButtons *customButtons.Button, v validation, sc shortcutCreation, m manageInstruction, s bufferize) *button {
	return &button{
		handler:           handler,
		customButtons:     customButtons,
		validation:        v,
		shortcutCreation:  sc,
		manageInstruction: m,
		bufferize:         s,
	}
}

func (b button) getInstructions(refresh refreshWindow, w fyne.Window) customButtons.InstructionSet {
	instructions := b.handler.GetAllInstructions()
	instructionButtons := b.customButtons.InstructionsWithPaths(instructions)

	instructionButtons.NewInstructionButton.OnTapped = func() {
		b.newInstruction(refresh)
	}

	b.instructionEditButtons(instructions, instructionButtons, refresh, w)

	return instructionButtons
}

func (b button) instructionEditButtons(instructions []model.Instruction, instructionButtons customButtons.InstructionSet, refresh refreshWindow, w fyne.Window) customButtons.InstructionSet {
	for _, instruction := range instructions {
		instructionName := instruction.Name
		paths := instruction.PathFile

		instructionButtons.InstructionList.EditButtonCanvas[instructionName] = append(instructionButtons.InstructionList.EditButtonCanvas[instructionName],
			layout.NewSpacer(),
			widget.NewButton("shortcut", func() {
				b.shortcutCreation.openWindowShortcut(paths)
			}),
			widget.NewButton("launch", func() {
				b.handler.Execute(paths)
			}),
			widget.NewButton("rename", func() {
				b.cleanInstruction()
				b.manageInstruction.renameInstruction(instructionName, w, refresh)
			}),
			widget.NewButton("delete", func() {
				b.cleanInstruction()
				b.manageInstruction.removeInstruction(instructionName, w, refresh)
			}),
			layout.NewSpacer())

		b.pathEditButtons(instructionButtons, instructionName, refresh, w)
	}
	return instructionButtons
}

func (b button) pathEditButtons(instructionButtons customButtons.InstructionSet, instructionName string, refresh refreshWindow, w fyne.Window) {
	for _, path := range instructionButtons.PathList[instructionName].ButtonCanvas {
		pathURI := path.(*widget.Button).Text
		if entry, ok := instructionButtons.PathList[instructionName]; ok {
			entry.EditButtonCanvas[pathURI] = append(entry.EditButtonCanvas[pathURI],
				layout.NewSpacer(),
				widget.NewButton("launch", func() {
					if b.validation.pathDoesNotExist(pathURI) {
						b.manageInstruction.removeNonExistentPath(instructionName, pathURI, refresh, w)
						return
					}
					b.handler.ExecuteOne(pathURI)
				}),
				widget.NewButton("change", func() {
					b.manageInstruction.changePath(instructionName, pathURI, refresh, w)
				}),
				widget.NewButton("remove", func() {
					b.manageInstruction.removePath(instructionName, pathURI, refresh, w)
				}),
				layout.NewSpacer())

			instructionButtons.PathList[instructionName] = entry
		}
	}
}
