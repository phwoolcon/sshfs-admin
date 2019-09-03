package sshfs

import (
	"sshfs-admin/pkg/base"
)

func sshfsExec(command string) (result []string) {
	return base.LocalExec("./scripts/sshfs", command)
}
