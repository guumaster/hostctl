package host

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProfile(t *testing.T) {

	t.Run("NewProfileFromReader", func(t *testing.T) {
		mem := createBasicFS(t)
		f, err := mem.Open("/etc/hosts")
		assert.NoError(t, err)

		p, err := NewProfileFromReader(f)
		assert.NoError(t, err)
		assert.Equal(t, []string{"localhost", "first.loc", "second.loc"}, p.Routes[localhost.String()].HostNames)
	})

}
