package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/render"
	"github.com/guumaster/hostctl/pkg/types"
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

func TestContainsDefault(t *testing.T) {
	err := containsDefault([]string{"default"})
	assert.EqualError(t, err, types.ErrDefaultProfile.Error())

	err = containsDefault([]string{"awesome"})
	assert.NoError(t, err)
}

func TestGetRenderer(t *testing.T) {
	t.Run("Markdown", func(t *testing.T) {
		cmd := NewRootCmd()
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.SetArgs([]string{"list", "--out", "md"})

		err := cmd.Execute()
		assert.NoError(t, err)

		r := getRenderer(cmd, nil)

		assert.IsType(t, render.TableRenderer{}, r)
		assert.IsType(t, render.Markdown, r.(render.TableRenderer).Type)

		cmd.SetOut(b)
		cmd.SetArgs([]string{"list", "--out", "markdown"})

		err = cmd.Execute()
		assert.NoError(t, err)

		r = getRenderer(cmd, nil)

		assert.IsType(t, render.TableRenderer{}, r)
		assert.IsType(t, render.Markdown, r.(render.TableRenderer).Type)
	})

	t.Run("Raw", func(t *testing.T) {
		cmd := NewRootCmd()
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.SetArgs([]string{"list", "--out", "raw"})

		err := cmd.Execute()
		assert.NoError(t, err)

		r := getRenderer(cmd, nil)

		assert.IsType(t, render.TableRenderer{}, r)
		assert.IsType(t, render.Raw, r.(render.TableRenderer).Type)
	})

	t.Run("JSON", func(t *testing.T) {
		cmd := NewRootCmd()
		b := bytes.NewBufferString("")
		cmd.SetOut(b)
		cmd.SetArgs([]string{"--out", "json"})

		err := cmd.Execute()
		assert.NoError(t, err)

		r := getRenderer(cmd, nil)

		assert.IsType(t, render.JSONRenderer{}, r)
		assert.IsType(t, render.JSON, r.(render.JSONRenderer).Type)
	})
}

func TestIsValidURL(t *testing.T) {
	valid := isValidURL("no valid")
	assert.Equal(t, valid, false)

	valid = isValidURL("/tmp/hosts")
	assert.Equal(t, valid, false)

	valid = isValidURL("http://localhost")
	assert.Equal(t, valid, true)

	valid = isValidURL("http://github.com/golang/go")
	assert.Equal(t, valid, true)
}

func TestReadFromURL(t *testing.T) {
	t.SkipNow()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	go func() {
		_ = http.ListenAndServe(":9998", nil)
	}()

	r, err := readerFromURL("http://0.0.0.0:9998/test")
	assert.NoError(t, err)

	c, _ := ioutil.ReadAll(r)

	assert.Equal(t, c, "Hello, test!")
}

func TestHelperCmd(t *testing.T) {
	info := newInfoCmd()
	helper := isHelperCmd(info)
	assert.True(t, helper)

	list := newListCmd()
	noHelper := isHelperCmd(list)
	assert.False(t, noHelper)
}
