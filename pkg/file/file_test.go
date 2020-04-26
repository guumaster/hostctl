package file

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/profile"
	"github.com/guumaster/hostctl/pkg/types"
)

// nolint:gochecknoglobals
var (
	defaultProfile = "127.0.0.1 localhost\n"

	testEnabledProfile = `
# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end
`
	testDisabledProfile = `
# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# end
`
	onlyEnabled  = defaultProfile + Banner + "\n" + testEnabledProfile
	fullHostfile = onlyEnabled + testDisabledProfile
)

func TestNewWithFs(t *testing.T) {
	t.Run("With file", func(t *testing.T) {
		src := makeTempHostsFile(t, "etc_hosts")
		fs := afero.NewOsFs()

		m, err := NewWithFs(src.Name(), fs)
		assert.NoError(t, err)

		err = m.Flush()
		assert.NoError(t, err)

		assert.Equal(t, src.Name(), m.src.Name())
	})
	t.Run("With invalid file", func(t *testing.T) {
		fs := afero.NewOsFs()

		m, err := NewWithFs("/tmp/invalid_random_name", fs)
		assert.Nil(t, m)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Without fs", func(t *testing.T) {
		src := makeTempHostsFile(t, "etc_hosts")

		m, err := NewWithFs(src.Name(), nil)
		assert.NoError(t, err)

		err = m.Flush()
		assert.NoError(t, err)

		assert.Equal(t, src.Name(), m.src.Name())
	})
}

func TestNewWithFsError(t *testing.T) {
}

func TestManagerStatus(t *testing.T) {
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
			expected := map[string]types.Status{
				"profile1": types.Enabled,
				"profile2": types.Disabled,
			}
			assert.Equal(t, expected, actual)
		})

		t.Run("GetStatus unknown", func(t *testing.T) {
			actual := m.GetStatus([]string{"unknown"})
			expected := map[string]types.Status{}
			assert.Equal(t, expected, actual)
		})
	})
}

func TestManagerRoutes(t *testing.T) {
	t.Run("AddRoutes", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		r := strings.NewReader(`3.3.3.4 some.profile.loc`)
		p, err := profile.NewProfileFromReader(r, true)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		hosts, err := p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)

		err = f.AddRoute("profile2", types.NewRoute("3.3.3.4", hosts...))
		assert.NoError(t, err)

		err = f.Flush()
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)

		assert.Contains(t, string(c), onlyEnabled)
		var added = `
# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# 3.3.3.4 some.profile.loc
# end
`
		assert.Contains(t, string(c), added)
	})

	t.Run("AddRoutes new", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		err = f.AddRoute("awesome", types.NewRoute("3.3.3.4", "host1.loc", "host2.loc"))
		assert.NoError(t, err)

		err = f.Flush()
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)

		assert.Contains(t, string(c), onlyEnabled)
		var added = `
# profile.on awesome
3.3.3.4 host1.loc
3.3.3.4 host2.loc
# end
`
		assert.Contains(t, string(c), added)
	})

	t.Run("RemoveHostnames one", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		removed, err := f.RemoveHostnames("profile2", []string{"second.loc"})
		assert.NoError(t, err)
		assert.Equal(t, false, removed)

		err = f.Flush()
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)
		assert.Contains(t, string(c), onlyEnabled)

		var added = `
# profile.off profile2
# 127.0.0.1 first.loc
# end
`
		assert.Contains(t, string(c), added)
	})

	t.Run("RemoveHostnames multi", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		removed, err := f.RemoveHostnames("profile2", []string{"first.loc", "second.loc"})
		assert.NoError(t, err)
		assert.Equal(t, true, removed)

		err = f.Flush()
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)
		assert.Contains(t, string(c), onlyEnabled)
		assert.NotContains(t, string(c), testDisabledProfile)
	})
}

func TestManagerWrite(t *testing.T) {
	t.Run("writeToFile", func(t *testing.T) {
		mem := createBasicFS(t)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		err = f.writeToFile(h)
		assert.NoError(t, err)
		f.Close()

		c, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)
		assert.Contains(t, string(c), fullHostfile)
	})

	t.Run("WriteTo", func(t *testing.T) {
		f := makeTempHostsFile(t, "WriteTo")
		defer os.Remove(f.Name())

		h, err := NewFile(f.Name())
		assert.NoError(t, err)

		err = h.WriteTo(f.Name())
		assert.NoError(t, err)
		f.Close()

		c, err := ioutil.ReadFile(f.Name())
		assert.NoError(t, err)
		assert.Contains(t, string(c), fullHostfile)
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

		assert.Contains(t, string(content), Banner)
	})

	t.Run("writeBanner once", func(t *testing.T) {
		mem := createBasicFS(t)
		h, _ := mem.OpenFile("/tmp/etc/hosts", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)

		f, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		f.writeBanner(h)
		f.writeBanner(h)
		h.Close()

		content, err := afero.ReadFile(mem, h.Name())
		assert.NoError(t, err)

		assert.Contains(t, string(content), Banner)
	})
}
