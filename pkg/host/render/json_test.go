package render

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewJSONRenderer(t *testing.T) {
	b := bytes.NewBufferString("")
	r := NewJSONRenderer(&JSONRendererOptions{
		Writer:  b,
		Columns: []string{"profile", "status", "ip"},
	})

	r.AppendRow(&Row{
		Comment: "this is not rendered",
	})
	r.AppendRow(&Row{
		Profile: "awesome",
		Status:  "on",
		IP:      "127.0.0.1",
		Host:    "awesome.loc",
	})

	err := r.Render()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	expected := `
[{"Profile":"awesome","Status":"on","IP":"127.0.0.1","Host":"awesome.loc"}]
`

	actual := "\n" + string(out)

	assert.Equal(t, expected, actual)
}
