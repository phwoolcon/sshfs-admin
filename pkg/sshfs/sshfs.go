package sshfs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func sshfsExec(command string) (result []string) {
	cmd := exec.Command("./scripts/sshfs", command)
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	output := strings.TrimSpace(string(cmdOutput.Bytes()))
	result = strings.Split(output, "\n")
	return result
}
