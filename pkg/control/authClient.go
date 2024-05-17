package control

import (
    "flag"
    "os"

    "github.com/sirupsen/logrus"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Logger setup
var log = logrus.New()

func init() {
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

// ParseFlags parses command line flags
func ParseFlags() (string, string, string) {
    flag.Parse()
    command := flag.Arg(0)
    namespace := flag.Arg(1)
    clusterName := flag.Arg(2)
    return command, namespace, clusterName
}

// CreateClientset creates a Kubernetes clientset
func CreateClientset() (*kubernetes.Clientset, *clientcmdapi.Config) {
    kubeconfig := clientcmd.NewDefaultPathOptions()
    kubeconfigFile := kubeconfig.GetDefaultFilename()

    // Allow kubeconfig file to be specified via an environment variable
    if envKubeconfig := os.Getenv("KUBECONFIG"); envKubeconfig != "" {
        kubeconfigFile = envKubeconfig
    }

    // Load the kubeconfig file
    config, err := kubeconfig.GetStartingConfig()
    if err != nil {
        log.WithError(err).Fatal("Failed to load kubeconfig")
    }

    // Create a *rest.Config object
    restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigFile)
    if err != nil {
        log.WithError(err).Fatal("Failed to build config")
    }

    // Create the clientset
    clientset, err := kubernetes.NewForConfig(restConfig)
    if err != nil {
        log.WithError(err).Fatal("Failed to create clientset")
    }

    return clientset, config
}
