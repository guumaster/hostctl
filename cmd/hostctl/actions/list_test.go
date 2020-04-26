package actions

import (
	"testing"
)

func Test_List(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "list")
	defer r.Clean()

	t.Run("List all", func(t *testing.T) {
		r.Run("hostctl list").
			Contains(`
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
			`)
	})

	t.Run("List enabled", func(t *testing.T) {
		r.Run("hostctl list enabled").
			Contains(`
				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| default  | on     | 127.0.0.1 | localhost  |
				+----------+--------+-----------+------------+
				| profile1 | on     | 127.0.0.1 | first.loc  |
				| profile1 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`)
	})

	t.Run("List disabled", func(t *testing.T) {
		r.Run("hostctl list disabled").
			Contains(`
				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| default  | on     | 127.0.0.1 | localhost  |
				+----------+--------+-----------+------------+
				| profile2 | off    | 127.0.0.1 | first.loc  |
				| profile2 | off    | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`)
	})
}
