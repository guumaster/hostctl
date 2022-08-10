package render

import (
	"bytes"
	"io"
	"testing"

	"github.com/guumaster/hostctl/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMarkdownRenderer(t *testing.T) {
	b := bytes.NewBufferString("")
	r := NewMarkdownRenderer(&TableRendererOptions{
		Writer:  b,
		Columns: []string{"profile", "status", "ip"},
	})

	r.AppendRow(&types.Row{
		Comment: "this is not rendered",
	})
	r.AppendRow(&types.Row{
		Profile: "awesome",
		Status:  "on",
		IP:      "127.0.0.1",
		Host:    "awesome.loc",
	})

	err := r.Render()
	assert.NoError(t, err)

	out, err := io.ReadAll(b)
	assert.NoError(t, err)

	expected := `
| PROFILE | STATUS |    IP     |
|---------|--------|-----------|
| awesome | on     | 127.0.0.1 |
`

	actual := "\n" + string(out)

	assert.Equal(t, expected, actual)
}
