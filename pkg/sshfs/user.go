package sshfs

import (
	"fmt"
	"sshfs-admin/pkg/base"
)

func CreateUser(name string, department string) (result []string) {
	sshKey := base.LocalExec("./scripts/sshkey_gen", name)[0]
	result = sshfsExec(fmt.Sprintf(`sf_user_add_only "%s" "%s" "%s"`, name, department, sshKey))
	if result[0] == "ok" {
		base.LocalExec("./scripts/sshkey_up", name)
	}
	return result
}

func GetUserDepartments(name string) (result []string) {
	return sshfsExec(fmt.Sprintf(`sf_user_dept_list "%s"`, name))
}

func GetUsers() (users []string) {
	return sshfsExec("sf_user_list")
}
