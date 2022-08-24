package file

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/types"
)

func TestFile_EnableDisable(t *testing.T) {
	mem := createBasicFS(t)
	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	t.Run("Enable", func(t *testing.T) {
		err = m.Enable([]string{"profile2"})
		assert.NoError(t, err)
		assert.Contains(t, m.GetEnabled(), "profile2")
	})

	t.Run("Enable Only", func(t *testing.T) {
		err = m.Enable([]string{"profile1", "profile2"})
		err = m.EnableOnly([]string{"default", "profile2"})
		assert.NoError(t, err)
		assert.Contains(t, m.GetEnabled(), "profile2")
		assert.Contains(t, m.GetDisabled(), "profile1")
	})

	t.Run("Enable All", func(t *testing.T) {
		err = m.EnableAll()
		assert.NoError(t, err)
		enabled := m.GetEnabled()
		assert.Contains(t, enabled, "profile1")
		assert.Contains(t, enabled, "profile2")
	})

	t.Run("Enable error", func(t *testing.T) {
		err = m.Enable([]string{"unknown"})
		assert.EqualError(t, err, types.ErrUnknownProfile.Error())
	})

	t.Run("Disable", func(t *testing.T) {
		err = m.Disable([]string{"profile2"})
		assert.NoError(t, err)
		assert.Contains(t, m.GetDisabled(), "profile2")
	})

	t.Run("Disable Only", func(t *testing.T) {
		err = m.Disable([]string{"profile1", "profile2"})
		err = m.DisableOnly([]string{"default", "profile2"})
		assert.NoError(t, err)
		assert.Contains(t, m.GetEnabled(), "profile1")
		assert.Contains(t, m.GetDisabled(), "profile2")
	})

	t.Run("Disable All", func(t *testing.T) {
		err = m.DisableAll()
		assert.NoError(t, err)
		Disabled := m.GetDisabled()
		assert.Contains(t, Disabled, "profile1")
		assert.Contains(t, Disabled, "profile2")
	})

	t.Run("Disable error", func(t *testing.T) {
		err = m.Disable([]string{"unknown"})
		assert.EqualError(t, err, types.ErrUnknownProfile.Error())
	})
}
