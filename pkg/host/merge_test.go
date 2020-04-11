package host

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_MergeProfiles(t *testing.T) {
	mem := createBasicFS(t)

	f, err := mem.Open("/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	ip3 := net.ParseIP("2.2.2.2")
	ip4 := net.ParseIP("3.3.3.3")

	c := &Content{
		DefaultProfile: nil,
		ProfileNames:   []string{"profile2", "profile3"},
		Profiles: map[string]Profile{
			"profile2": {
				Name:   "profile2",
				Status: Enabled,
				Routes: map[string]*Route{
					ip3.String(): {IP: ip3, HostNames: []string{"third.new.loc"}},
				},
			},
			"profile3": {
				Name:   "profile3",
				Status: Enabled,
				Routes: map[string]*Route{
					ip4.String(): {IP: ip4, HostNames: []string{"third.new.loc", "fourth.new.loc"}},
				},
			},
		},
	}
	m.MergeProfiles(c)

	assert.Equal(t, []string{"profile1", "profile3"}, m.GetEnabled())
	assert.Equal(t, []string{"profile2"}, m.GetDisabled())

	p3, err := m.GetProfile("profile3")
	assert.NoError(t, err)
	assert.Equal(t, c.Profiles["profile3"], *p3)

	p2, err := m.GetProfile("profile2")
	assert.NoError(t, err)

	modP2 := c.Profiles["profile2"]
	modP2.Routes[localhost.String()] = &Route{IP: localhost, HostNames: []string{"first.loc", "second.loc"}}
	modP2.Status = Disabled
	assert.Equal(t, modP2, *p2)
}
