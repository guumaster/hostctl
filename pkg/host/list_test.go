package host

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_List(t *testing.T) {
	mem := createBasicFS(t)

	f, err := mem.Open("/etc/hosts")
	assert.NoError(t, err)

	m, err := NewWithFs(f.Name(), mem)
	assert.NoError(t, err)

	t.Run("Table", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		m.List(&ListOptions{Writer: out})
		expected := `
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
		assertListOutput(t, out.String(), expected)
	})

	t.Run("Table Column order", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		m.List(&ListOptions{Writer: out, Columns: []string{"domain", "ip", "status"}})
		expected := `
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
		assertListOutput(t, out.String(), expected)
	})

	t.Run("Table Raw", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		m.List(&ListOptions{Writer: out, RawTable: true})

		expected := `
PROFILE 	STATUS	IP       	DOMAIN
default 	on    	127.0.0.1	localhost
profile1	on    	127.0.0.1	first.loc
profile1	on    	127.0.0.1	second.loc
profile2	off   	127.0.0.1	first.loc
profile2	off   	127.0.0.1	second.loc
`
		assertListOutput(t, out.String(), expected)
	})

	t.Run("Table Raw Filtered", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		m.List(&ListOptions{Writer: out, Profiles: []string{"profile1"}, RawTable: true})

		expected := `
PROFILE 	STATUS	IP       	DOMAIN
profile1	on    	127.0.0.1	first.loc
profile1	on    	127.0.0.1	second.loc
`
		assertListOutput(t, out.String(), expected)
	})
	t.Run("Table Raw Filtered", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		m.List(&ListOptions{
			Writer:   out,
			Columns:  []string{"ip", "domain"},
			Profiles: []string{"profile1"},
			RawTable: true,
		})

		expected := `
IP       	DOMAIN
127.0.0.1	first.loc
127.0.0.1	second.loc
`
		assertListOutput(t, expected, out.String())
	})

	t.Run("Profile status", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		m.ProfileStatus(&ListOptions{Writer: out})

		expected := `
+----------+--------+
| PROFILE  | STATUS |
+----------+--------+
| profile1 | on     |
| profile2 | off    |
+----------+--------+
`
		assertListOutput(t, out.String(), expected)
	})

	t.Run("Profiles status Raw", func(t *testing.T) {
		out := bytes.NewBufferString("\n")
		m.ProfileStatus(&ListOptions{Writer: out, RawTable: true})

		expected := `
PROFILE 	STATUS
profile1	on
profile2	off
`
		assertListOutput(t, out.String(), expected)
	})
}

func assertListOutput(t *testing.T, actual, expected string) {
	t.Helper()
	compact := regexp.MustCompile(`[ \t]+`)
	got := compact.ReplaceAllString(actual, "")
	want := compact.ReplaceAllString(expected, "")

	assert.Contains(t, want, got)
}
