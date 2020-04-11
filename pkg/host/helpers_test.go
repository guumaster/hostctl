package host

import (
	"net"
	"testing"

	"github.com/spf13/afero"
)

var localhost = net.ParseIP("127.0.0.1")

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

func createBasicFS(t *testing.T) afero.Fs {
	t.Helper()
	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/etc", 0755)
	f, _ := appFS.Create("/etc/hosts")
	defer f.Close()

	_, _ = f.WriteString(defaultProfile + banner + testEnabledProfile + testDisabledProfile)

	return appFS
}
