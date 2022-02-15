package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Aserose/ModusOperandi/pkg/client/customButtons"
	"github.com/Aserose/ModusOperandi/pkg/client/handler"
	"github.com/Aserose/ModusOperandi/pkg/config"
	"image/color"
)

type refreshWindow func()

type instructCreation interface {
	getWindowCreation() fyne.Window
}

type shortcutCreation interface {
	getWindowShortcut() fyne.Window
	openWindowShortcut(instructionPaths []string)
}

type notification interface {
	createNotification(content fyne.CanvasObject)
}

type buttons interface {
	getInstructions(refresh refreshWindow, w fyne.Window) customButtons.InstructionSet
}

type manageInstruction interface {
	newInstruction(r refreshWindow)
	renameInstruction(oldName string, w fyne.Window, r refreshWindow) (modal *widget.PopUp)
	removeInstruction(instructionName string, w fyne.Window, r refreshWindow) (modal *widget.PopUp)
	removeNonExistentPath(instructionName, pathURI string, refresh refreshWindow, w fyne.Window) (modal *widget.PopUp)
	removePath(instructionName, pathURI string, refresh refreshWindow, w fyne.Window) (modal *widget.PopUp)
	addPaths(instructionName string, r refreshWindow, w fyne.Window) (modal *widget.PopUp)
	changePath(instructionName, oldPath string, r refreshWindow, w fyne.Window) (modal *widget.PopUp)
}

type bufferize interface {
	memorizeInstruction(InstructionListContainer *fyne.Container, instructionName string)
	memorizeInstructionName(name string)
	memorizePathContainer(pathContainer *fyne.Container)
	memorizePathName(path string)
	getPathContainer() *fyne.Container
	getPathName() string
	getInstructionName() string
	getInstructionContainer() *fyne.Container
	getInstruction() (*fyne.Container, string)
	cleanInstruction()
	cleanInstructionContainer()
	cleanPathURI()
	cleanPathContainer()
}

type validation interface {
	instructionNameIsBlank(instructionName string) bool
	instructionNameExist(instructionName string) bool
	pathAlreadyAdded(paths []string, input string) bool
	pathAlreadyExist(paths []string, input string) bool
	pathDoesNotValid(path string) bool
	pathDoesNotExist(path string) bool
}

type windowMain struct {
	w fyne.Window
	buttons
	manageInstruction
	bufferize
}

func NewMainWindow(w fyne.Window, handler *handler.Handler, cfg config.CfgClient) *windowMain {

	n := NewNotificationWindow(cfg)
	v := NewValidator(handler, n)
	s := NewBuffer()
	m := NewInstructionManager(handler, NewInstructionCreation(cfg), s, v)

	return &windowMain{
		w:                 w,
		buttons:           NewButton(handler, customButtons.NewButtons(), v, NewWindowShortcut(handler, v, n, cfg), m, s),
		manageInstruction: m,
		bufferize:         s,
	}
}

func (w *windowMain) showInstructions() {
	//creating a set of buttons by the number of available instructions
	instructionButtons := w.buttons.getInstructions(w.showInstructions, w.w)

	for i, button := range instructionButtons.InstructionList.ButtonCanvas {
		before := instructionButtons.InstructionList.ButtonCanvas[:i]
		selectedBtn := instructionButtons.InstructionList.ButtonCanvas[i]
		after := instructionButtons.InstructionList.ButtonCanvas[i+1:]
		instructionName := button.(*widget.Button).Text

		button.(*widget.Button).OnTapped = func() {
			if instructionButtons.InstructionList.IsTapped == true {

				//check that the button is not tapped repeatedly
				if instructionName != w.bufferize.getInstructionName() {
					//open the list of instruction paths
					w.showAndEditInstructionPaths(
						w.wrapSelectedListIntoContainer(instructionButtons.NewInstructionButton, before,
							selectedBtn, instructionButtons.InstructionList.EditButtonCanvas[instructionName],
							after), instructionButtons, instructionName)
					//cache information about the new selected button
					w.bufferize.memorizeInstructionName(instructionName)

					return
				}

				//the re-selected button closes the filepath lists and edit buttons
				//container without selected buttons
				w.w.SetContent(w.baseContainer(nil, nil, instructionButtons, ""))

				instructionButtons.InstructionList.IsTapped = false
				return
			}

			instructionButtons.InstructionList.IsTapped = true
			//cache information about the selected button
			w.bufferize.memorizeInstructionName(instructionName)
			//list of instruction filepaths
			w.showAndEditInstructionPaths(
				w.wrapSelectedListIntoContainer(instructionButtons.NewInstructionButton, before,
					selectedBtn, instructionButtons.InstructionList.EditButtonCanvas[instructionName],
					after), instructionButtons, instructionName)

		}
	}
	//re-opens the updated list of selected instruction filepaths
	if w.bufferize.getInstructionName() != "" {
		w.toSelectedInstructionAfterEditPath(instructionButtons)
	} else {
		//container without selected buttons
		w.w.SetContent(w.baseContainer(nil, nil, instructionButtons, ""))
	}
}

func (w windowMain) toSelectedInstructionAfterEditPath(instructionButtons customButtons.InstructionSet) {
	for i, button := range instructionButtons.InstructionList.ButtonCanvas {
		before := instructionButtons.InstructionList.ButtonCanvas[:i]
		selectedBtn := instructionButtons.InstructionList.ButtonCanvas[i]
		after := instructionButtons.InstructionList.ButtonCanvas[i+1:]
		instructionName := button.(*widget.Button).Text

		if instructionName == w.bufferize.getInstructionName() {
			w.showAndEditInstructionPaths(
				w.wrapSelectedListIntoContainer(instructionButtons.NewInstructionButton, before,
					selectedBtn, instructionButtons.InstructionList.EditButtonCanvas[instructionName],
					after), instructionButtons, instructionName)

			w.cleanInstructionContainer()
			break
		}
	}
}

func (w *windowMain) showAndEditInstructionPaths(InstructionListContainer *fyne.Container, instructionButtons customButtons.InstructionSet, instructionName string) {
	w.memorizeInstruction(InstructionListContainer, instructionName)

	for i, button := range instructionButtons.PathList[instructionName].ButtonCanvas {
		before := instructionButtons.PathList[instructionName].ButtonCanvas[:i]
		selectedBtn := instructionButtons.PathList[instructionName].ButtonCanvas[i]
		after := instructionButtons.PathList[instructionName].ButtonCanvas[i+1:]
		pathIndex := button.(*widget.Button).Text

		instructionButtons.NewPathButton.OnTapped = func() { w.addPaths(instructionName, w.showInstructions, w.w) }

		button.(*widget.Button).OnTapped = func() {

			if instructionButtons.PathList[instructionName].IsTapped == true {
				//check that the button is not tapped repeatedly
				if pathIndex != w.bufferize.getPathName() {
					//container with the selected instruction and path buttons
					w.w.SetContent(
						w.baseContainer(InstructionListContainer,
							w.wrapSelectedListIntoContainer(instructionButtons.NewPathButton, before, selectedBtn,
								instructionButtons.PathList[instructionName].EditButtonCanvas[pathIndex], after), instructionButtons, instructionName))

					w.bufferize.memorizePathName(pathIndex)
					return
				}

				//re-selected button closes the list of edit buttons (delete/remove, rename, etc)
				//container without the selected path button
				w.w.SetContent(w.baseContainer(InstructionListContainer, nil, instructionButtons, instructionName))

				if entry, ok := instructionButtons.PathList[instructionName]; ok {
					entry.IsTapped = false
					instructionButtons.PathList[instructionName] = entry
				}

				return
			}

			if entry, ok := instructionButtons.PathList[instructionName]; ok {
				entry.IsTapped = true
				instructionButtons.PathList[instructionName] = entry
			}

			w.bufferize.memorizePathName(pathIndex)

			//container with the selected instruction and path buttons
			w.w.SetContent(
				w.baseContainer(InstructionListContainer,
					w.wrapSelectedListIntoContainer(instructionButtons.NewPathButton, before, selectedBtn,
						instructionButtons.PathList[instructionName].EditButtonCanvas[pathIndex], after), instructionButtons, instructionName))

		}
	}

	//container with the selected button from the instruction list
	w.w.SetContent(w.baseContainer(InstructionListContainer, nil, instructionButtons, instructionName))

}

func (w windowMain) baseContainer(instructionListContainer *fyne.Container, pathListContainer *fyne.Container,
	ins customButtons.InstructionSet, instructionName string) *fyne.Container {

	//container without selected buttons
	if instructionListContainer == nil {
		return container.NewGridWithColumns(2,
			container.NewVBox(
				container.NewVBox(ins.NewInstructionButton),
				container.NewVBox(ins.InstructionList.ButtonCanvas...)))
	}
	//container without the selected path button
	if pathListContainer == nil {
		return container.NewGridWithColumns(2,
			instructionListContainer,
			container.NewVBox(
				container.NewVBox(ins.NewPathButton),
				container.NewVBox(ins.PathList[instructionName].ButtonCanvas...)))
	}
	//container with the selected instruction and path buttons
	return container.NewGridWithColumns(2,
		instructionListContainer,
		pathListContainer)

}

func (w windowMain) wrapSelectedListIntoContainer(new *widget.Button, before []fyne.CanvasObject, selectedBtn fyne.CanvasObject,
	edit []fyne.CanvasObject, after []fyne.CanvasObject) *fyne.Container {

	btn_color := canvas.NewRectangle(color.NRGBA{R: 36, G: 37, B: 42, A: 155})

	return container.NewVBox(
		container.NewVBox(new),
		container.NewVBox(before...),
		container.NewVBox(container.New(layout.NewMaxLayout(), btn_color, selectedBtn)),
		container.NewHBox(edit...),
		container.NewVBox(after...))

}
