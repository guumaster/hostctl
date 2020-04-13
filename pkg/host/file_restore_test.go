package host

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_Restore(t *testing.T) {
	mem := createBasicFS(t)

	h, err := NewWithFs("/tmp/etc/hosts", mem)
	assert.NoError(t, err)

	_ = mem.Mkdir("/tmp", 0755)
	backup, err := mem.Create("/tmp/TestFile_Restore")
	assert.NoError(t, err)

	_, err = io.Copy(backup, h.src)
	assert.NoError(t, err)

	err = h.Restore(backup.Name())
	assert.NoError(t, err)

}
