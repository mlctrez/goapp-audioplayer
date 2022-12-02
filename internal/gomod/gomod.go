package gomod

import (
	"encoding/json"
	"os/exec"
)

func ModuleName() (name string, err error) {
	var bytes []byte
	command := exec.Command("go", "mod", "edit", "-json")
	if bytes, err = command.CombinedOutput(); err != nil {
		return
	}
	editOutput := &EditOutput{}
	if err = json.Unmarshal(bytes, editOutput); err != nil {
		return
	}
	return editOutput.Module.Path, nil
}

type EditOutput struct {
	Module struct {
		Path string `json:"Path"`
	} `json:"Module"`
	Go      string `json:"Go"`
	Require []struct {
		Path     string `json:"Path"`
		Version  string `json:"Version"`
		Indirect bool   `json:"Indirect,omitempty"`
	} `json:"Require"`
}
