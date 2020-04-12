package host

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostFile(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		appFS := createBasicFS(t)

		f, _ := appFS.Open("/tmp/etc/hosts")
		data, err := Parse(f)
		assert.NoError(t, err)
		assert.Equal(t, data.ProfileNames, []string{"profile1", "profile2"})
		assert.Equal(t, Enabled, data.Profiles["profile1"].Status)
		assert.Equal(t, Disabled, data.Profiles["profile2"].Status)
		assert.EqualValues(t, &Route{
			localhost,
			[]string{"first.loc", "second.loc"},
		}, data.Profiles["profile1"].Routes["127.0.0.1"])
	})
}

func TestParser(t *testing.T) {
	t.Run("appendLine enabled", func(t *testing.T) {
		p := &Profile{
			Name:   "test",
			Routes: map[string]*Route{},
		}
		appendLine(p, "127.0.0.1 first.loc")
		appendLine(p, "127.0.0.1 second.loc")
		assert.Len(t, p.Routes["127.0.0.1"].HostNames, 2)
	})

	t.Run("appendLine disabled", func(t *testing.T) {
		p := &Profile{
			Name:   "test",
			Routes: map[string]*Route{},
		}
		appendLine(p, "# 127.0.0.1 first.loc")
		assert.Len(t, p.Routes["127.0.0.1"].HostNames, 1)
	})

	t.Run("appendLine invalid lines", func(t *testing.T) {
		p := &Profile{
			Name:   "test",
			Routes: map[string]*Route{},
		}
		appendLine(p, "")
		appendLine(p, "3333 asdfasdfa")
		assert.Len(t, p.Routes, 0)
	})

	t.Run("parseProfileHeader", func(t *testing.T) {
		list := map[string]*Profile{
			"# profile.on something ": {
				Name:   "something",
				Status: Enabled,
			},
			"# profile.off thisoneoff ": {
				Name:   "thisoneoff",
				Status: Disabled,
			},
			"# profile another ": {
				Name:   "another",
				Status: Enabled,
			},
			"# profile another with spaces": {
				Name:   "another with spaces",
				Status: Enabled,
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

	t.Run("cleanLines", func(t *testing.T) {
		l := cleanLine("#     127.0.0.1 dirty.line        ")
		assert.Equal(t, "# 127.0.0.1 dirty.line", l)

		l = cleanLine("127.0.0.1\t\t dirty.loc \t second.loc      ")
		assert.Equal(t, "127.0.0.1 dirty.loc second.loc", l)
	})
}
