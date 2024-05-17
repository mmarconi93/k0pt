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

// GetAzureCredentials retrieves Azure credentials
func getAzureCredentials() (*azidentity.DefaultAzureCredential, error) {
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        log.WithError(err).Error("Failed to obtain Azure credentials")
        return nil, err
    }
    return cred, nil
}

// GetEnvironmentVariable retrieves the value of an environment variable
func getEnvironmentVariable(name string) (string, error) {
    value := os.Getenv(name)
    if value == "" {
        err := fmt.Errorf("%s environment variable is not set", name)
        log.WithError(err).Errorf("%s is required", name)
        return "", err
    }
    return value, nil
}

// GetClusterState retrieves the state of the Azure Kubernetes Service (AKS) cluster
func GetClusterState(clusterName string) {
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

    res, err := client.Get(ctx, resourceGroupName, clusterName, nil)
    if err != nil {
        callbacks.PrintErrorMessage(err)
        log.WithError(err).Error("Failed to get cluster state")
        return
    }

    if res.ManagedCluster.Properties.PowerState != nil {
        fmt.Printf("AKS cluster '%v' is %s\n", clusterName, *res.ManagedCluster.Properties.PowerState.Code)
    } else {
        fmt.Println("Cluster power state is not available")
    }
}
