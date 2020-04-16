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
		cmd.SetArgs([]string{"add", "awesome", "--uniq", "--from", tmp.Name(), "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		const expected = `
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
		cmd.SetArgs([]string{"add", "awesome", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		const expected = `
`
		assert.Contains(t, actual, expected)
	})
}

func Test_ReplaceStdin(t *testing.T) {
	cmd := rootCmd

	tmp := makeTempHostsFile(t, "replaceStdinCmd")
	defer os.Remove(tmp.Name())

	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	in := strings.NewReader(`3.3.3.3 stdin.replaced.loc`)
	cmd.SetIn(in)

	cmd.SetArgs([]string{"replace", "profile1", "--host-file", tmp.Name()})

	err := cmd.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	actual := "\n" + string(out)
	assert.Contains(t, actual, `
+----------+--------+---------+--------------------+
| PROFILE  | STATUS |   IP    |       DOMAIN       |
+----------+--------+---------+--------------------+
| profile1 | on     | 3.3.3.3 | stdin.replaced.loc |
+----------+--------+---------+--------------------+
`)
}

func Test_ReplaceFile(t *testing.T) {
	cmd := rootCmd

	in := strings.NewReader(`
5.5.5.5 replaced.loc
5.5.5.6 replaced2.loc
`)

	tmp := makeTempHostsFile(t, "replaceFileCmd")
	defer os.Remove(tmp.Name())

	b := bytes.NewBufferString("")

	cmd.SetIn(in)
	cmd.SetOut(b)
	cmd.SetArgs([]string{"replace", "awesome", "--host-file", tmp.Name()})

	err := cmd.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	actual := "\n" + string(out)
	assert.Contains(t, actual, `
+---------+--------+---------+---------------+
| PROFILE | STATUS |   IP    |    DOMAIN     |
+---------+--------+---------+---------------+
| awesome | on     | 5.5.5.5 | replaced.loc  |
| awesome | on     | 5.5.5.6 | replaced2.loc |
+---------+--------+---------+---------------+
`)
}
