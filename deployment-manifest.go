package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type Manifest struct {
	ID           string `json:"id"`
	Kind         string `json:"kind"`
	Version      string `json:"version"`
	CustomerName string `json:"customer_name"`
}

func getDeploymentManifests(database string, outputFormat string) {

	fmt.Println("Executing getManifest with the following parameters:")
	fmt.Printf("---\ndatabase: %s\noutputFormat: %s\n---\n", database, outputFormat)

	// new function
	// queryResponse, err := queryManifest("ec-provisioner", "manifests", "52314b92-cecd-4b11-aef8-f0cda6d3bb98")
	// if err != nil {
	// 	log.Printf("readItem failed: %s\n", err)
	// }

	// fmt.Println(queryResponse)
}

func readDeploymentManifest(databaseName string, containerName string, partitionKey string, itemId string) error {

	client := createClient(
		os.Getenv("AZURE_COSMOS_ENDPOINT"),
		os.Getenv("AZURE_COSMOS_KEY"),
	)

	// Create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return fmt.Errorf("failed to create a container client: %s", err)
	}

	// Specifies the value of the partiton key
	pk := azcosmos.NewPartitionKeyString(partitionKey)

	// Read an item
	ctx := context.TODO()
	itemResponse, err := containerClient.ReadItem(ctx, pk, itemId, nil)
	if err != nil {
		return err
	}

	itemResponseBody := struct {
		ID           string `json:"id"`
		Kind         string `json:"kind"`
		Version      string `json:"version"`
		CustomerName string `json:"customer_name"`
	}{}

	err = json.Unmarshal(itemResponse.Value, &itemResponseBody)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(itemResponseBody, "", "    ")
	if err != nil {
		return err
	}
	// fmt.Printf("Read item with Metadata UUID %s\n", itemResponseBody.Kind)
	fmt.Printf("%s\n", b)

	// log.Printf("Status %d. Item %v read. ActivityId %s. Consuming %v Request Units.\n", itemResponse.RawResponse.StatusCode, pk, itemResponse.ActivityID, itemResponse.RequestCharge)

	return nil
}

func queryDeploymentManifest(databaseName string, containerName string, partitionKey string) ([]Manifest, error) {
	client := createClient(
		os.Getenv("AZURE_COSMOS_ENDPOINT"),
		os.Getenv("AZURE_COSMOS_KEY"),
	)

	// Create container client
	containerClient, err := client.NewContainer(databaseName, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create a container client: %s", err)
	}

	// Create a new context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	// Make sure to cancel the context when you're done to avoid a memory leak
	defer cancel()

	// Specifies the value of the partition key
	pk := azcosmos.NewPartitionKeyString(partitionKey)
	queryPager := containerClient.NewQueryItemsPager("SELECT * FROM c", pk, nil)

	var manifests []Manifest
	for queryPager.More() {
		queryResponse, err := queryPager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range queryResponse.Items {
			var manifest Manifest
			if err := json.Unmarshal(item, &manifest); err != nil {
				return nil, err
			}
			manifests = append(manifests, manifest)
		}
	}

	return manifests, nil
}
