package control

import (
    "context"
    "fmt"
    "os"

    "github.com/mmarconi93/k0pt/pkg/callbacks"
    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
    "github.com/sirupsen/logrus"
)

// Logger setup
var log = logrus.New()

func init() {
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

// getAzureCredentials retrieves Azure credentials
func getAzureCredentials() (*azidentity.DefaultAzureCredential, error) {
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        log.WithError(err).Error("Failed to obtain Azure credentials")
        return nil, err
    }
    return cred, nil
}

// getEnvironmentVariable retrieves the value of an environment variable
func getEnvironmentVariable(name string) (string, error) {
    value := os.Getenv(name)
    if value == "" {
        err := fmt.Errorf("%s environment variable is not set", name)
        log.WithError(err).Errorf("%s is required", name)
        return "", err
    }
    return value, nil
}

// stopCluster stops the Azure Kubernetes Service (AKS) cluster
func stopCluster(ctx context.Context, client *armcontainerservice.ManagedClustersClient, resourceGroupName, clusterName string) error {
    poller, err := client.BeginStop(ctx, resourceGroupName, clusterName, nil)
    if err != nil {
        callbacks.PrintErrorMessage(err)
        log.WithError(err).Error("Failed to begin stopping the cluster")
        return err
    }
    _, err = poller.PollUntilDone(ctx, nil)
    if err != nil {
        log.WithError(err).Error("Failed to stop the cluster")
        return err
    }
    log.Infof("Cluster %s stopped successfully", clusterName)
    return nil
}

// startCluster starts the Azure Kubernetes Service (AKS) cluster
func startCluster(ctx context.Context, client *armcontainerservice.ManagedClustersClient, resourceGroupName, clusterName string) error {
    poller, err := client.BeginStart(ctx, resourceGroupName, clusterName, nil)
    if err != nil {
        callbacks.PrintErrorMessage(err)
        log.WithError(err).Error("Failed to begin starting the cluster")
        return err
    }
    _, err = poller.PollUntilDone(ctx, nil)
    if err != nil {
        log.WithError(err).Error("Failed to start the cluster")
        return err
    }
    log.Infof("Cluster %s started successfully", clusterName)
    return nil
}

// ClusterStop stops an AKS cluster
func ClusterStop(clusterName string) {
    if clusterName == "" {
        log.Error("Cluster name argument is required")
        return
    }

    cred, err := getAzureCredentials()
    if err != nil {
        return
    }

    ctx := context.Background()

    subscriptionID, err := getEnvironmentVariable("AZURE_SUBSCRIPTION_ID")
    if err != nil {
        return
    }

    client, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
    if err != nil {
        log.WithError(err).Error("Failed to create AKS client")
        return
    }

    resourceGroupName, err := getEnvironmentVariable("AZURE_RESOURCE_GROUP_NAME")
    if err != nil {
        return
    }

    if err := stopCluster(ctx, client, resourceGroupName, clusterName); err != nil {
        log.WithError(err).Error("Failed to stop AKS cluster")
    }
}

// ClusterStart starts an AKS cluster
func ClusterStart(clusterName string) {
    if clusterName == "" {
        log.Error("Cluster name argument is required")
        return
    }

    cred, err := getAzureCredentials()
    if err != nil {
        return
    }

    ctx := context.Background()

    subscriptionID, err := getEnvironmentVariable("AZURE_SUBSCRIPTION_ID")
    if err != nil {
        return
    }

    client, err := armcontainerservice.NewManagedClustersClient(subscriptionID, cred, nil)
    if err != nil {
        log.WithError(err).Error("Failed to create AKS client")
        return
    }

    resourceGroupName, err := getEnvironmentVariable("AZURE_RESOURCE_GROUP_NAME")
    if err != nil {
        return
    }

    if err := startCluster(ctx, client, resourceGroupName, clusterName); err != nil {
        log.WithError(err).Error("Failed to start AKS cluster")
    }
}
