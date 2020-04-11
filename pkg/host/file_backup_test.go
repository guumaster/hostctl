package host

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFile_Backup(t *testing.T) {

	mem := createBasicFS(t)
	h, err := NewWithFs("/etc/hosts", mem)
	assert.NoError(t, err)

	fname, err := h.Backup("/tmp")
	assert.NoError(t, err)

	want, err := afero.ReadFile(mem, "/etc/hosts")
	got, err := afero.ReadFile(mem, fname)

	assert.Equal(t, want, got)
}
