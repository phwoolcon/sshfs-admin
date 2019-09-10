package sshfs

func GetDiskUsage() []string {
	return sshfsExec("sf_disk_usages")
}
