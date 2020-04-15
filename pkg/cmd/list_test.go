package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_List(t *testing.T) {
	cmd := rootCmd

	tmp := makeTempHostsFile(t, "listCmd")
	defer os.Remove(tmp.Name())

	t.Run("List all", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"list", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := `
+----------+--------+-----------+------------+
| PROFILE  | STATUS |    IP     |   DOMAIN   |
+----------+--------+-----------+------------+
| default  | on     | 127.0.0.1 | localhost  |
+----------+--------+-----------+------------+
| profile1 | on     | 127.0.0.1 | first.loc  |
| profile1 | on     | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
| profile2 | off    | 127.0.0.1 | first.loc  |
| profile2 | off    | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
		assert.Contains(t, actual, expected)
	})
	for _, filter := range []string{"enabled", "disabled"} {
		t.Run("List "+filter, func(t *testing.T) {
			b := bytes.NewBufferString("")

			cmd.SetOut(b)
			cmd.SetArgs([]string{"list", filter, "--host-file", tmp.Name()})

			err := cmd.Execute()
			assert.NoError(t, err)

			out, err := ioutil.ReadAll(b)
			assert.NoError(t, err)

			expected := ""
			actual := "\n" + string(out)
			if filter == "enabled" {
				expected = `
+----------+--------+-----------+------------+
| PROFILE  | STATUS |    IP     |   DOMAIN   |
+----------+--------+-----------+------------+
| default  | on     | 127.0.0.1 | localhost  |
+----------+--------+-----------+------------+
| profile1 | on     | 127.0.0.1 | first.loc  |
| profile1 | on     | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
			} else {
				expected = `
+----------+--------+-----------+------------+
| PROFILE  | STATUS |    IP     |   DOMAIN   |
+----------+--------+-----------+------------+
| default  | on     | 127.0.0.1 | localhost  |
+----------+--------+-----------+------------+
| profile2 | off    | 127.0.0.1 | first.loc  |
| profile2 | off    | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
			}
			assert.Contains(t, actual, expected)
		})

	}
}
