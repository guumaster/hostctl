package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddDomains(t *testing.T) {
	cmd := rootCmd

	t.Run("Add domains", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "addDomainCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"add", "domains", "profile1", "arg.domain.loc", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		const expected = `
+----------+--------+-----------+----------------+
| PROFILE  | STATUS |    IP     |     DOMAIN     |
+----------+--------+-----------+----------------+
| profile1 | on     | 127.0.0.1 | first.loc      |
| profile1 | on     | 127.0.0.1 | second.loc     |
| profile1 | on     | 127.0.0.1 | arg.domain.loc |
+----------+--------+-----------+----------------+
`
		assert.Contains(t, actual, expected)
	})

	t.Run("Add domains new profile", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "addDomainCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"add", "domains", "newprofile", "arg.domain.loc", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		const expected = `
+------------+--------+-----------+----------------+
|  PROFILE   | STATUS |    IP     |     DOMAIN     |
+------------+--------+-----------+----------------+
| newprofile | on     | 127.0.0.1 | arg.domain.loc |
+------------+--------+-----------+----------------+
`
		assert.Contains(t, actual, expected)
	})

	t.Run("Add domains with IP", func(t *testing.T) {
		tmp := makeTempHostsFile(t, "addDomainCmd")
		defer os.Remove(tmp.Name())
		b := bytes.NewBufferString("")

		cmd.SetOut(b)
		cmd.SetArgs([]string{"add", "domains", "profile1", "--ip", "5.5.5.5", "arg2.domain.loc", "--host-file", tmp.Name()})

		err := cmd.Execute()
		assert.NoError(t, err)

		out, err := ioutil.ReadAll(b)
		assert.NoError(t, err)

		actual := "\n" + string(out)
		const expected = `
+----------+--------+-----------+-----------------+
| PROFILE  | STATUS |    IP     |     DOMAIN      |
+----------+--------+-----------+-----------------+
| profile1 | on     | 127.0.0.1 | first.loc       |
| profile1 | on     | 127.0.0.1 | second.loc      |
| profile1 | on     | 5.5.5.5   | arg2.domain.loc |
+----------+--------+-----------+-----------------+
`
		assert.Contains(t, actual, expected)
	})
}
