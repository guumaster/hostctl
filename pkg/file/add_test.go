package file

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/parser"
	"github.com/guumaster/hostctl/pkg/types"
)

func TestFile_AddProfile(t *testing.T) {
	mem := createBasicFS(t)
	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	t.Run("Add new", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)
		r := strings.NewReader(`127.0.0.1 added.loc`)

		p, err := parser.ParseProfile(r)
		assert.NoError(t, err)
		p.Name = "awesome"
		p.Status = types.Enabled

		err = m.AddProfile(p)
		assert.NoError(t, err)

		assert.Equal(t, []string{"profile1", "awesome"}, m.GetEnabled())

		added, err := m.GetProfile("awesome")
		assert.NoError(t, err)

		assert.Equal(t, added, p)
	})

	t.Run("Add existing", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)
		r := strings.NewReader(`127.0.0.1 added.loc`)

		p, err := parser.ParseProfile(r)
		assert.NoError(t, err)
		p.Name = "profile1"

		err = m.AddProfile(p)
		assert.NoError(t, err)

		assert.Equal(t, []string{"profile1"}, m.GetEnabled())

		added, err := m.GetProfile("profile1")
		assert.NoError(t, err)
		hosts, err := added.GetHostNames("127.0.0.1")
		assert.NoError(t, err)

		assert.Equal(t, hosts, []string{"first.loc", "second.loc", "added.loc"})
	})

	t.Run("Add default error", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)
		r := strings.NewReader(`127.0.0.1 added.loc`)

		p, err := parser.ParseProfile(r)
		assert.NoError(t, err)
		p.Name = "default"

		err = m.AddProfile(p)
		assert.Error(t, err, types.ErrDefaultProfile)
	})
}
