package actions

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Info(t *testing.T) {
	cmd := NewRootCmd()

	b := bytes.NewBufferString("")

	cmd.SetOut(b)
	cmd.SetArgs([]string{"info"})

	err := cmd.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	actual := "\n" + string(out)
	assert.Contains(t, actual, "Your dev tool to manage /etc/hosts like a pro")
}
