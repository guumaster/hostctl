package actions

import (
	"testing"
)

func Test_AddDomains(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "addDomains")
	defer r.Clean()

	t.Run("Add domains", func(t *testing.T) {
		r.Run("hostctl add domains profile1 arg.domain.loc").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Domains 'arg.domain.loc' added.

				+----------+--------+-----------+----------------+
				| PROFILE  | STATUS |    IP     |     DOMAIN     |
				+----------+--------+-----------+----------------+
				| profile1 | on     | 127.0.0.1 | first.loc      |
				| profile1 | on     | 127.0.0.1 | second.loc     |
				| profile1 | on     | 127.0.0.1 | arg.domain.loc |
				+----------+--------+-----------+----------------+
`, r.Hostfile())
	})

	t.Run("Add domains new profile", func(t *testing.T) {
		r.Run("hostctl add domains newprofile arg.domain.loc").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Domains 'arg.domain.loc' added.

        +------------+--------+-----------+----------------+
        |  PROFILE   | STATUS |    IP     |     DOMAIN     |
        +------------+--------+-----------+----------------+
				| newprofile | on     | 127.0.0.1 | arg.domain.loc |
				+------------+--------+-----------+----------------+
`, r.Hostfile())
	})

	t.Run("Add domains with IP", func(t *testing.T) {
		r := NewRunner(t, cmd, "addWithIP")
		defer r.Clean()
		r.Run("hostctl add domains profile1 --ip 5.5.5.5 arg.domain.loc").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Domains 'arg.domain.loc' added.

				+----------+--------+-----------+----------------+
				| PROFILE  | STATUS |    IP     |     DOMAIN     |
				+----------+--------+-----------+----------------+
				| profile1 | on     | 127.0.0.1 | first.loc      |
				| profile1 | on     | 127.0.0.1 | second.loc     |
				| profile1 | on     | 5.5.5.5   | arg.domain.loc |
				+----------+--------+-----------+----------------+
`, r.Hostfile())
	})
}
