package cprinter

import (
	"fmt"
	"os"
	"os/exec"
)

func Print(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %s", filePath)
	}
	defer file.Close()

	cmd := exec.Command("cat", filePath)
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting command: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for command: %v", err)
	}

	return nil
}
