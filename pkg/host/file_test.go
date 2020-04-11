package host

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {

	t.Run("Get Status", func(t *testing.T) {
		mem := createBasicFS(t)

		m, err := NewWithFs("/etc/hosts", mem)
		assert.NoError(t, err)

		t.Run("GetEnabled", func(t *testing.T) {
			enabled := m.GetEnabled()
			assert.Equal(t, []string{"profile1"}, enabled)
		})

		t.Run("GetDisabled", func(t *testing.T) {
			disabled := m.GetDisabled()
			assert.Equal(t, []string{"profile2"}, disabled)
		})
		t.Run("GetStatus", func(t *testing.T) {
			actual := m.GetStatus([]string{"profile1", "profile2"})
			expected := map[string]ProfileStatus{
				"profile1": Enabled,
				"profile2": Disabled,
			}
			assert.Equal(t, expected, actual)
		})
	})

	t.Run("WriteToFile", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/etc/hosts", mem)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		err = f.writeToFile(h)
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.Contains(t, string(c), defaultProfile)
		assert.Contains(t, string(c), banner)
		assert.Contains(t, string(c), testEnabledProfile)
		assert.Contains(t, string(c), testDisabledProfile)
	})

	t.Run("writeBanner", func(t *testing.T) {
		mem := createBasicFS(t)
		h, _ := mem.OpenFile("/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		f, err := NewWithFs("/etc/hosts", mem)
		assert.NoError(t, err)

		f.writeBanner(h)
		content, _ := ioutil.ReadFile(h.Name())

		assert.Contains(t, string(content), banner)
	})
}
