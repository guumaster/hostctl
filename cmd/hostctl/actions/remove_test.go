package actions

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/types"
)

func Test_Remove(t *testing.T) {
	cmd := NewRootCmd()

	t.Run("Remove", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "addCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"remove", "profile2", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := listHeader
		assert.NotContains(t, expected, actual)
	})

	t.Run("Remove multiple", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "addCmd")
		defer os.Remove(tmp.Name())

		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"remove", "profile1", "profile2", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := listHeader
		assert.NotContains(t, actual, expected)
	})

	t.Run("Remove unknown", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "addCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"remove", "unknown", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.EqualError(t, err, types.ErrUnknownProfile.Error())
	})

	t.Run("Remove all", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "addCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"remove", "--all", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := listHeader
		assert.NotContains(t, actual, expected)
	})

	t.Run("Remove all bad", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "addCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"remove", "profile1", "--all", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.EqualError(t, err, "args must be empty with --all flag")
	})
}
