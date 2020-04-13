package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

var defaultProfile = "127.0.0.1 localhost\n"

var testEnabledProfile = `
# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end
`

var testDisabledProfile = `
# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# end
`

var listHeader = `
+---------+--------+-----------+-----------+
| PROFILE | STATUS |    IP     |  DOMAIN   |
+---------+--------+-----------+-----------+
`

func makeTempHostsFile(t *testing.T, pattern string) *os.File {
	t.Helper()

	file, err := ioutil.TempFile("/tmp", pattern+"_")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = file.WriteString(defaultProfile + testEnabledProfile + testDisabledProfile)
	defer file.Close()

	return file
}
