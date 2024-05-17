package main

import (
    "context"
    "fmt"
    "os"

    "github.com/mmarconi93/k0pt/pkg/analysis"
    "github.com/mmarconi93/k0pt/pkg/control"
    "github.com/mmarconi93/k0pt/pkg/info"
    "github.com/mmarconi93/k0pt/pkg/namespaces"
    "github.com/mmarconi93/k0pt/pkg/recommendations"
    "github.com/sirupsen/logrus"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

var log = logrus.New()

func init() {
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
    command, namespace, clusterName := control.ParseFlags()
    fmt.Printf("Command: %s, Namespace: %s, ClusterName: %s\n", command, namespace, clusterName)

    clientset, config := control.CreateClientset()

    switch command {
    case "list-pods":
        listPods(clientset, namespace)
    case "get-cluster-state":
        control.GetClusterState(clusterName)
    case "get-admin-credentials":
        control.GetAdminCredentials(clusterName)
    case "delete-namespace":
        namespaces.DeleteNamespace(clientset, namespace)
    case "switch-namespace":
        namespaces.CheckoutNamespace(clientset, config, namespace)
    case "calculate-cost-savings":
        recommendations.CalculateCostSavings(clientset)
    case "optimize-resources":
        recommendations.OptimizeResources(clientset)
    case "analyze-resource-usage":
        analysis.AnalyzeResourceUsage()
    case "help":
        info.DisplayHelp()
    case "version":
        info.VersionInfo()
    default:
        fmt.Println("Unknown command. Available commands: list-pods, get-cluster-state, get-admin-credentials, delete-namespace, switch-namespace, calculate-cost-savings, optimize-resources, analyze-resource-usage, help, version")
    }
}

func listPods(clientset *kubernetes.Clientset, namespace string) {
    pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
    if err != nil {
        log.WithError(err).Error("Failed to list pods")
        return
    }

    for _, pod := range pods.Items {
        fmt.Printf("Pod: %s\n", pod.Name)
    }
}