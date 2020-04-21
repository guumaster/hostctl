package file

import (
	"io/ioutil"
	"net"
	"os"
	"testing"

	"github.com/spf13/afero"
)

var Localhost = net.ParseIP("127.0.0.1")

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

func createBasicFS(t *testing.T) afero.Fs {
	t.Helper()

	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/tmp/etc", 0755)

	f, _ := appFS.Create("/tmp/etc/hosts")
	defer f.Close()

	_, _ = f.WriteString(defaultProfile + Banner + testEnabledProfile + testDisabledProfile)

	return appFS
}
