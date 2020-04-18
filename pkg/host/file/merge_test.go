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

	profiles := []*types.Profile{
		{
			Name:   "profile2",
			Status: types.Enabled,
			Routes: map[string]*types.Route{
				ip3.String(): {IP: ip3, HostNames: []string{"third.new.loc"}},
			},
		},
		{
			Name:   "profile3",
			Status: types.Enabled,
			Routes: map[string]*types.Route{
				ip4.String(): {IP: ip4, HostNames: []string{"third.new.loc", "fourth.new.loc"}},
			},
		},
	}
	m.MergeProfiles(profiles)

	assert.Equal(t, []string{"profile1", "profile3"}, m.GetEnabled())
	assert.Equal(t, []string{"profile2"}, m.GetDisabled())

	p3, err := m.GetProfile("profile3")
	assert.NoError(t, err)
	assert.Equal(t, profiles[1], p3)

	p2, err := m.GetProfile("profile2")
	assert.NoError(t, err)

	modP2 := profiles[0]
	modP2.IPList = []string{"127.0.0.1", "2.2.2.2"}
	modP2.Routes[Localhost.String()] = &types.Route{IP: Localhost, HostNames: []string{"first.loc", "second.loc"}}
	modP2.Status = types.Disabled
	assert.Equal(t, modP2, p2)
}
