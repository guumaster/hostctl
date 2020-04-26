package file

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func makeTempHostsFile(t *testing.T, pattern string) *os.File {
	t.Helper()

	file, err := ioutil.TempFile("/tmp", pattern+"_")
	if err != nil {
		t.Fatal(err)
	}

	_, _ = file.WriteString(`
127.0.0.1 localhost
# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end
# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# end
`)
	defer file.Close()

	return file
}

func createBasicFS(t *testing.T) afero.Fs {
	t.Helper()

	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/tmp/etc", 0755)

	f, _ := appFS.Create("/tmp/etc/hosts")
	defer f.Close()

	_, _ = f.WriteString(`
127.0.0.1 localhost
` + Banner + `
# profile.on profile1
127.0.0.1 first.loc
127.0.0.1 second.loc
# end
# profile.off profile2
# 127.0.0.1 first.loc
# 127.0.0.1 second.loc
# end
`)

	return appFS
}
