package k8s

import (
	"net"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Ingress contains information about an ingress rule
type Ingress struct {
	IP       net.IP
	Hostname string
}

// GetIngresses returns a list of ingresses presents on a namespace
func GetIngresses(cli *kubernetes.Clientset, ns string) ([]Ingress, error) {
	if ns == "" {
		ns = v1.NamespaceAll
	}

	ing, err := cli.ExtensionsV1beta1().Ingresses(ns).List(v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var list []Ingress

	for _, i := range ing.Items {
		ip := i.Status.LoadBalancer.Ingress[0].IP

		for _, r := range i.Spec.Rules {
			list = append(list, Ingress{
				IP:       net.ParseIP(ip),
				Hostname: r.Host,
			})
		}
	}

	return list, nil
}
