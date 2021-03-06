package sshfs

import (
	"fmt"
	"sshfs-admin/pkg/base"
)

func CreateUser(name, department string) (result []string) {
	sshKey := base.LocalExec("./scripts/sshkey_gen", name)[0]
	result = sshfsExec(fmt.Sprintf(`sf_user_add_only "%s" "%s" "%s"`, name, department, sshKey))
	if result[0] == "ok" {
		base.LocalExec("./scripts/sshkey_up", name)
	}
	return result
}

func GetUserCount() (users []string) {
	return sshfsExec("sf_user_list | wc -l")
}

func GetUserDepartments(name string) (result []string) {
	return sshfsExec(fmt.Sprintf(`sf_user_dept_list "%s"`, name))
}

func GetUsersWithUsages() (users []string) {
	return sshfsExec("sf_user_usages_list")
}

func RenameUser(name string, newName string) []string {
	return sshfsExec(fmt.Sprintf(`sf_user_rename "%s" "%s"`, name, newName))
}

func RegenerateKey(name string) (result []string) {
	sshKey := base.LocalExec("./scripts/sshkey_gen", name)[0]
	result = sshfsExec(fmt.Sprintf(`sf_user_update "%s" "%s" "%s"`, name, "", sshKey))
	if result[0] == "ok" {
		base.LocalExec("./scripts/sshkey_up", name)
	}
	return result
}

func UpdateUserDepartment(name, department string) []string {
	return sshfsExec(fmt.Sprintf(`sf_user_update "%s" "%s"`, name, department))
}

func UserExists(name string) bool {
	return sshfsExec(fmt.Sprintf(`sf_user_exists "%s"`, name))[0] == "1"
}
