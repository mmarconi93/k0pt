package main

import (
    "context"
    "fmt"
    "os"

    "github.com/joho/godotenv"
    "github.com/mmarconi93/k0pt/pkg/analysis"
    "github.com/mmarconi93/k0pt/pkg/control"
    "github.com/mmarconi93/k0pt/pkg/info"
    "github.com/mmarconi93/k0pt/pkg/kubeconfig"
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
    if err := godotenv.Load(); err != nil {
        log.Warn("No .env file found")
    }
}

func main() {
    command, namespace, clusterName := control.ParseFlags()
    fmt.Printf("Command: %s, Namespace: %s, ClusterName: %s\n", command, namespace, clusterName)

    clientset, err := kubeconfig.LoadKubeClient("")
    if err != nil {
        log.Fatalf("Failed to load Kubernetes client: %v", err)
    }

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
        config, err := kubeconfig.LoadKubeconfig("")
        if err != nil {
            log.Fatalf("Failed to load kubeconfig: %v", err)
        }
        namespaces.CheckoutNamespace(clientset, config, namespace)
    case "calculate-cost-savings":
        recommendations.CalculateCostSavings(clientset)
    case "optimize-resources":
        recommendations.OptimizeResources(clientset)
    case "analyze-resource-usage":
        analysis.AnalyzeResourceUsage(clientset)
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