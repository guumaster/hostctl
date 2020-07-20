package minikube

// ProfileListResponse represents the json response of `minikube profile list`
type ProfileListResponse struct {
	Valid []struct {
		Config struct {
			Driver string `json:"Driver"`
			Name   string `json:"Name"`
			Nodes  []struct {
				IP   string `json:"IP"`
				Name string `json:"Name"`
			} `json:"Nodes"`
		} `json:"Config"`
		Name   string `json:"Name"`
		Status string `json:"Status"`
	} `json:"valid"`
}

// AddonsResponse represents the json response of `minikube addons list`
type AddonsResponse struct {
	Ingress struct {
		Status string `json:"Status"`
	} `json:"ingress"`
	IngressDNS struct {
		Status string `json:"Status"`
	} `json:"ingress-dns"`
}
