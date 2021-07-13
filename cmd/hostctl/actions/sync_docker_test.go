package actions

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/docker/docker/client"
	"github.com/guumaster/hostctl/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func testGetOptions(t *testing.T, cli *client.Client) getOptionsFn {
	t.Helper()

	return func(cmd *cobra.Command, profiles []string) (*docker.Options, error) {
		opts, err := getDockerOptions(cmd, nil)
		assert.NoError(t, err)

		opts.Cli = cli

		return opts, nil
	}
}

type transportFunc func(*http.Request) (*http.Response, error)

func (tf transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return tf(req)
}

func newClientWithResponse(t *testing.T, resp map[string]string) *client.Client {
	t.Helper()

	v := "1.22"
	c, err := client.NewClient("tcp://fake:2345", v,
		&http.Client{
			Transport: transportFunc(func(req *http.Request) (*http.Response, error) {
				url := req.URL.Path
				b, ok := resp[url]
				if !ok {
					b = "{}"
				}
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte(b))),
				}, nil
			}),
		},
		map[string]string{})

	assert.NoError(t, err)

	return c
}

func Test_SyncDocker(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "remove")
	defer r.Clean()

	cli := newClientWithResponse(t, map[string]string{
		"/v1.22/networks": `[
{"Id": "networkID1", "Name": "networkName1" }
]`,
		"/v1.22/containers/json": `[{
	"Id": "container_id-first", "Names": ["first"],
	"NetworkSettings": { "Networks": { "networkName2": { "NetworkID": "networkID1", "IPAddress": "172.0.0.2" }}}
}, {
	"Id": "container_id-second", "Names": ["second"],
	"NetworkSettings": { "Networks": { "networkName2": { "NetworkID": "networkID1", "IPAddress": "172.0.0.3" }}}
}]`,
	})

	opts := testGetOptions(t, cli)
	cmdSync := newSyncDockerCmd(nil, opts)
	cmdSync.Use = "test-sync-docker"

	cmd.AddCommand(cmdSync)

	r.Run("hostctl test-sync-docker profile2").
		Containsf(`
				[ℹ] Using hosts file: %s

				+----------+--------+-----------+------------+
				| PROFILE  | STATUS |    IP     |   DOMAIN   |
				+----------+--------+-----------+------------+
				| profile2 | on     | 172.0.0.2 | first.loc  |
				| profile2 | on     | 172.0.0.3 | second.loc |
				+----------+--------+-----------+------------+
			`, r.Hostfile())
}
