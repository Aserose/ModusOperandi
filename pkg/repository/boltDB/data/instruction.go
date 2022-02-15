package data

import (
	"encoding/json"
	"github.com/Aserose/ModusOperandi/pkg/config"
	"github.com/Aserose/ModusOperandi/pkg/logger"
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
	"github.com/boltdb/bolt"
)

type instructionStorage struct {
	db  *bolt.DB
	cfg config.CfgBolt
	log logger.Logger
}

func NewInstructionStorage(db *bolt.DB, cfg config.CfgBolt, log logger.Logger) *instructionStorage {
	return &instructionStorage{
		db:  db,
		cfg: cfg,
		log: log,
	}
}

func (i *instructionStorage) PutInstruction(newInstruction model.Instruction) error {

	jsonInstruction, err := json.Marshal(newInstruction)
	if err != nil {
		i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
		return err
	}

	if err := i.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(i.cfg.BucketInstructionName)).Put([]byte(newInstruction.Name), jsonInstruction)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
		return err
	}
	return nil
}

func (i instructionStorage) RenameInstruction(newName, oldName string) {
	instruction := i.GetInstruction(oldName)
	i.DeleteInstruction(oldName)

	instruction.Name = newName

	if err := i.PutInstruction(instruction); err != nil {
		i.log.Errorf("%s %s", i.log.CallInfoStr(), err.Error())
	}
}

func (i instructionStorage) AddPaths(updatedInstruction model.Instruction) {
	instruction := i.GetInstruction(updatedInstruction.Name)

	instruction.PathFile = append(instruction.PathFile, updatedInstruction.PathFile...)

	if err := i.PutInstruction(instruction); err != nil {
		i.log.Errorf("%s %s", i.log.CallInfoStr(), err.Error())
	}
}

func (i instructionStorage) ChangePath(updatedInstruction model.Instruction, oldPath string) {
	instruction := i.GetInstruction(updatedInstruction.Name)

	for i, p := range instruction.PathFile {
		if p == oldPath {
			instruction.PathFile[i] = updatedInstruction.PathFile[0]
		}
	}

	if err := i.PutInstruction(instruction); err != nil {
		i.log.Errorf("%s %s", i.log.CallInfoStr(), err.Error())
	}
}

func (i *instructionStorage) GetAll() []model.Instruction {
	var instructions []model.Instruction

	if err := i.db.View(func(tx *bolt.Tx) error {
		var decodedInstruction model.Instruction
		err := tx.Bucket([]byte(i.cfg.BucketInstructionName)).ForEach(func(k []byte, v []byte) error {

			err := json.Unmarshal(v, &decodedInstruction)
			if err != nil {
				i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
			}
			instructions = append(instructions, decodedInstruction)
			decodedInstruction = model.Instruction{}
			return nil
		})
		if err != nil {
			i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
		}
		return nil
	}); err != nil {
		i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
	}

	return instructions
}

func (i instructionStorage) GetInstruction(key string) model.Instruction {
	var instruction model.Instruction

	if err := i.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket([]byte(i.cfg.BucketInstructionName)).Get([]byte(key))
		if len(data) == 0 {
			return nil
		}
		if err := json.Unmarshal(data, &instruction); err != nil {
			i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
		}
		return nil
	}); err != nil {
		i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
	}

	return instruction
}

func (i instructionStorage) DeletePath(instructionName, path string) {
	instruction := i.GetInstruction(instructionName)

	for i, p := range instruction.PathFile {
		if p == path {
			instruction.PathFile = append(instruction.PathFile[:i], instruction.PathFile[i+1:]...)
		}
	}

	i.PutInstruction(instruction)
}

func (i instructionStorage) DeleteAll() {

}

func (i instructionStorage) DeleteInstruction(key string) {
	if err := i.db.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte(i.cfg.BucketInstructionName)).Delete([]byte(key))
		if err != nil {
			i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
		}
		return nil
	}); err != nil {
		i.log.Panicf("%s: %s", i.log.CallInfoStr(), err.Error())
	}
}
