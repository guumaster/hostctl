package actions

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var composeFile = `
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

  another-thing:
    container_name: db
    image: db:1.0
    networks:
      - networkName1

networks:
  networkName1:
`

var dockerComposeResponse = map[string]string{
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
},{
	"Id": "container_id_db", "Names": ["/testing-app_db"],
	"NetworkSettings": {
		"Networks": { "testing-app_networkName1": { "NetworkID": "testing-app_networkID1", "IPAddress": "172.0.0.4" }}
  }
}]`,
}

func Test_SyncDockerCompose(t *testing.T) {
	c := makeTempHostsFile(t, "docker-compose.yml")

	_, _ = c.WriteString(composeFile)
	defer os.Remove(c.Name())

	cli := newClientWithResponse(t, dockerComposeResponse)

	cmd := NewRootCmd()

	opts := testGetOptions(t, cli)
	cmdSync := newSyncDockerCmd(nil, opts)
	cmdSync.Use = "test-sync-docker-compose"

	cmd.AddCommand(cmdSync)

	tmp := makeTempHostsFile(t, "syncDockerCmd")
	defer os.Remove(tmp.Name())

	b := bytes.NewBufferString("")

	cmd.SetOut(b)
	cmd.SetArgs([]string{"test-sync-docker-compose", "profile2", "--host-file", tmp.Name()})

	err := cmd.Execute()
	assert.NoError(t, err)

	out, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	actual := "\n" + string(out)

	const expected = `
+----------+--------+-----------+------------------------------+
| PROFILE  | STATUS |    IP     |            DOMAIN            |
+----------+--------+-----------+------------------------------+
| profile2 | on     | 172.0.0.2 | testing-app_container1_1.loc |
| profile2 | on     | 172.0.0.3 | testing-app_container2_1.loc |
| profile2 | on     | 172.0.0.4 | testing-app_db.loc           |
+----------+--------+-----------+------------------------------+
`

	assert.Contains(t, actual, expected)
}
