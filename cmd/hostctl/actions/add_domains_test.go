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

	t.Run("Add domains uniq", func(t *testing.T) {
		r := NewRunner(t, cmd, "add_uniq")
		defer r.Clean()

		r.Run("hostctl add domains awesome same.domain.loc --ip 3.3.3.3").
			Run("hostctl add domains awesome same.domain.loc --ip 3.3.3.3").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Domains 'same.domain.loc' added.

        +---------+--------+---------+-----------------+
        | PROFILE | STATUS |   IP    |     DOMAIN      |
        +---------+--------+---------+-----------------+
        | awesome | on     | 3.3.3.3 | same.domain.loc |
        +---------+--------+---------+-----------------+
			`, r.Hostfile())
	})

	t.Run("Add domains uniq per profile", func(t *testing.T) {
		r := NewRunner(t, cmd, "multi_profile")
		defer r.Clean()

		r.Run("hostctl add domains another same.domain.loc --ip 3.3.3.3").
			Run("hostctl add domains another same.domain.loc --ip 3.3.3.3").
			Run("hostctl add domains awesome same.domain.loc --ip 3.3.3.3").
			Run("hostctl add domains awesome same.domain.loc --ip 3.3.3.3").
			Run("hostctl list another").
			Containsf(`
				[ℹ] Using hosts file: %s

				+---------+--------+---------+-----------------+
				| PROFILE | STATUS |   IP    |     DOMAIN      |
				+---------+--------+---------+-----------------+
				| another | on     | 3.3.3.3 | same.domain.loc |
				+---------+--------+---------+-----------------+
			`, r.Hostfile()).
			Run("hostctl list awesome").
			Containsf(`
				[ℹ] Using hosts file: %s

				+---------+--------+---------+-----------------+
				| PROFILE | STATUS |   IP    |     DOMAIN      |
				+---------+--------+---------+-----------------+
				| awesome | on     | 3.3.3.3 | same.domain.loc |
				+---------+--------+---------+-----------------+
			`, r.Hostfile())
	})
}
