package host

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {

	t.Run("String", func(t *testing.T) {
		p := Profile{
			Name:   "awesome",
			Status: Enabled,
		}
		assert.Equal(t, "[on]awesome", p.String())
		p.Status = Disabled
		assert.Equal(t, "[off]awesome", p.String())
	})

	t.Run("Render", func(t *testing.T) {
		mem := createBasicFS(t)

		h, err := NewWithFs("/tmp/etc/hosts", mem)
		assert.NoError(t, err)

		b, err := mem.Create("memory")
		assert.NoError(t, err)

		p := h.data.Profiles["profile1"]
		err = p.Render(b)
		assert.NoError(t, err)

		c, err := afero.ReadFile(mem, b.Name())
		assert.NoError(t, err)

		assert.Contains(t, string(c), testEnabledProfile)
	})

	t.Run("AddRoute", func(t *testing.T) {
		r := strings.NewReader(`3.3.3.4 some.profile.loc`)
		p, err := NewProfileFromReader(r)
		assert.NoError(t, err)

		p.AddRoute("1.1.1.1", "added.loc")
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc"}, names)
	})

	t.Run("AddRoutes", func(t *testing.T) {
		r := strings.NewReader(`3.3.3.4 some.profile.loc`)
		p, err := NewProfileFromReader(r)
		assert.NoError(t, err)

		p.AddRoutes("1.1.1.1", []string{"added.loc", "another.loc"})
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc", "another.loc"}, names)
	})

	t.Run("RemoveRoutes", func(t *testing.T) {
		r := strings.NewReader("3.3.3.4 some.profile.loc\n5.5.5.5 another.profile.loc")
		p, err := NewProfileFromReader(r)
		assert.NoError(t, err)

		p.RemoveRoutes([]string{"another.profile.loc"})
		names, err := p.GetAllHostNames()
		assert.NoError(t, err)

		assert.Equal(t, []string{"some.profile.loc"}, names)
	})

	t.Run("GetHostnames", func(t *testing.T) {
		r := strings.NewReader(`3.3.3.4 some.profile.loc`)
		p, err := NewProfileFromReader(r)
		assert.NoError(t, err)

		names, err := p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)

		assert.Equal(t, []string{"some.profile.loc"}, names)

		_, err = p.GetHostNames("13333t")
		assert.EqualError(t, err, "invalid ip '13333t'")
	})
}
