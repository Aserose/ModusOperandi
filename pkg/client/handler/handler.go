package handler

import (
	"fmt"
	"github.com/Aserose/ModusOperandi/pkg/logger"
	"github.com/Aserose/ModusOperandi/pkg/repository/model"
	"github.com/Aserose/ModusOperandi/pkg/service"
	"os"
	"os/exec"
	"strings"
)

type Handler struct {
	service *service.Service
	log     logger.Logger
}

func NewHandler(service *service.Service, log logger.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h Handler) RenameInstruction(newName, oldName string) {
	h.service.RenameInstruction(newName, oldName)
}

func (h Handler) GetInstruction(instructionName string) model.Instruction {
	return h.service.GetInstruction(instructionName)
}

func (h Handler) InstructionsExist(instructionName string) bool {
	return h.service.InstructionsExist(instructionName)
}

func (h Handler) AddPaths(updatedInstruction model.Instruction) {
	h.service.AddPaths(updatedInstruction)
}

func (h Handler) ChangePath(updatedInstruction model.Instruction, oldPath string) {
	h.service.ChangePath(updatedInstruction, oldPath)
}

func (h Handler) GetAllInstructions() []model.Instruction {
	return h.service.Instruction.GetAllInstructions()
}

func (h Handler) RemoveInstruction(name string) {
	h.service.RemoveInstruction(name)
}

func (h Handler) RemovePath(instructionName, path string) {
	h.service.RemovePath(instructionName, path)
}

func (h Handler) NewInstruction(instruction model.Instruction) {
	if err := h.service.Instruction.Put(instruction); err != nil {
		h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
	}
}

func (h Handler) CreateShortcut(instructionPaths []string, fileName, directory string) {
	var filePath string
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.Mkdir(directory, 0755); err != nil {
			h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
		}
	}
	filePath = fmt.Sprintf("%s/%s.ps1", directory, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
	}
	if err := out.Close(); err != nil {
		h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
	}
	var result string
	for i := 0; i <= len(instructionPaths)-1; i++ {
		path := instructionPaths[i]

		result += fmt.Sprintf("& %s\n", path)
	}
	a, _ := os.OpenFile(filePath, os.O_WRONLY, 0755)
	if _, err := a.Write([]byte(result)); err != nil {
		h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
	}
	if err := a.Close(); err != nil {
		h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
	}

	if err := exec.Command("powershell", "ps2exe", filePath, strings.Replace(filePath, "ps1", "exe", 1)).Run(); err != nil {
		h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
	}

	if err := os.Remove(filePath); err != nil {
		h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
	}
}

func (h Handler) ExecuteOne(instructionPath string) {
	if err := exec.Command("powershell", "&", instructionPath).Run(); err != nil {
		h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
	}
}

func (h Handler) Execute(instructionPaths []string) {
	for i := 0; i <= len(instructionPaths)-1; i++ {
		path := instructionPaths[i]

		go func() {
			if err := exec.Command("powershell", "&", path).Run(); err != nil {
				h.log.Printf("%s : %s", h.log.CallInfoStr(), err.Error())
			}
		}()
	}
}
