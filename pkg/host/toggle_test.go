package host

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Toggle(t *testing.T) {
	mem := createBasicFS(t)

	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	t.Run("Toggle", func(t *testing.T) {
		err = m.Toggle([]string{"profile1", "profile2"})
		assert.NoError(t, err)
		assert.Contains(t, m.GetEnabled(), "profile2")
		assert.Contains(t, m.GetDisabled(), "profile1")
	})

	t.Run("Toggle error", func(t *testing.T) {
		err = m.Toggle([]string{"unknown"})
		assert.EqualError(t, err, ErrUnknownProfile.Error())
	})
}
