package template_parser

import (
	"fmt"
	"os/exec"
)

func RunGoimportsOnFile(path string) {
	cmd := exec.Command("goimports", "-w", path)

	err := cmd.Run()

	if err != nil {
		fmt.Println("goimports run on ", path, "failed", err)
	}
}
