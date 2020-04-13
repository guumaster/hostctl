package host

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {

	t.Run("Get Status", func(t *testing.T) {
		mem := createBasicFS(t)

		m, err := NewWithFs("/tmp/etc/hosts", mem)
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

	t.Run("AddRoutes", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		r := strings.NewReader(`3.3.3.4 some.profile.loc`)
		p, err := NewProfileFromReader(r)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		hosts, err := p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)

		err = f.AddRoutes("profile2", "3.3.3.4", hosts)
		assert.NoError(t, err)

		err = f.Flush()
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)

		assert.Contains(t, string(c), defaultProfile)
		assert.Contains(t, string(c), banner)
		assert.Contains(t, string(c), testEnabledProfile)
		var added = `
# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# 3.3.3.4 some.profile.loc
# end
`
		assert.Contains(t, string(c), added)
	})

	t.Run("RemoveRoutes", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		removed, err := f.RemoveRoutes("profile2", []string{"second.loc"})
		assert.NoError(t, err)
		assert.Equal(t, false, removed)

		err = f.Flush()
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)
		assert.Contains(t, string(c), defaultProfile)
		assert.Contains(t, string(c), banner)
		assert.Contains(t, string(c), testEnabledProfile)
		var added = `
# profile.off profile2
# 127.0.0.1 first.loc
# end
`
		assert.Contains(t, string(c), added)
	})

	t.Run("RemoveRoutes", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		removed, err := f.RemoveRoutes("profile2", []string{"first.loc", "second.loc"})
		assert.NoError(t, err)
		assert.Equal(t, true, removed)

		err = f.Flush()
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)
		assert.Contains(t, string(c), defaultProfile)
		assert.Contains(t, string(c), banner)
		assert.Contains(t, string(c), testEnabledProfile)
		assert.NotContains(t, string(c), testDisabledProfile)
	})

	t.Run("WriteToFile", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		err = f.writeToFile(h)
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)
		assert.Contains(t, string(c), defaultProfile)
		assert.Contains(t, string(c), banner)
		assert.Contains(t, string(c), testEnabledProfile)
		assert.Contains(t, string(c), testDisabledProfile)
	})

	t.Run("writeBanner", func(t *testing.T) {
		mem := createBasicFS(t)
		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		f.writeBanner(h)
		h.Close()

		content, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)

		assert.Contains(t, string(content), banner)
	})
}
