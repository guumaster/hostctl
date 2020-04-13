package host

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_RemoveProfile(t *testing.T) {
	mem := createBasicFS(t)
	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	t.Run("Remove", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)

		err = m.RemoveProfile("profile2")
		assert.NoError(t, err)

		assert.Equal(t, []string{"profile1"}, m.GetEnabled())
		assert.Equal(t, []string{}, m.GetDisabled())

		_, err = m.GetProfile("profile2")
		assert.EqualError(t, err, UnknownProfileError.Error())
	})

	t.Run("Remove unknown", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)
		err = m.RemoveProfile("unknown")
		assert.EqualError(t, err, UnknownProfileError.Error())
	})

	t.Run("Remove profiles", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)
		err = m.RemoveProfiles([]string{"profile1", "profile2"})
		assert.NoError(t, err)

		assert.Equal(t, []string{}, m.GetEnabled())
	})

}
