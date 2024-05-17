package namespaces

import (
    "context"
    "fmt"

    "github.com/mmarconi93/k0pt/pkg/callbacks"
    "github.com/sirupsen/logrus"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// Logger setup
var log = logrus.New()

func init() {
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

// DeleteNamespace deletes a namespace in the Kubernetes cluster
func DeleteNamespace(clientset *kubernetes.Clientset, namespace string) {
    if namespace == "" {
        log.Error("Namespace argument is required")
        return
    }

    // Attempt to delete the namespace
    err := clientset.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
    if err != nil {
        callbacks.PrintErrorMessage(err)
        log.WithError(err).Errorf("Failed to delete namespace: %s", namespace)
        return
    }

    log.Infof("Successfully deleted namespace: %s", namespace)
    fmt.Printf("❌ Deleted namespace => %s\n", namespace)
}
