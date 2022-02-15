package customButtons

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
)

type instructionButtons struct {
}

type InstructionSet struct {
	NewInstructionButton *widget.Button
	InstructionList      buttonSet
	NewPathButton        *widget.Button
	PathList             map[string]buttonSet
}

type buttonSet struct {
	ButtonCanvas     []fyne.CanvasObject
	EditButtonCanvas map[interface{}][]fyne.CanvasObject
	IsTapped         bool
}

func NewInstructionSet() *instructionButtons {
	return &instructionButtons{}
}

func (ib instructionButtons) InstructionsWithPaths(instructions []model.Instruction) InstructionSet {
	instructionButtons := ib.initInstructionVariable()

	ib.createButtons(instructions, &instructionButtons)

	return instructionButtons
}

func (ib instructionButtons) initInstructionVariable() InstructionSet {
	var instructionButtons InstructionSet

	instructionButtons.NewInstructionButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})
	instructionButtons.NewInstructionButton.Importance = 1
	instructionButtons.NewPathButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {})
	instructionButtons.NewPathButton.Importance = 1

	instructionButtons.PathList = make(map[string]buttonSet)
	instructionButtons.InstructionList.EditButtonCanvas = make(map[interface{}][]fyne.CanvasObject)

	return instructionButtons
}

func (ib instructionButtons) createButtons(instructions []model.Instruction, ins *InstructionSet) {
	for _, instruction := range instructions {
		instructionName := instruction.Name
		pathList := instruction.PathFile

		ib.appendInstructionButtons(instructionName, ins)
		ib.appendPathListButtons(pathList, instructionName, ins)
	}
}

func (ib instructionButtons) appendInstructionButtons(instructionName string, ins *InstructionSet) {
	ins.InstructionList.ButtonCanvas = append(
		ins.InstructionList.ButtonCanvas,
		widget.NewButton(instructionName, func() {}))
}

func (ib instructionButtons) appendPathListButtons(paths []string, instructionName string, ins *InstructionSet) {
	var pathList buttonSet
	pathList.EditButtonCanvas = make(map[interface{}][]fyne.CanvasObject)

	for _, elem := range paths {
		pathList.ButtonCanvas = append(pathList.ButtonCanvas, widget.NewButton(elem, func() {}))
	}

	ins.PathList[instructionName] = pathList
}
