package analysis

import (
    "context"
    "fmt"
    "os"

    "github.com/sirupsen/logrus"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)

// Logger setup
var log = logrus.New()

func init() {
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

// PodResourceUsage holds the resource usage details of a pod
type PodResourceUsage struct {
    Namespace  string
    PodName    string
    CPUUsage   string
    MemoryUsage string
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

        usage, err := clientset.MetricsV1beta1().PodMetricses(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
        if err != nil {
            log.WithError(err).Errorf("Failed to get metrics for pod: %s", podName)
            continue
        }

        for _, container := range usage.Containers {
            cpuUsage := container.Usage.Cpu().String()
            memoryUsage := container.Usage.Memory().String()

            resourceUsages = append(resourceUsages, PodResourceUsage{
                Namespace:  namespace,
                PodName:    podName,
                CPUUsage:   cpuUsage,
                MemoryUsage: memoryUsage,
            })

            log.Infof("Pod: %s, Namespace: %s, CPU Usage: %s, Memory Usage: %s", podName, namespace, cpuUsage, memoryUsage)
        }
    }

    fmt.Println("Resource Usage Analysis:")
    for _, usage := range resourceUsages {
        fmt.Printf("Namespace: %s, Pod: %s, CPU Usage: %s, Memory Usage: %s\n", usage.Namespace, usage.PodName, usage.CPUUsage, usage.MemoryUsage)
    }
}