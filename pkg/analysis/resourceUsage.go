package analysis

import (
    "context"
    "fmt"
    "os"

    "github.com/mmarconi93/k0pt/pkg/cloud"
    "github.com/sirupsen/logrus"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/metrics/pkg/client/clientset/versioned"
)

// Logger setup
var log = logrus.New()

func init() {
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

// PodResourceUsage holds the resource usage details of a pod
type PodResourceUsage struct {
    Namespace   string
    PodName     string
    CPUUsage    string
    MemoryUsage string
    CPUCost     float64
    MemoryCost  float64
}

// AnalyzeResourceUsage analyzes the resource usage of all pods in the cluster
func AnalyzeResourceUsage() {
    log.Info("Analyzing resource usage...")

    config, err := rest.InClusterConfig()
    if err != nil {
        log.WithError(err).Error("Failed to get in-cluster config")
        return
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.WithError(err).Error("Failed to create Kubernetes client")
        return
    }

    metricsClient, err := versioned.NewForConfig(config)
    if err != nil {
        log.WithError(err).Error("Failed to create metrics client")
        return
    }

    // Fetch Azure pricing
    prices, err := cloud.FetchAzurePricing()
    if err != nil {
        log.WithError(err).Error("Failed to fetch Azure pricing")
        return
    }

    pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.WithError(err).Error("Failed to list pods")
        return
    }

    log.Infof("There are %d pods in the cluster", len(pods.Items))

    var resourceUsages []PodResourceUsage

    for _, pod := range pods.Items {
        podName := pod.Name
        namespace := pod.Namespace

        usage, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
        if err != nil {
            log.WithError(err).Errorf("Failed to get metrics for pod: %s", podName)
            continue
        }

        for _, container := range usage.Containers {
            cpuUsage := container.Usage.Cpu().MilliValue()
            memoryUsage := container.Usage.Memory().Value() / (1024 * 1024) // in MiB

            cpuCost := float64(cpuUsage) / 1000 * prices["cpu"]
            memoryCost := float64(memoryUsage) * prices["memory"]

            resourceUsages = append(resourceUsages, PodResourceUsage{
                Namespace:   namespace,
                PodName:     podName,
                CPUUsage:    fmt.Sprintf("%dm", cpuUsage),
                MemoryUsage: fmt.Sprintf("%dMi", memoryUsage),
                CPUCost:     cpuCost,
                MemoryCost:  memoryCost,
            })

            log.Infof("Pod: %s, Namespace: %s, CPU Usage: %dm, Memory Usage: %dMi, CPU Cost: $%.2f, Memory Cost: $%.2f",
                podName, namespace, cpuUsage, memoryUsage, cpuCost, memoryCost)
        }
    }

    fmt.Println("Resource Usage Analysis:")
    for _, usage := range resourceUsages {
        fmt.Printf("Namespace: %s, Pod: %s, CPU Usage: %s, Memory Usage: %s, CPU Cost: $%.2f, Memory Cost: $%.2f\n",
            usage.Namespace, usage.PodName, usage.CPUUsage, usage.MemoryUsage, usage.CPUCost, usage.MemoryCost)
    }
}