package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Status(t *testing.T) {
	cmd := rootCmd

	tmp := makeTempHostsFile(t, "statusCmd")
	defer os.Remove(tmp.Name())

	t.Run("Status", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"status", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		expected := `
+----------+--------+
| PROFILE  | STATUS |
+----------+--------+
| profile1 | on     |
| profile2 | off    |
+----------+--------+
`
		assert.Contains(t, actual, expected)
	})
}
