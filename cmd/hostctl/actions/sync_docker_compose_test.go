package actions

import (
	"os"
	"testing"

	"github.com/docker/docker/client"
)

func Test_SyncDockerCompose(t *testing.T) {
	cmd := NewRootCmd()

	r := NewRunner(t, cmd, "remove")
	defer r.Clean()

	cli := prepareComposeCli(t, r)

	opts := testGetOptions(t, cli)
	cmdSync := newSyncDockerCmd(nil, opts)
	cmdSync.Use = "test-sync-docker-compose"

	cmd.AddCommand(cmdSync)

	r.Run("hostctl test-sync-docker-compose profile2").
		Containsf(`
				[ℹ] Using hosts file: %s

        +----------+--------+-----------+------------------------------+
        | PROFILE  | STATUS |    IP     |            DOMAIN            |
        +----------+--------+-----------+------------------------------+
        | profile2 | on     | 172.0.0.2 | testing-app_container1_1.loc |
        | profile2 | on     | 172.0.0.3 | testing-app_container2_1.loc |
        | profile2 | on     | 172.0.0.4 | testing-app_db.loc           |
        +----------+--------+-----------+------------------------------+
			`, r.Hostfile())
}

func prepareComposeCli(t *testing.T, r Runner) *client.Client {
	c := r.TempHostfile("docker-compose.yml")
	defer os.Remove(c.Name())

	_, _ = c.WriteString(`
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
`)

	dockerResponse := map[string]string{
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

	return newClientWithResponse(t, dockerResponse)
}
