package kubeconfig

import (
    "fmt"
    "os"

    "github.com/sirupsen/logrus"
    "k8s.io/client-go/kubernetes"
    _ "k8s.io/client-go/plugin/pkg/client/auth"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

// Logger setup
var log = logrus.New()

func init() {
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

// LoadKubeconfig attempts to load a kubeconfig based on default locations or an explicit path.
func LoadKubeconfig(path string) (*rest.Config, error) {
    // Use the default load order: KUBECONFIG env > $HOME/.kube/config > In cluster config
    loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
    if path != "" {
        loadingRules.ExplicitPath = path
    }
    loader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
    config, err := loader.ClientConfig()
    if err != nil {
        return nil, fmt.Errorf("loading kubeconfig: %w", err)
    }
    config.UserAgent = "k0pt"
    // Use protobuf for faster serialization instead of default JSON
    config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
    config.ContentType = "application/vnd.kubernetes.protobuf"
    return config, nil
}

// LoadKubeClient accepts a path to a kubeconfig to load and returns the clientset.
func LoadKubeClient(path string) (*kubernetes.Clientset, error) {
    config, err := LoadKubeconfig(path)
    if err != nil {
        log.WithError(err).Error("Failed to load kubeconfig")
        return nil, err
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.WithError(err).Error("Failed to create Kubernetes client")
        return nil, err
    }
    return clientset, nil
}