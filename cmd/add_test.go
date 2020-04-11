package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	cmd := rootCmd

	tmp := makeTempHostsFile(t, "addCmd")
	defer os.Remove(tmp.Name())

	t.Run("Add from file", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"add-to", "awesome", "--from", tmp.Name(), "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := `
+---------+--------+-----------+------------+
| PROFILE | STATUS |    IP     |   DOMAIN   |
+---------+--------+-----------+------------+
| awesome | on     | 127.0.0.1 | localhost  |
| awesome | on     | 127.0.0.1 | first.loc  |
| awesome | on     | 127.0.0.1 | second.loc |
+---------+--------+-----------+------------+
`
		assert.Contains(t, actual, expected)
	})

	t.Run("Add from stdin", func(t *testing.T) {
		b := bytes.NewBufferString("")

		in := strings.NewReader(`3.3.3.3 stdin.loc`)
		cmd.SetOut(b)
		cmd.SetIn(in)
		cmd.SetArgs([]string{"add-to", "awesome", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := `
`
		assert.Contains(t, actual, expected)
	})

}
