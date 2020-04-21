package profile

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/types"
)

func TestProfile(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		p := types.Profile{
			Name:   "awesome",
			Status: types.Enabled,
		}
		assert.Equal(t, "[on]awesome", p.String())
		p.Status = types.Disabled
		assert.Equal(t, "[off]awesome", p.String())
	})

	t.Run("Render", func(t *testing.T) {
		testEnabledProfile := `
# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end
`
		r := strings.NewReader(testEnabledProfile)
		p, err := NewProfileFromReader(r, true)
		assert.NoError(t, err)

		p.Name = "profile1"
		p.Status = types.Enabled
		b := bytes.NewBufferString("")

		err = p.Render(b)
		assert.NoError(t, err)

		c, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		assert.Contains(t, string(c), testEnabledProfile)
	})

	t.Run("AddRoute empty profile", func(t *testing.T) {
		r := strings.NewReader(``)
		p, err := NewProfileFromReader(r, true)
		assert.NoError(t, err)

		p.AddRoute("1.1.1.1", "added.loc")
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc"}, names)
	})

	t.Run("AddRoute same IP", func(t *testing.T) {
		r := strings.NewReader(``)
		p, err := NewProfileFromReader(r, true)
		assert.NoError(t, err)

		p.AddRoute("1.1.1.1", "added.loc")
		p.AddRoute("1.1.1.1", "second.loc")
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc", "second.loc"}, names)
	})

	t.Run("AddRoute", func(t *testing.T) {
		r := strings.NewReader(`3.3.3.4 some.profile.loc`)
		p, err := NewProfileFromReader(r, true)
		assert.NoError(t, err)

		p.AddRoute("1.1.1.1", "added.loc")
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc"}, names)
	})

	t.Run("AddRoutes", func(t *testing.T) {
		r := strings.NewReader(`3.3.3.4 some.profile.loc`)
		p, err := NewProfileFromReader(r, true)
		assert.NoError(t, err)

		p.AddRoutes("1.1.1.1", []string{"added.loc", "another.loc"})
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc", "another.loc"}, names)
	})

	t.Run("RemoveRoutes", func(t *testing.T) {
		r := strings.NewReader("3.3.3.4 some.profile.loc\n5.5.5.5 another.profile.loc")
		p, err := NewProfileFromReader(r, true)
		assert.NoError(t, err)

		p.RemoveRoutes([]string{"another.profile.loc"})
		names, err := p.GetAllHostNames()
		assert.NoError(t, err)

		assert.Equal(t, []string{"some.profile.loc"}, names)
	})

	t.Run("GetHostnames", func(t *testing.T) {
		r := strings.NewReader(`3.3.3.4 some.profile.loc`)
		p, err := NewProfileFromReader(r, true)
		assert.NoError(t, err)

		names, err := p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)

		assert.Equal(t, []string{"some.profile.loc"}, names)

		_, err = p.GetHostNames("13333t")
		assert.EqualError(t, err, "invalid ip '13333t'")
	})
}

func Test_appendIP(t *testing.T) {
	r := strings.NewReader(`3.3.3.4 some.profile.loc`)
	p, err := NewProfileFromReader(r, true)
	assert.NoError(t, err)

	p.AppendIP("3.3.3.4")
	p.AppendIP("3.3.3.4")
	p.AppendIP("3.3.3.5")

	assert.Equal(t, p.IPList, []string{"3.3.3.4", "3.3.3.5"})
}

func TestDefaultProfile_Render(t *testing.T) {
	b := bytes.NewBufferString("")

	d := types.DefaultProfile{}

	d = append(d, &types.Row{
		Comment: "# This is a comment",
	}, &types.Row{
		IP:   "127.0.0.1",
		Host: "localhost",
	})

	err := d.Render(b)
	assert.NoError(t, err)

	c, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	out := "\n" + string(c)
	assert.Equal(t, `
# This is a comment
127.0.0.1 localhost
`, out)
}
