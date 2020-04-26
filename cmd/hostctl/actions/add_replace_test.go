package actions

import (
	"os"
	"strings"
	"testing"
)

func Test_Add(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "add")
	defer r.Clean()

	tmp := r.TempHostfile("source")
	defer os.Remove(tmp.Name())

	t.Run("Add from file", func(t *testing.T) {
		r.Runf("hostctl add awesome --uniq --from %s", tmp.Name()).
			Contains(`
				+---------+--------+-----------+------------+
				| PROFILE | STATUS |    IP     |   DOMAIN   |
				+---------+--------+-----------+------------+
				| awesome | on     | 127.0.0.1 | localhost  |
				| awesome | on     | 127.0.0.1 | first.loc  |
				| awesome | on     | 127.0.0.1 | second.loc |
				+---------+--------+-----------+------------+
			`)
	})

	t.Run("Add from stdin", func(t *testing.T) {
		r := NewRunner(t, cmd, "add")
		defer r.Clean()

		tmp := r.TempHostfile("source")
		defer os.Remove(tmp.Name())

		in := strings.NewReader(`3.3.3.3 stdin.loc`)
		cmd.SetIn(in)

		r.Run("hostctl add awesome").
			Containsf(`
				[ℹ] Using hosts file: %s

        +---------+--------+---------+-----------+
        | PROFILE | STATUS |   IP    |  DOMAIN   |
        +---------+--------+---------+-----------+
        | awesome | on     | 3.3.3.3 | stdin.loc |
        +---------+--------+---------+-----------+
			`, r.Hostfile())
	})
}

func Test_ReplaceStdin(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "replace")
	defer r.Clean()

	tmp := r.TempHostfile("source")
	defer os.Remove(tmp.Name())

	in := strings.NewReader(`3.3.3.3 stdin.replaced.loc`)
	cmd.SetIn(in)

	r.Run("hostctl replace profile1").
		Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+---------+--------------------+
				| PROFILE  | STATUS |   IP    |       DOMAIN       |
				+----------+--------+---------+--------------------+
				| profile1 | on     | 3.3.3.3 | stdin.replaced.loc |
				+----------+--------+---------+--------------------+
			`, r.Hostfile())
}

func Test_ReplaceFile(t *testing.T) {
	cmd := NewRootCmd()

	in := strings.NewReader(`
		5.5.5.5 replaced.loc
		5.5.5.6 replaced2.loc
`)
	cmd.SetIn(in)

	r := NewRunner(t, cmd, "replace")
	defer r.Clean()

	tmp := r.TempHostfile("source")
	defer os.Remove(tmp.Name())

	r.Run("hostctl replace awesome").
		Containsf(`
				[ℹ] Using hosts file: %s

				+---------+--------+---------+---------------+
				| PROFILE | STATUS |   IP    |    DOMAIN     |
				+---------+--------+---------+---------------+
				| awesome | on     | 5.5.5.5 | replaced.loc  |
				| awesome | on     | 5.5.5.6 | replaced2.loc |
				+---------+--------+---------+---------------+
			`, r.Hostfile())
}
