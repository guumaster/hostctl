package file

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFile_Backup(t *testing.T) {
	mem := CreateBasicFS(t)
	h, err := NewWithFs("/tmp/etc/hosts", mem)
	assert.NoError(t, err)

	fname, err := h.Backup("/tmp")
	assert.NoError(t, err)

	want, err := afero.ReadFile(mem, "/tmp/etc/hosts")
	assert.NoError(t, err)
	got, err := afero.ReadFile(mem, fname)
	assert.NoError(t, err)

	assert.Equal(t, want, got)
}
