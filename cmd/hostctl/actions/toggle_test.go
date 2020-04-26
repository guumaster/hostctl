package actions

import (
	"testing"

	"github.com/guumaster/hostctl/pkg/types"
)

func Test_Toggle(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "remove")
	defer r.Clean()

	t.Run("Toggle", func(t *testing.T) {
		r.Run("hostctl toggle profile2").
			Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| profile2 | on     | 127.0.0.1 | first.loc  |
				| profile2 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile())
	})

	t.Run("Toggle unknown", func(t *testing.T) {
		r.RunE("hostctl toggle unknown", types.ErrUnknownProfile).
			Containsf(`
				[ℹ] Using hosts file: %s
			`, r.Hostfile())
	})
}
