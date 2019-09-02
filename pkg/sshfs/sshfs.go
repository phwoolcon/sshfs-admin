package sshfs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func localExec(command string, arg ...string) (result []string) {
	fmt.Println("exec: ", command, arg)
	cmd := exec.Command(command, arg...)
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

func sshfsExec(command string) (result []string) {
	return localExec("./scripts/sshfs", command)
}
