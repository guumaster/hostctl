package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Restore(t *testing.T) {
	cmd := rootCmd

	from := makeTempHostsFile(t, "restoreFrom")
	defer os.Remove(from.Name())

	to, err := ioutil.TempFile("/tmp", "restoreTo")
	assert.NoError(t, err)
	defer os.Remove(to.Name())

	b := bytes.NewBufferString("")

	cmd.SetOut(b)
	cmd.SetArgs([]string{"restore", "--from", from.Name(), "--host-file", to.Name()})

	err = cmd.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	toData, _ := ioutil.ReadFile(to.Name())
	fromData, _ := ioutil.ReadFile(from.Name())
	assert.Equal(t, string(toData), string(fromData))

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
}
