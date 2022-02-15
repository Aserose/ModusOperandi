package client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Aserose/ModusOperandi/pkg/client/handler"
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
)

type instructionManager struct {
	handler *handler.Handler
	instructCreation
	validation
	bufferize
}

func NewInstructionManager(handler *handler.Handler, c instructCreation, s bufferize, v validation) *instructionManager {
	return &instructionManager{
		handler:          handler,
		instructCreation: c,
		validation:       v,
		bufferize:        s,
	}
}

func (im instructionManager) newInstruction(refresh refreshWindow) {
	creationWindow := im.instructCreation.getWindowCreation()

	var newInstructionList model.Instruction

	entryForNameInstruction, entryForInstructionPaths := im.creationWindowEntry()

	confirmButton := widget.NewButton("confirm", func() {
		if im.validation.instructionNameExist(entryForNameInstruction.Text) {
			return
		}
		newInstructionList.Name = entryForNameInstruction.Text
		im.handler.NewInstruction(newInstructionList)
		refresh()
		creationWindow.Close()
	})

	addButton, wipeButton, removeButton := im.buttonsForEditPaths(entryForInstructionPaths, confirmButton, &newInstructionList)

	showAddedInstructionPaths := im.showAddedPaths(removeButton, confirmButton, &newInstructionList)

	creationWindow.SetContent(
		container.NewVSplit(
			im.creationWindowForms(entryForInstructionPaths, entryForNameInstruction, creationWindow, addButton, wipeButton),
			container.NewVSplit(
				showAddedInstructionPaths,
				container.NewVBox(
					removeButton,
					layout.NewSpacer(),
					confirmButton),
			),
		),
	)

	creationWindow.Show()
}

func (im instructionManager) creationWindowEntry() (entryForNameInstruction, entryForInstructionPaths *widget.Entry) {
	return widget.NewEntry(), widget.NewEntry()
}

func (im instructionManager) creationWindowForms(entryForInstructionPaths, entryForNameInstruction *widget.Entry, creationWindow fyne.Window, addButton, wipeButton *widget.Button) *widget.Form {
	return widget.NewForm(
		widget.NewFormItem("instruction name", entryForNameInstruction),
		im.setFilePathFormItem(entryForInstructionPaths, creationWindow, addButton, wipeButton))
}

func (im instructionManager) removeInstruction(instructionName string, w fyne.Window, refresh refreshWindow) (modal *widget.PopUp) {
	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Are you sure?"),
			container.NewGridWithColumns(2,
				widget.NewButton("No", func() { modal.Hide() }),
				widget.NewButton("Yes", func() { im.handler.RemoveInstruction(instructionName); refresh(); modal.Hide() }),
			),
		), w.Canvas())

	modal.Show()

	return modal
}

func (im instructionManager) renameInstruction(oldName string, w fyne.Window, refresh refreshWindow) (modal *widget.PopUp) {
	entryForNewName := widget.NewEntry()

	modal = widget.NewModalPopUp(
		container.NewVSplit(

			widget.NewForm(widget.NewFormItem("new name", entryForNewName)),

			container.NewVBox(
				layout.NewSpacer(),
				container.NewGridWithColumns(2,
					widget.NewButton("cancel", func() { modal.Hide() }),
					widget.NewButton("confirm", func() {
						if im.validation.instructionNameExist(entryForNewName.Text) {
							return
						}
						if im.validation.instructionNameIsBlank(entryForNewName.Text) {
							return
						}
						im.handler.RenameInstruction(entryForNewName.Text, oldName)
						refresh()
						modal.Hide()
					}),
				),
			),
		),
		w.Canvas())

	modal.Show()

	return modal
}

func (im instructionManager) showAddedPaths(removeButton, confirmButton *widget.Button, newInstructionList *model.Instruction) *widget.List {
	removeButton.Disable()

	addedInstructionPaths := widget.NewList(
		func() int {
			return len(newInstructionList.PathFile)
		},
		func() fyne.CanvasObject {
			if len(newInstructionList.PathFile) <= 0 {
				confirmButton.Disable()
			}
			return widget.NewLabel("")
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(newInstructionList.PathFile[lii])
		})

	addedInstructionPaths.OnSelected = func(id widget.ListItemID) {
		removeButton.Enable()
		removeButton.OnTapped = func() {
			newInstructionList.PathFile = append(newInstructionList.PathFile[:id], newInstructionList.PathFile[id+1:]...)
			addedInstructionPaths.UnselectAll()
			removeButton.Disable()
		}
	}

	return addedInstructionPaths
}

func (im instructionManager) formattingAddPath(path []rune) string {
	if path[0] != '"' && path[len(path)-1] != '"' {
		path = append([]rune{'"'}, path...)
		path = append(path, '"')
		return string(path)
	}

	if path[0] == '"' {
		if path[len(path)-1] != '"' {
			path = append(path, '"')
		}
	} else {
		path = append([]rune{'"'}, path...)
	}

	return string(path)
}

func (im instructionManager) buttonsForEditPath(entryForInstructionPaths *widget.Entry, confirmButton *widget.Button,
	newInstructionList *model.Instruction) (addButton, wipeButton, removeButton *widget.Button) {

	existedInstruction := im.handler.GetInstruction(newInstructionList.Name)

	return widget.NewButton("add", func() {

			if entryForInstructionPaths.Text == "" { return }

			addingPath := im.formattingAddPath([]rune(entryForInstructionPaths.Text))

			if im.validation.pathDoesNotValid(addingPath) { return }
			if im.validation.pathAlreadyAdded(newInstructionList.PathFile, addingPath) { return }
			if existedInstruction.Name != "" {
				if im.validation.pathAlreadyExist(existedInstruction.PathFile, addingPath) {
					return
				}
			}

			newInstructionList.PathFile[0] = addingPath

			confirmButton.Enable()

			entryForInstructionPaths.SetText("")
		}),

		widget.NewButton("wipe", func() {
			entryForInstructionPaths.SetText("")
		}),

		widget.NewButton("remove path", func() {})
}

func (im instructionManager) buttonsForEditPaths(entryForInstructionPaths *widget.Entry, confirmButton *widget.Button,
	newInstructionList *model.Instruction) (addButton, wipeButton, removeButton *widget.Button) {
	var existedInstruction model.Instruction

	if newInstructionList.Name != "" {
		existedInstruction = im.handler.GetInstruction(newInstructionList.Name)
	}

	return widget.NewButton("add", func() {
			if entryForInstructionPaths.Text == "" { return }

			addingPath := im.formattingAddPath([]rune(entryForInstructionPaths.Text))

			if im.validation.pathDoesNotValid(addingPath) { return }
			if im.validation.pathAlreadyAdded(newInstructionList.PathFile, addingPath) { return }
			if existedInstruction.Name != "" {
				if im.validation.pathAlreadyExist(existedInstruction.PathFile, addingPath) {
					return
				}
			}
			newInstructionList.PathFile = append(newInstructionList.PathFile, addingPath)
			confirmButton.Enable()
			entryForInstructionPaths.SetText("")
		}),

		widget.NewButton("wipe", func() { entryForInstructionPaths.SetText("") }),

		widget.NewButton("remove path", func() {})
}

func (im instructionManager) setFilePathFormItem(entry *widget.Entry, a fyne.Window, addButton, wipeButton *widget.Button) *widget.FormItem {
	return widget.NewFormItem("file path", container.NewVBox(
		container.NewHSplit(
			entry,
			widget.NewButtonWithIcon("", theme.FolderNewIcon(), func() {

				openFileDialog := dialog.NewFileOpen(
					func(r fyne.URIReadCloser, _ error) {
						if r == nil {
							return
						} else {
							entry.SetText(r.URI().Path())
						}
					}, a)

				openFileDialog.Show()
			})), container.NewAdaptiveGrid(2, addButton, wipeButton)))
}

func (im instructionManager) removeNonExistentPath(instructionName, pathURI string, refresh refreshWindow, w fyne.Window) (modal *widget.PopUp) {
	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("File not found.\nRemove the path?"),
			container.NewGridWithColumns(2,
				widget.NewButton("No", func() { modal.Hide() }),
				widget.NewButton("Yes", func() { im.handler.RemovePath(instructionName, pathURI); refresh(); modal.Hide() }),
			),
		), w.Canvas())

	modal.Show()

	return modal
}

func (im instructionManager) removePath(instructionName, pathURI string, refresh refreshWindow, w fyne.Window) (modal *widget.PopUp) {
	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Are you sure?"),
			container.NewGridWithColumns(2,
				widget.NewButton("No", func() { modal.Hide() }),
				widget.NewButton("Yes", func() { im.handler.RemovePath(instructionName, pathURI); refresh(); modal.Hide() }),
			),
		), w.Canvas())

	modal.Show()

	return modal
}

func (im instructionManager) addPaths(instructionName string, refresh refreshWindow, w fyne.Window) (modal *widget.PopUp) {
	updatedInstructionList := model.Instruction{Name: instructionName}

	entryForInstructionPaths := widget.NewEntry()

	confirmButton := widget.NewButton("confirm", func() {
		im.handler.AddPaths(updatedInstructionList)
		refresh()
		modal.Hide()
	})

	addButton, wipeButton, removeButton := im.buttonsForEditPaths(entryForInstructionPaths, confirmButton, &updatedInstructionList)

	showAddedPaths := im.showAddedPaths(removeButton, confirmButton, &updatedInstructionList)

	modal = widget.NewModalPopUp(
		container.NewVSplit(

			widget.NewForm(
				im.setFilePathFormItem(entryForInstructionPaths, w, addButton, wipeButton)),

			container.NewVSplit(
				showAddedPaths,

				container.NewVBox(
					removeButton,
					layout.NewSpacer(),
					container.NewGridWithColumns(2,
						widget.NewButton("cancel", func() { modal.Hide() }),
						confirmButton,
					),
				),
			),
		),
		w.Canvas())

	modal.Resize(fyne.NewSize(400, 350))

	modal.Show()
	return modal

}

func (im instructionManager) changePath(instructionName, oldPath string, refresh refreshWindow, w fyne.Window) (modal *widget.PopUp) {
	updatedInstructionList := model.Instruction{Name: instructionName, PathFile: make([]string, 1)}

	entryForInstructionPaths := widget.NewEntry()

	confirmButton := widget.NewButton("confirm", func() {
		im.handler.ChangePath(updatedInstructionList, oldPath)
		im.cleanInstructionContainer()
		refresh()
		modal.Hide()
	})

	addButton, wipeButton, removeButton := im.buttonsForEditPath(entryForInstructionPaths, confirmButton, &updatedInstructionList)

	showAddedPaths := im.showAddedPaths(removeButton, confirmButton, &updatedInstructionList)

	modal = widget.NewModalPopUp(
		container.NewVSplit(

			widget.NewForm(
				im.setFilePathFormItem(entryForInstructionPaths, w, addButton, wipeButton)),

			container.NewVSplit(
				showAddedPaths,

				container.NewVBox(
					removeButton,
					layout.NewSpacer(),
					container.NewGridWithColumns(2,
						widget.NewButton("cancel", func() { modal.Hide() }),
						confirmButton,
					),
				),
			),
		),
		w.Canvas())

	modal.Resize(fyne.NewSize(400, 350))

	modal.Show()
	return modal
}
