package docker

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ComposeData represents data in a docker-compose.yml file
type composeData struct {
	Services map[string]composeService `yaml:"services"`
}

// ComposeService represents one service from a yml file
type composeService struct {
	ContainerName string `yaml:"container_name"`
}

// ParseComposeFile returns a list of containers from a docker-compose.yml file
func ParseComposeFile(r io.Reader, projectName string) ([]string, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	data := &composeData{}

	err = yaml.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	var containers []string

	for serv, data := range data.Services {
		name := data.ContainerName
		if data.ContainerName == "" {
			name = fmt.Sprintf("%s_%s", projectName, serv)
		}

		containers = append(containers, name)
	}

	return containers, nil
}
