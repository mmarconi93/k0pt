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

// CostSavingRecommendations holds recommendations for cost savings
type CostSavingRecommendations struct {
    Namespace           string
    UnderutilizedPods   []string
    OverprovisionedPods []string
    SuggestedActions    []string
}

// CalculateCostSavings calculates cost savings based on resource usage in the Kubernetes cluster
func CalculateCostSavings(clientset *kubernetes.Clientset) {
    log.Info("Calculating cost savings...")

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

    var recommendations []CostSavingRecommendations

    for _, ns := range namespaces.Items {
        namespaceName := ns.Name
        pods, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            log.WithError(err).Errorf("Failed to list pods in namespace: %s", namespaceName)
            continue
        }

        nsRecommendations := CostSavingRecommendations{Namespace: namespaceName}

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

                cpuCost := float64(cpuUsage) / 1000 * prices["cpu"]
                memoryCost := float64(memoryUsage) * prices["memory"]

                // Logic to identify underutilized and overprovisioned pods
                if cpuUsage < 100 && memoryUsage < 100 {
                    nsRecommendations.UnderutilizedPods = append(nsRecommendations.UnderutilizedPods, podName)
                    nsRecommendations.SuggestedActions = append(nsRecommendations.SuggestedActions, fmt.Sprintf("Consider scaling down pod: %s in namespace: %s. Estimated cost savings: $%.2f", podName, namespaceName, cpuCost+memoryCost))
                } else if cpuUsage > 500 || memoryUsage > 512 {
                    nsRecommendations.OverprovisionedPods = append(nsRecommendations.OverprovisionedPods, podName)
                    nsRecommendations.SuggestedActions = append(nsRecommendations.SuggestedActions, fmt.Sprintf("Consider scaling up pod: %s in namespace: %s. Estimated cost increase: $%.2f", podName, namespaceName, cpuCost+memoryCost))
                }
            }
        }

        if len(nsRecommendations.UnderutilizedPods) > 0 || len(nsRecommendations.OverprovisionedPods) > 0 {
            recommendations = append(recommendations, nsRecommendations)
        }
    }

    for _, rec := range recommendations {
        fmt.Printf("Namespace: %s\n", rec.Namespace)
        fmt.Printf("Underutilized Pods: %v\n", rec.UnderutilizedPods)
        fmt.Printf("Overprovisioned Pods: %v\n", rec.OverprovisionedPods)
        fmt.Printf("Suggested Actions: %v\n", rec.SuggestedActions)
    }
}