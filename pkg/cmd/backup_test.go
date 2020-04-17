package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Backup(t *testing.T) {
	cmd := rootCmd

	tmp := makeTempHostsFile(t, "backupCmd")
	defer os.Remove(tmp.Name())

	b := bytes.NewBufferString("")

	cmd.SetOut(b)
	cmd.SetArgs([]string{"backup", "--path", "/tmp", "--host-file", tmp.Name()})

	err := cmd.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	actual := "\n" + string(out)
	assert.Contains(t, actual, `
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
`)
}
