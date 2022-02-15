package model

type Instruction struct {
	Name     string   `json:"name"`
	PathFile []string `json:"path_file"`
}
