package actions

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_Backup(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "backup")
	defer r.Clean()

	date := time.Now().UTC().Format("20060102")

	backupFile := fmt.Sprintf("%s.%s", r.Hostfile(), date)
	defer os.Remove(backupFile)

	r.Run("hostctl backup --path /tmp").
		Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Backup '%s' created.

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| default  | on     | 127.0.0.1 | localhost  |
				+----------+--------+-----------+------------+
				| profile1 | on     | 127.0.0.1 | first.loc  |
				| profile1 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
				| profile2 | off    | 127.0.0.1 | first.loc  |
				| profile2 | off    | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile(), backupFile)
}

func Test_Backup_Deep_Target(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "/deep/backup-deep")
	defer r.Clean()

	date := time.Now().UTC().Format("20060102")

	backupFile := fmt.Sprintf("%s.%s", r.Hostfile(), date)
	defer os.Remove(backupFile)

	r.Run("hostctl backup --path /tmp/deep").
		Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Backup '%s' created.

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| default  | on     | 127.0.0.1 | localhost  |
				+----------+--------+-----------+------------+
				| profile1 | on     | 127.0.0.1 | first.loc  |
				| profile1 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
				| profile2 | off    | 127.0.0.1 | first.loc  |
				| profile2 | off    | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile(), backupFile)
}
