package sshfs

func CreateDepartment(name string) (result []string) {
	return sshfsExec("sf_dept_add " + name)
}

func GetDepartmentCount() (departments []string) {
	return sshfsExec("sf_dept_list | wc -l")
}

func GetDepartmentsWithUsages() (departments []string) {
	return sshfsExec("sf_dept_usages_list")
}

func GetDepartmentUsers(department string) (users []string) {
	return sshfsExec("sf_dept_user_list " + department)
}
