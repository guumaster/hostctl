package actions

import (
	"testing"

	"github.com/guumaster/hostctl/pkg/types"
)

func Test_Disable(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "disable")
	defer r.Clean()

	t.Run("Disable", func(t *testing.T) {
		r.Run("hostctl disable profile1").
			Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| profile1 | off    | 127.0.0.1 | first.loc  |
				| profile1 | off    | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile())
	})

	t.Run("Disable Only", func(t *testing.T) {
		r.Run("hostctl disable profile1 --only").
			Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| profile1 | off    | 127.0.0.1 | first.loc  |
				| profile1 | off    | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile()).
			Run("hostctl list").
			Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| default  | on     | 127.0.0.1 | localhost  |
				+----------+--------+-----------+------------+
				| profile1 | off    | 127.0.0.1 | first.loc  |
				| profile1 | off    | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
				| profile2 | on     | 127.0.0.1 | first.loc  |
				| profile2 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile())
	})
}

func Test_EnableDisableErrors(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "enableDisableErrors")
	defer r.Clean()

	t.Run("Enable/Disable all error", func(t *testing.T) {
		r.RunE("hostctl disable something --all", ErrIncompatibleAllFlag).Empty()
		r.RunE("hostctl enable something --all", ErrIncompatibleAllFlag).Empty()
	})
}

func Test_EnableDisableUnknown(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "enableDisableUnknown")
	defer r.Clean()

	t.Run("Enable unknown", func(t *testing.T) {
		r.RunE("hostctl enable unknown", types.ErrUnknownProfile).
			Containsf(`
				[ℹ] Using hosts file: %s
			`, r.Hostfile())
	})

	t.Run("Disable unknown", func(t *testing.T) {
		r.RunE("hostctl disable unknown", types.ErrUnknownProfile).
			Containsf(`
				[ℹ] Using hosts file: %s
			`, r.Hostfile())
	})
}

func Test_EnableDisableAll(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "disable")
	defer r.Clean()

	t.Run("Disable All", func(t *testing.T) {
		r.Run("hostctl disable --all").
			Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| default  | on     | 127.0.0.1 | localhost  |
				+----------+--------+-----------+------------+
				| profile1 | off    | 127.0.0.1 | first.loc  |
				| profile1 | off    | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
				| profile2 | off    | 127.0.0.1 | first.loc  |
				| profile2 | off    | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile())
	})

	t.Run("Enable All", func(t *testing.T) {
		r.Run("hostctl enable --all").
			Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| default  | on     | 127.0.0.1 | localhost  |
				+----------+--------+-----------+------------+
				| profile1 | on     | 127.0.0.1 | first.loc  |
				| profile1 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
				| profile2 | on     | 127.0.0.1 | first.loc  |
				| profile2 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile())
	})
}

func Test_Enable(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "disable")
	defer r.Clean()

	t.Run("Enable", func(t *testing.T) {
		r.Run("hostctl enable profile1").
			Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| profile1 | on     | 127.0.0.1 | first.loc  |
				| profile1 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile())
	})

	t.Run("Enable Only", func(t *testing.T) {
		r.Run("hostctl enable profile1 --only").
			Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| profile1 | on     | 127.0.0.1 | first.loc  |
				| profile1 | on     | 127.0.0.1 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile()).
			Run("hostctl list").
			Containsf(`
				[ℹ] Using hosts file: %s

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
			`, r.Hostfile())
	})
}
