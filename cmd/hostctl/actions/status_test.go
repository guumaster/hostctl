package actions

import (
	"testing"
)

func Test_Status(t *testing.T) {
	cmd := NewRootCmd()
	r := NewRunner(t, cmd, "status")
	defer r.Clean()
	r.Run("hostctl status").
		Contains(`
			+----------+--------+
			| PROFILE  | STATUS |
			+----------+--------+
			| profile1 | on     |
			| profile2 | off    |
			+----------+--------+
		`)
}
