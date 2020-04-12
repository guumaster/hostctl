package host

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_ReplaceProfile(t *testing.T) {
	mem := createBasicFS(t)
	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	t.Run("Replace", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)

		r := strings.NewReader(`4.4.4.4 replaced.loc`)

		p, err := NewProfileFromReader(r)
		assert.NoError(t, err)
		p.Name = "profile1"
		p.Status = Enabled

		err = m.ReplaceProfile(*p)
		assert.NoError(t, err)

		replaced, err := m.GetProfile("profile1")
		assert.NoError(t, err)
		hosts, err := replaced.GetHostNames("4.4.4.4")
		assert.Equal(t, hosts, []string{"replaced.loc"})
	})

	t.Run("Replace new", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)

		r := strings.NewReader(`4.4.4.4 replaced.loc`)

		p, err := NewProfileFromReader(r)
		assert.NoError(t, err)
		p.Name = "awesome"
		p.Status = Enabled

		err = m.ReplaceProfile(*p)
		assert.NoError(t, err)

		added, err := m.GetProfile("awesome")
		assert.NoError(t, err)
		hosts, err := added.GetHostNames("4.4.4.4")
		assert.Equal(t, hosts, []string{"replaced.loc"})
	})
}
