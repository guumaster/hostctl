package host

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Disable(t *testing.T) {
	mem := createBasicFS(t)

	f, err := mem.Open("/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	t.Run("Disable one", func(t *testing.T) {
		err = m.Disable([]string{"profile1"})
		assert.NoError(t, err)
		assert.Contains(t, m.GetDisabled(), "profile1")
	})

	t.Run("Disable All", func(t *testing.T) {
		err = m.DisableAll()
		assert.NoError(t, err)
		disabled := m.GetDisabled()
		assert.Contains(t, disabled, "profile1")
		assert.Contains(t, disabled, "profile2")
	})

	t.Run("Disable error", func(t *testing.T) {
		err = m.Disable([]string{"unknown"})
		assert.EqualError(t, err, UnknownProfileError.Error())
	})
}
