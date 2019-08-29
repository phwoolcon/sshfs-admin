package sshfs

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetDepartments() (departments []string) {
	return sshfsExec("sf_get_departments")
}

func GetDepartmentUsers(department string) (users []string) {
	return sshfsExec("sf_get_department_users " + department)
}

func GetUsers() (users []string) {
	return sshfsExec("sf_get_users")
}

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
