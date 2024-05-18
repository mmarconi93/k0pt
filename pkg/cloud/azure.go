package cloud

import (
    "context"
    "fmt"
    "os"

    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/consumption/armconsumption"
    "github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
    log.Out = os.Stdout
    log.SetFormatter(&logrus.JSONFormatter{})
}

// FetchAzurePricing fetches the current pricing details from Azure
func FetchAzurePricing() (map[string]float64, error) {
    cred, err := azidentity.NewDefaultAzureCredential(nil)
    if err != nil {
        return nil, fmt.Errorf("failed to obtain a credential: %v", err)
    }

    client, err := armconsumption.NewPricesClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), cred, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %v", err)
    }

    ctx := context.Background()
    prices, err := client.ListByBillingPeriod(ctx, "latest", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get prices: %v", err)
    }

    priceMap := make(map[string]float64)
    for _, item := range prices.Value {
        priceMap[*item.MeterID] = *item.UnitPrice
    }

    return priceMap, nil
}