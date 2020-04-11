package host

import (
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

		h, err := NewWithFs("/etc/hosts", mem)
		assert.NoError(t, err)

		b, err := mem.Create("memory")

		p := h.data.Profiles["profile1"]
		err = p.Render(b)
		assert.NoError(t, err)

		c, err := afero.ReadFile(mem, b.Name())
		assert.NoError(t, err)

		assert.Contains(t, string(c), testEnabledProfile)
	})
}
