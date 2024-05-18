package recommendations

import (
    "context"
    "fmt"
    "github.com/mmarconi93/k0pt/pkg/cloud"
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

// OptimizeResources optimizes resource usage in the Kubernetes cluster
func OptimizeResources(clientset *kubernetes.Clientset) {
    log.Info("Optimizing resources...")

    // Fetch Azure pricing
    prices, err := cloud.FetchAzurePricing()
    if err != nil {
        log.WithError(err).Error("Failed to fetch Azure pricing")
        return
    }

    namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.WithError(err).Error("Failed to list namespaces")
        return
    }

    for _, ns := range namespaces.Items {
        namespaceName := ns.Name
        pods, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            log.WithError(err).Errorf("Failed to list pods in namespace: %s", namespaceName)
            continue
        }

        for _, pod := range pods.Items {
            podName := pod.Name
            usage, err := clientset.MetricsV1beta1().PodMetricses(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
            if err != nil {
                log.WithError(err).Errorf("Failed to get metrics for pod: %s", podName)
                continue
            }

            for _, container := range usage.Containers {
                cpuUsage := container.Usage.Cpu().MilliValue()
                memoryUsage := container.Usage.Memory().Value() / (1024 * 1024) // in MiB

                // Logic to optimize resource requests and limits
                if cpuUsage < 100 && memoryUsage < 100 {
                    // Scale down pod
                    newCPURequest := "50m"
                    newMemoryRequest := "50Mi"
                    optimizePod(clientset, namespaceName, podName, newCPURequest, newMemoryRequest, prices)
                } else if cpuUsage > 500 || memoryUsage > 512 {
                    // Scale up pod
                    newCPURequest := "1"
                    newMemoryRequest := "1Gi"
                    optimizePod(clientset, namespaceName, podName, newCPURequest, newMemoryRequest, prices)
                }
            }
        }
    }
}

// optimizePod updates the resource requests and limits for a pod
func optimizePod(clientset *kubernetes.Clientset, namespace, podName, cpuRequest, memoryRequest string, prices map[string]float64) {
    pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
    if err != nil {
        log.WithError(err).Errorf("Failed to get pod: %s", podName)
        return
    }

    for i := range pod.Spec.Containers {
        pod.Spec.Containers[i].Resources.Requests["cpu"] = cpuRequest
        pod.Spec.Containers[i].Resources.Requests["memory"] = memoryRequest
    }

    _, err = clientset.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
    if err != nil {
        log.WithError(err).Errorf("Failed to update pod: %s", podName)
        return
    }

    cpuCost := prices["cpu"]
    memoryCost := prices["memory"]

    log.Infof("Optimized pod: %s in namespace: %s with CPU request: %s and Memory request: %s. Estimated cost: $%.2f",
        podName, namespace, cpuRequest, memoryRequest, cpuCost+memoryCost)
    fmt.Printf("✅ Optimized pod: %s in namespace: %s\n", podName, namespace)
}