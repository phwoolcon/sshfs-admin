package sshfs

func GetUsers() (users []string) {
	return sshfsExec("sf_user_list")
}
