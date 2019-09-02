package sshfs

import "fmt"

func GetUsers() (users []string) {
	return sshfsExec("sf_user_list")
}

func CreateUser(name string, department string) (result []string) {
	sshKey := localExec("./scripts/sshkey_gen", name)[0]
	result = sshfsExec(fmt.Sprintf(`sf_user_add_only "%s" "%s" "%s"`, name, department, sshKey))
	if result[0] == "ok" {
		localExec("./scripts/sshkey_up", name)
	}
	return result
}
