package actions

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/types"
)

func Test_Toggle(t *testing.T) {
	cmd := NewRootCmd()

	tmp := makeTempHostsFile(t, "enableCmd")
	defer os.Remove(tmp.Name())

	t.Run("Toggle", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"toggle", "profile2", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		const expected = `
+----------+--------+-----------+------------+
| PROFILE  | STATUS |    IP     |   DOMAIN   |
+----------+--------+-----------+------------+
| profile2 | on     | 127.0.0.1 | first.loc  |
| profile2 | on     | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
		assert.Contains(t, actual, expected)
	})

	t.Run("Toggle unknown", func(t *testing.T) {
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"toggle", "unknown", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.EqualError(t, err, types.ErrUnknownProfile.Error())
	})
}
