package docker

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestGetContainerList(t *testing.T) {
	cli := newClientWithResponse(t, map[string]string{
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

	list, err := GetContainerList(context.Background(), cli, "")
	assert.NoError(t, err)

	assert.Len(t, list, 2)

	assert.IsType(t, types.Container{}, list[0])
	assert.IsType(t, types.Container{}, list[1])
	assert.Equal(t, "networkID1", list[1].NetworkSettings.Networks["networkName1"].NetworkID)

	assert.Equal(t, types.Container{
		ID:    "container_id1",
		Names: []string{"container1"},

		NetworkSettings: list[0].NetworkSettings, // simplify the comparison
	}, list[0])
	assert.Equal(t, types.Container{
		ID:    "container_id2",
		Names: []string{"container2"},

		NetworkSettings: list[1].NetworkSettings, // simplify the comparison
	}, list[1])
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
