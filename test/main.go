package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func main() {

	checkEnvVars([]string{"AZURE_COSMOS_ENDPOINT", "AZURE_COSMOS_KEY"})

	client := createClient(
		os.Getenv("AZURE_COSMOS_ENDPOINT"),
		os.Getenv("AZURE_COSMOS_KEY"),
	)

	databaseName := "ec-provisioner"
	containerName := "manifests"

	// databaseClient, err := client.NewDatabase(databaseName)
	// if err != nil {
	// 	panic(err)
	// }

	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		log.Fatal(err)
	}

	containerResponse, err := containerClient.Read(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", containerResponse)
}

func checkEnvVars(varNames []string) {
	missingVars := []string{}
	for _, varName := range varNames {
		val := os.Getenv(varName)
		if val == "" {
			missingVars = append(missingVars, varName)
		}
	}
	if len(missingVars) > 0 {
		fmt.Printf("Please set the following environment variables:\n%s\n", strings.Join(missingVars, ", "))
		os.Exit(1)
	}
}

func createClient(endpoint string, key string) *azcosmos.Client {
	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		log.Fatalf("Failed to create a credential: %v", err)
	}

	// Create a CosmosDB client
	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		log.Fatalf("Failed to create Azure Cosmos DB client: %v", err)
	}

	return client
}
