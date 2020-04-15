package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/host"
)

func Test_Disable(t *testing.T) {
	cmd := rootCmd

	tmp := makeTempHostsFile(t, "disableCmd")
	defer os.Remove(tmp.Name())

	t.Run("Disable", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"disable", "profile1", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := `
+----------+--------+-----------+------------+
| PROFILE  | STATUS |    IP     |   DOMAIN   |
+----------+--------+-----------+------------+
| profile1 | off    | 127.0.0.1 | first.loc  |
| profile1 | off    | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
		assert.Contains(t, actual, expected)
	})

	t.Run("Disable unknown", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"disable", "unknown", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.EqualError(t, err, host.UnknownProfileError.Error())
	})

	t.Run("Disable All", func(t *testing.T) {
		t.SkipNow()
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"disable", "--all", "--host-file", tmp.Name()})

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
| profile1 | off    | 127.0.0.1 | first.loc  |
| profile1 | off    | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
| profile2 | off    | 127.0.0.1 | first.loc  |
| profile2 | off    | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
		assert.Contains(t, actual, expected)
	})

	t.Run("Disable Only", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"disable", "profile1", "--only", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		cmd.SetArgs([]string{"list", "--host-file", tmp.Name()})

		err = cmd.Execute()
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
| profile1 | off    | 127.0.0.1 | first.loc  |
| profile1 | off    | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
| profile2 | on     | 127.0.0.1 | first.loc  |
| profile2 | on     | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
		assert.Contains(t, actual, expected)
	})

	t.Run("Disable all error", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"disable", "any", "--all", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.EqualError(t, err, "args must be empty with --all flag")
	})

}
