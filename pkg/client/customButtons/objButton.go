package customButtons

import (
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
)

type Instruction interface {
	InstructionsWithPaths(instructions []model.Instruction) InstructionSet
}

type Button struct {
	Instruction
}

func NewButtons() *Button {
	return &Button{
		Instruction: NewInstructionSet(),
	}
}
