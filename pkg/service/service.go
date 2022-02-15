package service

import (
	"github.com/Aserose/ModusOperandi/pkg/logger"
	"github.com/Aserose/ModusOperandi/pkg/repository"
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
)

type Instruction interface {
	Put(instruction model.Instruction) error
	RenameInstruction(newName, oldName string)
	InstructionsExist(instructionName string) bool
	AddPaths(updatedInstruction model.Instruction)
	ChangePath(updatedInstruction model.Instruction, oldPath string)
	RemoveInstruction(name string)
	RemovePath(instructionName, path string)
	GetAllInstructions() []model.Instruction
	GetInstruction(instructionName string) model.Instruction
}

type Service struct {
	Instruction
}

func NewService(db *repository.DB, log logger.Logger) *Service {
	return &Service{
		Instruction: NewInstruction(db, log),
	}
}
