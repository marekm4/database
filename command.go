package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func Exec(command string) (string, error) {
	if len(command) > 0 {
		args := strings.Fields(command)
		cmd := exec.Command(args[0], args[1:]...)
		output, err := cmd.Output()
		return fmt.Sprintf("%s", output), err
	}
	return "", nil
}
