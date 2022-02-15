package data

import (
	"github.com/Aserose/ModusOperandi/pkg/config"
	"github.com/Aserose/ModusOperandi/pkg/logger"
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
	"github.com/boltdb/bolt"
)

type InstructionStorage interface {
	PutInstruction(newInstruction model.Instruction) error
	RenameInstruction(newName, oldName string)
	AddPaths(updatedInstruction model.Instruction)
	ChangePath(updatedInstruction model.Instruction, oldPath string)
	GetAll() []model.Instruction
	GetInstruction(key string) model.Instruction
	DeletePath(instructionName, path string)
	DeleteInstruction(key string)
	DeleteAll()
}

type BoltData struct {
	InstructionStorage
}

func NewBoltData(db *bolt.DB, cfg config.CfgBolt, log logger.Logger) *BoltData {
	return &BoltData{
		InstructionStorage: NewInstructionStorage(db, cfg, log),
	}
}
