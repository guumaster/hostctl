package actions

import (
	"testing"

	"github.com/guumaster/hostctl/pkg/types"
)

func Test_Remove(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "remove")
	defer r.Clean()

	t.Run("Remove", func(t *testing.T) {
		r.Run("hostctl remove profile2").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Profile(s) 'profile2' removed.
			`, r.Hostfile())
	})

	t.Run("Remove unknown", func(t *testing.T) {
		r.RunE("hostctl remove unknown", types.ErrUnknownProfile).
			Containsf(`[ℹ] Using hosts file: %s`, r.Hostfile())
	})

	t.Run("Remove all bad", func(t *testing.T) {
		r.RunE("hostctl remove profile1 --all", ErrIncompatibleAllFlag).Empty()
	})

	t.Run("Remove multiple", func(t *testing.T) {
		cmd := NewRootCmd()

		r := NewRunner(t, cmd, "remove")
		defer r.Clean()

		r.Run("hostctl remove profile1 profile2").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Profile(s) 'profile1, profile2' removed.
			`, r.Hostfile())
	})

	t.Run("Remove all", func(t *testing.T) {
		cmd := NewRootCmd()

		r := NewRunner(t, cmd, "remove")
		defer r.Clean()

		r.Run("hostctl remove --all").
			Containsf(`
				[ℹ] Using hosts file: %s

				[✔] Profile(s) 'profile1, profile2' removed.
			`, r.Hostfile())
	})
}
