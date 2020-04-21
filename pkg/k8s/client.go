package k8s

import (
	"os"
	"path/filepath"

	"github.com/guumaster/cligger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClientset returns a new clienset to interact with Kubernetes
func NewClientset() (*kubernetes.Clientset, error) {
	kubeConfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	if fromEnv := os.Getenv("KUBECONFIG"); fromEnv != "" {
		kubeConfigPath = fromEnv
		cligger.Info("Using config from %s", kubeConfigPath)
	}

	// use the current context in kubeConfigPath
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, cligger.Errorf("fatal error kubernetes config: %s", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, cligger.Errorf("error with kubernetes config:", err)
	}

	return clientset, nil
}
