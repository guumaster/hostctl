package host

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProfile(t *testing.T) {
	t.Run("NewProfileFromReader", func(t *testing.T) {
		r := strings.NewReader(`
3.3.3.4 some.profile.loc
3.3.3.4 first.loc
`)
		p, err := NewProfileFromReader(r, true)
		assert.NoError(t, err)
		hosts, err := p.GetHostNames("3.3.3.4")
		assert.NoError(t, err)
		assert.Equal(t, []string{"some.profile.loc", "first.loc"}, hosts)
	})
}
