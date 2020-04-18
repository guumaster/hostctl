package file

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/host/errors"
)

func TestFile_Enable(t *testing.T) {
	mem := CreateBasicFS(t)

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
		err = m.EnableOnly([]string{"profile2"})
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
		assert.EqualError(t, err, errors.ErrUnknownProfile.Error())
	})
}
