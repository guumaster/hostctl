package file

import (
	"strings"
	"testing"

	"github.com/guumaster/hostctl/pkg/parser"
	"github.com/guumaster/hostctl/pkg/types"
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

		p, err := parser.ParseProfile(r)
		assert.NoError(t, err)
		p.Name = "profile1"
		p.Status = types.Enabled

		err = m.ReplaceProfile(p)
		assert.NoError(t, err)

		replaced, err := m.GetProfile("profile1")
		assert.NoError(t, err)
		hosts, err := replaced.GetHostNames("4.4.4.4")
		assert.NoError(t, err)
		assert.Equal(t, hosts, []string{"replaced.loc"})
	})

	t.Run("Replace new", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)

		r := strings.NewReader(`4.4.4.4 replaced.loc`)

		p, err := parser.ParseProfile(r)
		assert.NoError(t, err)
		p.Name = "awesome"
		p.Status = types.Enabled

		err = m.ReplaceProfile(p)
		assert.NoError(t, err)

		added, err := m.GetProfile("awesome")
		assert.NoError(t, err)
		hosts, err := added.GetHostNames("4.4.4.4")
		assert.NoError(t, err)
		assert.Equal(t, hosts, []string{"replaced.loc"})
	})

	t.Run("Replace unknown", func(t *testing.T) {
		m, err := NewWithFs(f.Name(), mem)
		assert.NoError(t, err)

		r := strings.NewReader(`4.4.4.4 replaced.loc`)

		p, err := parser.ParseProfile(r)
		assert.NoError(t, err)
		p.Name = types.Default
		p.Status = types.Enabled

		err = m.ReplaceProfile(p)
		assert.EqualError(t, err, types.ErrDefaultProfile.Error())
	})
}
