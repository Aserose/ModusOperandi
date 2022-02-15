package service

import (
	"github.com/Aserose/ModusOperandi/pkg/logger"
	"github.com/Aserose/ModusOperandi/pkg/repository"
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
)

type instructionService struct {
	db  *repository.DB
	log logger.Logger
}

func NewInstruction(db *repository.DB, log logger.Logger) *instructionService {
	return &instructionService{
		db:  db,
		log: log,
	}
}

func (i instructionService) GetInstruction(instructionName string) model.Instruction {
	return i.db.BoltDB.GetInstruction(instructionName)
}

func (i instructionService) RenameInstruction(newName, oldName string) {
	i.db.BoltDB.RenameInstruction(newName, oldName)
}

func (i instructionService) InstructionsExist(instructionName string) bool {
	if i.db.BoltDB.GetInstruction(instructionName).Name == "" {
		return false
	} else {
		return true
	}
}

func (i instructionService) Put(instruction model.Instruction) error {
	return i.db.BoltDB.PutInstruction(instruction)
}

func (i instructionService) RemovePath(instructionName, path string) {
	i.db.BoltDB.DeletePath(instructionName, path)
}

func (i instructionService) RemoveInstruction(name string) {
	i.db.BoltDB.DeleteInstruction(name)
}

func (i instructionService) AddPaths(updatedInstruction model.Instruction) {
	i.db.BoltDB.AddPaths(updatedInstruction)
}

func (i instructionService) ChangePath(updatedInstruction model.Instruction, oldPath string) {
	i.db.BoltDB.ChangePath(updatedInstruction, oldPath)
}

func (i instructionService) GetAllInstructions() []model.Instruction {
	return i.db.BoltDB.GetAll()
}
