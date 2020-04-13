package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RemoveDomains(t *testing.T) {
	cmd := rootCmd

	t.Run("Remove domains", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "removeDomainCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"remove", "domains", "profile1", "first.loc", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := `
+----------+--------+-----------+------------+
| PROFILE  | STATUS |    IP     |   DOMAIN   |
+----------+--------+-----------+------------+
| profile1 | on     | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
		assert.Contains(t, actual, expected)
	})

	t.Run("Remove domains and profile", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "removeDomainCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"remove", "domains", "profile1", "first.loc", "second.loc", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		// TODO: Improve list management to avoid printing empty lists
		expected := `
+---------+--------+----+--------+
| PROFILE | STATUS | IP | DOMAIN |
+---------+--------+----+--------+
+---------+--------+----+--------+
`
		assert.Contains(t, actual, expected)
	})

}
