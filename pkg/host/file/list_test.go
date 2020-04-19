package file

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guumaster/hostctl/pkg/host/render"
)

func TestFile_List(t *testing.T) {
	mem := createBasicFS(t)

	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	t.Run("Table", func(t *testing.T) {
		out := bytes.NewBufferString("\n")

		r := render.NewTableRenderer(&render.TableRendererOptions{Writer: out})

		m.List(r, &ListOptions{})

		const expected = `
+----------+--------+-----------+------------+
| PROFILE  | STATUS |    IP     |   DOMAIN   |
+----------+--------+-----------+------------+
| default  | on     | 127.0.0.1 | localhost  |
+----------+--------+-----------+------------+
| profile1 | on     | 127.0.0.1 | first.loc  |
| profile1 | on     | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
| profile2 | off    | 127.0.0.1 | first.loc  |
| profile2 | off    | 127.0.0.1 | second.loc |
+----------+--------+-----------+------------+
`
		assertListOutput(t, expected, out.String())
	})

	t.Run("Table Column order", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		tabOpts := &render.TableRendererOptions{
			Writer:  out,
			Columns: []string{"domain", "ip", "status"},
		}
		opts := &ListOptions{}

		r := render.NewTableRenderer(tabOpts)

		m.List(r, opts)

		const expected = `
+------------+-----------+--------+
| DOMAIN     | IP        | STATUS |
+------------+-----------+--------+
| localhost  | 127.0.0.1 | on     |
+------------+-----------+--------+
| first.loc  | 127.0.0.1 | on     |
| second.loc | 127.0.0.1 | on     |
+------------+-----------+--------+
| first.loc  | 127.0.0.1 | off    |
| second.loc | 127.0.0.1 | off    |
+------------+-----------+--------+
`
		assertListOutput(t, expected, out.String())
	})
}

func Test_TableRaw(t *testing.T) {
	mem := createBasicFS(t)

	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	t.Run("Table Raw", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		opts := &render.TableRendererOptions{Writer: out}

		r := render.NewRawRenderer(opts)

		m.List(r, &ListOptions{})

		const expected = `
PROFILE 	STATUS	IP       	DOMAIN
default 	on    	127.0.0.1	localhost
profile1	on    	127.0.0.1	first.loc
profile1	on    	127.0.0.1	second.loc
profile2	off   	127.0.0.1	first.loc
profile2	off   	127.0.0.1	second.loc
`
		assertListOutput(t, expected, out.String())
	})

	t.Run("Table Raw Filtered", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		opts := &render.TableRendererOptions{Writer: out}

		r := render.NewRawRenderer(opts)

		m.List(r, &ListOptions{Profiles: []string{"profile1"}})

		const expected = `
PROFILE 	STATUS	IP       	DOMAIN
profile1	on    	127.0.0.1	first.loc
profile1	on    	127.0.0.1	second.loc
`
		assertListOutput(t, expected, out.String())
	})

	t.Run("Table Raw Filtered with columns", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		opts := &render.TableRendererOptions{
			Writer:  out,
			Columns: []string{"ip", "domain"},
		}
		r := render.NewRawRenderer(opts)

		m.List(r, &ListOptions{
			Profiles: []string{"profile1"},
		})

		const expected = `
IP       	DOMAIN
127.0.0.1	first.loc
127.0.0.1	second.loc
`
		assertListOutput(t, expected, out.String())
	})
}

func Test_TableMarkdown(t *testing.T) {
	mem := createBasicFS(t)

	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	t.Run("Table Markdown", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		opts := &render.TableRendererOptions{Writer: out}

		r := render.NewMarkdownRenderer(opts)

		m.List(r, &ListOptions{})

		const expected = `
| PROFILE  | STATUS |    IP     |   DOMAIN   |
|----------|--------|-----------|------------|
| default  | on     | 127.0.0.1 | localhost  |
|----------|--------|-----------|------------|
| profile1 | on     | 127.0.0.1 | first.loc  |
| profile1 | on     | 127.0.0.1 | second.loc |
|----------|--------|-----------|------------|
| profile2 | off    | 127.0.0.1 | first.loc  |
| profile2 | off    | 127.0.0.1 | second.loc |
`
		assertListOutput(t, expected, out.String())
	})
}

func TestFile_ProfileStatus(t *testing.T) {
	mem := createBasicFS(t)

	f, err := mem.Open("/tmp/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	t.Run("Profile status", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		r := render.NewTableRenderer(&render.TableRendererOptions{
			Writer:  out,
			Columns: render.ProfilesOnlyColumns,
		})

		m.ProfileStatus(r, nil)

		const expected = `
+----------+--------+
| PROFILE  | STATUS |
+----------+--------+
| profile1 | on     |
| profile2 | off    |
+----------+--------+
`
		assertListOutput(t, expected, out.String())
	})

	t.Run("Profiles status Raw", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		r := render.NewRawRenderer(&render.TableRendererOptions{
			Writer:  out,
			Columns: render.ProfilesOnlyColumns,
		})

		m.ProfileStatus(r, nil)

		const expected = `
PROFILE 	STATUS
profile1	on
profile2	off
`
		assertListOutput(t, expected, out.String())
	})
}

func assertListOutput(t *testing.T, actual, expected string) {
	t.Helper()

	compact := regexp.MustCompile(`[ \t]+`)
	actual = compact.ReplaceAllString(actual, "")
	expected = compact.ReplaceAllString(expected, "")

	assert.Contains(t, expected, actual)
}
