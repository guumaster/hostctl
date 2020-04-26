package actions

import (
	"testing"
)

func Test_RemoveDomains(t *testing.T) {
	cmd := NewRootCmd()

	t.Run("Remove domains", func(t *testing.T) {
		r := NewRunner(t, cmd, "removeDomain")
		defer r.Clean()

		r.Run("hostctl remove domains profile1 first.loc").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Domains 'first.loc' removed.

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| profile1 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile())
	})

	t.Run("Remove domains and profile", func(t *testing.T) {
		r := NewRunner(t, cmd, "removeDomain")
		defer r.Clean()

		r.Run("hostctl remove domains profile1 first.loc second.loc").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Profile 'profile1' removed.
			`, r.Hostfile())
	})
}
