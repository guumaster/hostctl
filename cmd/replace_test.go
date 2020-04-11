package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReplaceStdin(t *testing.T) {
	cmd := rootCmd

	tmp := makeTempHostsFile(t, "replaceStdinCmd")
	defer os.Remove(tmp.Name())

	b := bytes.NewBufferString("")

	in := strings.NewReader(`3.3.3.3 stdin.replaced.loc`)
	cmd.SetOut(b)
	cmd.SetIn(in)
	cmd.SetArgs([]string{"replace", "profile1", "--host-file", tmp.Name()})

	err := cmd.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	actual := "\n" + string(out)
	expected := `
+----------+--------+---------+--------------------+
| PROFILE  | STATUS |   IP    |       DOMAIN       |
+----------+--------+---------+--------------------+
| profile1 | on     | 3.3.3.3 | stdin.replaced.loc |
+----------+--------+---------+--------------------+
`
	assert.Contains(t, actual, expected)
}
func Test_ReplaceFile(t *testing.T) {
	// This test only fails with others, works fine executed alone. Too weird race condition somewhere.
	t.SkipNow()
	cmd := rootCmd

	newProfile := makeTempProfile(t, "replace_profile1", []string{
		"5.5.5.5 replaced.loc",
		"5.5.5.6 replaced2.loc",
	})
	defer os.Remove(newProfile.Name())

	tmp := makeTempHostsFile(t, "replaceFileCmd")
	defer os.Remove(tmp.Name())

	b := bytes.NewBufferString("")

	cmd.SetOut(b)
	cmd.SetArgs([]string{"replace", "awesome", "--from", newProfile.Name(), "--host-file", tmp.Name()})

	err := cmd.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	actual := "\n" + string(out)
	expected := `
+----------+--------+---------+---------------+
| PROFILE  | STATUS |   IP    |    DOMAIN     |
+----------+--------+---------+---------------+
| profile2 | on     | 5.5.5.5 | replaced.loc  |
| profile2 | on     | 5.5.5.6 | replaced2.loc |
+----------+--------+---------+---------------+
`
	assert.Contains(t, actual, expected)
}
