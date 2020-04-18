package file

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/host/types"
)

func TestFile_MergeProfiles(t *testing.T) {
	mem := CreateBasicFS(t)

	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	ip3 := net.ParseIP("2.2.2.2")
	ip4 := net.ParseIP("3.3.3.3")

	c := &types.Content{
		DefaultProfile: nil,
		ProfileNames:   []string{"profile2", "profile3"},
		Profiles: map[string]*types.Profile{
			"profile2": {
				Name:   "profile2",
				Status: types.Enabled,
				Routes: map[string]*types.Route{
					ip3.String(): {IP: ip3, HostNames: []string{"third.new.loc"}},
				},
			},
			"profile3": {
				Name:   "profile3",
				Status: types.Enabled,
				Routes: map[string]*types.Route{
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
	assert.Equal(t, c.Profiles["profile3"], p3)

	p2, err := m.GetProfile("profile2")
	assert.NoError(t, err)

	modP2 := c.Profiles["profile2"]
	modP2.IPList = []string{"127.0.0.1", "2.2.2.2"}
	modP2.Routes[Localhost.String()] = &types.Route{IP: Localhost, HostNames: []string{"first.loc", "second.loc"}}
	modP2.Status = types.Disabled
	assert.Equal(t, modP2, p2)
}
