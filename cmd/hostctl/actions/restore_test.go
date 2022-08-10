package actions

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Restore(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "remove")
	defer r.Clean()

	from := r.TempHostfile("restoreFrom")
	defer os.Remove(from.Name())

	to, err := os.CreateTemp("/tmp", "restoreTo")
	assert.NoError(t, err)

	defer os.Remove(to.Name())

	r.Runf("hostctl restore --from %s --host-file %s", from.Name(), to.Name()).
		Containsf(`
			[ℹ] Using hosts file: %s

			[✔] File '%s' restored.

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
		`, to.Name(), from.Name())

	toData, _ := os.ReadFile(to.Name())
	fromData, _ := os.ReadFile(from.Name())
	assert.Equal(t, string(toData), string(fromData))
}
