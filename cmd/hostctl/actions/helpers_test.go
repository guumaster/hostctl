package actions

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/guumaster/hostctl/pkg/parser"
	"github.com/guumaster/hostctl/pkg/render"
	"github.com/guumaster/hostctl/pkg/types"
	"github.com/stretchr/testify/assert"
)

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

type handlerFn func(w http.ResponseWriter, r *http.Request)

type MyHandler struct {
	sync.Mutex
	fn handlerFn
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fn(w, r)
}

func TestReadFromURL(t *testing.T) {
	server := httptest.NewServer(&MyHandler{
		fn: func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`3.3.3.4 some.profile.loc`))
		},
	})
	defer server.Close()

	ctx := context.Background()
	r, err := readerFromURL(ctx, server.URL)
	assert.NoError(t, err)

	p, err := parser.ParseProfile(r)
	assert.NoError(t, err)

	hosts := p.GetAllHostNames()

	assert.Equal(t, []string{"some.profile.loc"}, hosts)
}

func TestHelperCmd(t *testing.T) {
	info := newInfoCmd()
	helper := isHelperCmd(info)
	assert.True(t, helper)

	list := newListCmd()
	noHelper := isHelperCmd(list)
	assert.False(t, noHelper)
}
