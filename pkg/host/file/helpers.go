package file

import (
	"net"
	"testing"

	"github.com/spf13/afero"

	"github.com/guumaster/hostctl/pkg/host/types"
)

var Localhost = net.ParseIP("127.0.0.1")

var DefaultProfile = "127.0.0.1 localhost\n"
var TestEnabledProfile = `
# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end
`
var TestDisabledProfile = `
# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# end
`

func CreateBasicFS(t *testing.T) afero.Fs {
	t.Helper()

	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/tmp/etc", 0755)

	f, _ := appFS.Create("/tmp/etc/hosts")
	defer f.Close()

	_, _ = f.WriteString(DefaultProfile + types.Banner + TestEnabledProfile + TestDisabledProfile)

	return appFS
}
