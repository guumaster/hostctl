package parser

import (
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/types"
)

func TestHostFile(t *testing.T) {
	testFile := `
127.0.0.1 localhost
##################################################################
# Content under this line is handled by hostctl. DO NOT EDIT.
##################################################################

# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end

# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# end
`

	t.Run("New", func(t *testing.T) {
		f := strings.NewReader(testFile)
		localhost := net.ParseIP("127.0.0.1")

		data, err := Parse(f)
		assert.NoError(t, err)
		assert.Equal(t, data.ProfileNames, []string{"profile1", "profile2"})
		assert.Equal(t, types.Enabled, data.Profiles["profile1"].Status)
		assert.Equal(t, types.Disabled, data.Profiles["profile2"].Status)
		assert.EqualValues(t, &types.Route{
			IP:        localhost,
			HostNames: []string{"first.loc", "second.loc"},
		}, data.Profiles["profile1"].Routes["127.0.0.1"])
	})
}

func appendLine(p *types.Profile, line string) {
	if line == "" {
		return
	}

	route, ok := parseRouteLine(line)
	if !ok {
		return
	}

	p.AddRoute(route)
}

func TestNewProfile(t *testing.T) {
	t.Run("parser.ParseProfile", func(t *testing.T) {
		r := strings.NewReader(`
3.3.3.4 some.profile.loc
3.3.3.4 first.loc
`)
		p, err := ParseProfile(r, true)
		assert.NoError(t, err)
		hosts, err := p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)
		assert.Equal(t, []string{"some.profile.loc", "first.loc"}, hosts)
	})

	t.Run("parser.ParseProfile non-uniq", func(t *testing.T) {
		r := strings.NewReader(`
3.3.3.4 some.profile.loc
# non-route-line
3.3.3.4 first.loc
3.3.3.4 first.loc
`)
		p, err := ParseProfile(r, false)
		assert.NoError(t, err)
		hosts, err := p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)
		assert.Equal(t, []string{"some.profile.loc", "first.loc", "first.loc"}, hosts)
	})
}

func TestParser(t *testing.T) {
	t.Run("appendLine enabled", func(t *testing.T) {
		p := &types.Profile{
			Name:   "test",
			Routes: map[string]*types.Route{},
		}
		appendLine(p, "127.0.0.1 first.loc")
		appendLine(p, "127.0.0.1 second.loc")
		assert.Len(t, p.Routes["127.0.0.1"].HostNames, 2)
	})

	t.Run("appendLine disabled", func(t *testing.T) {
		p := &types.Profile{
			Name:   "test",
			Routes: map[string]*types.Route{},
		}
		appendLine(p, "# 127.0.0.1 first.loc")
		assert.Len(t, p.Routes["127.0.0.1"].HostNames, 1)
	})

	t.Run("appendLine invalid lines", func(t *testing.T) {
		p := &types.Profile{
			Name:   "test",
			Routes: map[string]*types.Route{},
		}
		appendLine(p, "")
		appendLine(p, "3333 asdfasdfa")
		assert.Len(t, p.Routes, 0)
	})

	t.Run("parseProfileHeader", func(t *testing.T) {
		list := map[string]*types.Profile{
			"# profile.on something ": {
				Name:   "something",
				Status: types.Enabled,
			},
			"# profile.off thisoneoff ": {
				Name:   "thisoneoff",
				Status: types.Disabled,
			},
			"# profile another ": {
				Name:   "another",
				Status: types.Enabled,
			},
			"# profile another with spaces": {
				Name:   "another with spaces",
				Status: types.Enabled,
			},
			"wrong line": nil,
		}
		for header, wanted := range list {
			p, err := parseProfileHeader([]byte(header))
			if err != nil {
				assert.Equal(t, "invalid format for profile header", err.Error(), "unexpected error")
			}

			assert.Equal(t, wanted, p, "profile should match")
		}
	})
}
