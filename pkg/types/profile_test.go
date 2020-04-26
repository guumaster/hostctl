package types

import (
	"bytes"
	"io/ioutil"
	"testing"

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
		testEnabledProfile := `
# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end
`

		p := Profile{
			Name:   "profile1",
			Status: Enabled,
		}
		p.AddRoute(NewRoute("127.0.0.1", "first.loc", "second.loc"))

		b := bytes.NewBufferString("")

		err := p.Render(b)
		assert.NoError(t, err)

		c, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		assert.Contains(t, string(c), testEnabledProfile)
	})

	t.Run("AddRoute empty profile", func(t *testing.T) {
		p := Profile{}

		p.AddRoute(NewRoute("1.1.1.1", "added.loc"))
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc"}, names)
	})

	t.Run("AddRoute same IP", func(t *testing.T) {
		p := Profile{}

		p.AddRoute(NewRoute("1.1.1.1", "added.loc"))
		p.AddRoute(NewRoute("1.1.1.1", "second.loc"))
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc", "second.loc"}, names)
	})

	t.Run("AddRoute", func(t *testing.T) {
		p := Profile{
			IPList: []string{"3.3.3.4"},
			Routes: map[string]*Route{
				"3.3.3.4": NewRoute("3.3.3.4", "some.profile.loc"),
			},
		}

		p.AddRoute(NewRoute("1.1.1.1", "added.loc"))
		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"added.loc"}, names)

		names, err = p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)

		assert.Equal(t, []string{"some.profile.loc"}, names)
	})

	t.Run("AddRouteUniq", func(t *testing.T) {
		p := Profile{}

		p.AddRoute(NewRoute("1.1.1.1", "another.loc", "added.loc"))
		p.AddRouteUniq(NewRoute("1.1.1.1", "added.loc", "third.loc"))

		names, err := p.GetHostNames("1.1.1.1")
		assert.NoError(t, err)

		assert.Equal(t, []string{"another.loc", "added.loc", "third.loc"}, names)
	})

	t.Run("RemoveHostnames", func(t *testing.T) {
		p := Profile{
			IPList: []string{"3.3.3.4"},
			Routes: map[string]*Route{
				"3.3.3.4": NewRoute("3.3.3.4", "another.profile.loc", "some.profile.loc"),
			},
		}

		p.RemoveHostnames([]string{"another.profile.loc"})
		names := p.GetAllHostNames()

		assert.Equal(t, []string{"some.profile.loc"}, names)
	})

	t.Run("GetHostnames", func(t *testing.T) {
		p := Profile{}

		p.AddRoute(NewRoute("3.3.3.4", "some.profile.loc"))

		names, err := p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)

		assert.Equal(t, []string{"some.profile.loc"}, names)

		_, err = p.GetHostNames("13333t")
		assert.EqualError(t, err, "invalid ip '13333t'")
	})
}

func Test_appendIP(t *testing.T) {
	p := Profile{}

	p.AddRoute(NewRoute("3.3.3.4", "some.profile.loc"))

	p.appendIP("3.3.3.4")
	p.appendIP("3.3.3.4")
	p.appendIP("3.3.3.5")

	assert.Equal(t, p.IPList, []string{"3.3.3.4", "3.3.3.5"})
}

func TestDefaultProfile_Render(t *testing.T) {
	b := bytes.NewBufferString("")

	d := DefaultProfile{}

	d = append(d, &Row{
		Comment: "# This is a comment",
	}, &Row{
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
