package profile

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestNewProfileFromDocker(t *testing.T) {
	t.Run("All containers", func(t *testing.T) {
		c := newClientWithResponse(t, map[string]string{
			"/v1.22/networks": `[
{"Id": "networkId1", "Name": "networkName1" },
{"Id": "networkId2", "Name": "networkName2" }
]`,
			"/v1.22/containers/json": `[{
	"Id": "container_id1", "Names": ["container1"],
	"NetworkSettings": { "Networks": { "networkName1": { "NetworkID": "networkID1", "IPAddress": "172.0.0.2" }}}
}, {
	"Id": "container_id2", "Names": ["container2"],
	"NetworkSettings": { "Networks": { "networkName1": { "NetworkID": "networkID1", "IPAddress": "172.0.0.3" }}}
}]`,
		})

		p, err := NewProfileFromDocker(&DockerOptions{
			Domain: "test",
			Cli:    c,
		})

		assert.NoError(t, err)

		hosts := p.GetAllHostNames()

		assert.Equal(t, []string{"172.0.0.2", "172.0.0.3"}, p.IPList)
		assert.Equal(t, []string{"container1.test", "container2.test"}, hosts)
	})

	t.Run("Only one network", func(t *testing.T) {
		c := newClientWithResponse(t, map[string]string{
			"/v1.22/networks": `[
{"Id": "networkID2", "Name": "networkName2" }
]`,
			"/v1.22/containers/json": `[{
	"Id": "container_id2", "Names": ["container2"],
	"NetworkSettings": { "Networks": { "networkName2": { "NetworkID": "networkID2", "IPAddress": "172.0.0.3" }}}
}]`,
		})

		p, err := NewProfileFromDocker(&DockerOptions{
			Domain:  "test",
			Cli:     c,
			Network: "networkName2",
		})

		assert.NoError(t, err)

		hosts := p.GetAllHostNames()

		assert.Equal(t, []string{"172.0.0.3"}, p.IPList)
		assert.Equal(t, []string{"container2.test"}, hosts)
	})
}

func TestNewProfileFromDockerCompose(t *testing.T) {
	t.Run("Compose all", func(t *testing.T) {
		r := strings.NewReader(`
version: "3"
services:

  container1:
    image: some_image:3.5
    networks:
      - networkName1

  container2:
    image: some_other_image:1.0
    networks:
      - networkName1

networks:
  networkName1:
`)
		c := newClientWithResponse(t, map[string]string{
			"/v1.22/networks": `[
{"Id": "testing-app_networkID1", "Name": "testing-app_networkName1" }
]`,
			"/v1.22/containers/json": `[{
	"Id": "container_id1", "Names": ["/testing-app_container1_1"],
	"NetworkSettings": {
		"Networks": { "testing-app_networkName1": { "NetworkID": "testing-app_networkID1", "IPAddress": "172.0.0.2" }}
  }
},{
	"Id": "container_id2", "Names": ["/testing-app_container2_1"],
	"NetworkSettings": {
		"Networks": { "testing-app_networkName1": { "NetworkID": "testing-app_networkID1", "IPAddress": "172.0.0.3" }}
  }
}]`,
		})

		p, err := NewProfileFromDockerCompose(&DockerOptions{
			Domain:      "loc",
			Cli:         c,
			ComposeFile: r,
			ProjectName: "testing-app",
		})

		assert.NoError(t, err)

		hosts := p.GetAllHostNames()

		assert.Equal(t, []string{"172.0.0.2", "172.0.0.3"}, p.IPList)
		assert.Equal(t, []string{"container1_1.loc", "container2_1.loc"}, hosts)
	})
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
