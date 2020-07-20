package minikube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
)

// Profile contains information about a minikube profile
type Profile struct {
	Name   string
	IP     string
	Status Status
	Driver string

	IngressEnabled    bool
	IngressDNSEnabled bool
}

// Status represents the current state of a profile
type Status string

const (
	// Running status of minikube profile
	Running Status = "Running"
)

// GetProfile returns information about a minikube profile
func GetProfile(name string) (*Profile, error) {
	mini, err := exec.LookPath("minikube")
	if err != nil {
		return nil, err
	}

	b := bytes.NewBufferString("")
	c := exec.Command(mini, "profile", "list", "-o", "json")
	c.Stdout = b

	err = c.Run()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	var m ProfileListResponse

	_ = json.Unmarshal(data, &m)

	var profile *Profile

	for _, p := range m.Valid {
		if p.Name == name {
			profile = &Profile{
				Name:   p.Name,
				IP:     p.Config.Nodes[0].IP,
				Status: Status(p.Status),
				Driver: p.Config.Driver,
			}
		}
	}

	if profile == nil {
		return nil, fmt.Errorf("can't find profile '%s' on minikube", name)
	}

	b = bytes.NewBufferString("")
	c = exec.Command(mini, "addons", "list", "-o", "json", "-p", profile.Name)
	c.Stdout = b

	err = c.Run()
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	var a *AddonsResponse

	_ = json.Unmarshal(data, &a)

	if a == nil {
		return profile, nil
	}

	profile.IngressDNSEnabled = a.IngressDNS.Status == "enabled"
	profile.IngressEnabled = a.Ingress.Status == "enabled"

	return profile, nil
}
