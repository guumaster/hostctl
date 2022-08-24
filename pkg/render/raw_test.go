package render

import (
	"bytes"
	"io"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/types"
)

func TestNewRawRenderer(t *testing.T) {
	b := bytes.NewBufferString("")
	r := NewRawRenderer(&TableRendererOptions{
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
PROFILE	STATUS	IP
awesome	on    	127.0.0.1
`

	actual := "\n" + string(out)

	compact := regexp.MustCompile(`[ \t]+`)
	actual = compact.ReplaceAllString(actual, "")
	expected = compact.ReplaceAllString(expected, "")

	assert.Contains(t, expected, actual)
}
