package docker

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseComposeFile(t *testing.T) {
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

	r := strings.NewReader(composeFile)
	containers, err := ParseComposeFile(r, "test")
	assert.NoError(t, err)

	sort.Strings(containers)

	assert.EqualValues(t, []string{"db", "test_container1", "test_container2"}, containers)
}
