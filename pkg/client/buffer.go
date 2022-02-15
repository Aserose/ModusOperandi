package client

import "fyne.io/fyne/v2"

type buffer struct {
	pathContainer        *fyne.Container
	pathURI              string
	instructionContainer *fyne.Container
	instructionName      string
}

func NewBuffer() *buffer {
	return &buffer{}
}

func (b buffer) getInstructionContainer() *fyne.Container {
	return b.instructionContainer
}

func (b buffer) getInstruction() (*fyne.Container, string) {
	return b.instructionContainer, b.instructionName
}

func (b *buffer) memorizeInstruction(InstructionListContainer *fyne.Container, instructionName string) {
	b.instructionContainer = InstructionListContainer
	b.instructionName = instructionName
}

func (b *buffer) memorizePathContainer(pathContainer *fyne.Container) {
	b.pathContainer = pathContainer
}

func (b *buffer) memorizePathName(path string) {
	b.pathURI = path
}

func (b *buffer) memorizeInstructionName(name string) {
	b.instructionName = name
}

func (b buffer) getInstructionName() string {
	return b.instructionName
}

func (b *buffer) cleanPathURI() {
	b.pathURI = ""
}

func (b *buffer) cleanInstructionContainer() {
	b.instructionContainer = nil
}

func (b *buffer) cleanInstruction() {
	b.instructionContainer = nil
	b.instructionName = ""
}

func (b *buffer) cleanPathContainer() {
	b.pathContainer = nil
}

func (b buffer) getPathName() string {
	return b.pathURI
}

func (b buffer) getPathContainer() *fyne.Container {
	return b.pathContainer
}
