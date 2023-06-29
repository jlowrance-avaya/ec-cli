package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

type QueryParameter struct {
	// Name represents the name of the parameter in the parametrized query.
	Name string `json:"name"`
	// Value represents the value of the parameter in the parametrized query.
	Value any `json:"value"`
}

func getManifests(database string, outputFormat string) {
	var partitionKey = "/" + database

	fmt.Println("Executing getManifest with the following parameters:")
	fmt.Printf("---\ndatabase: %s\npartitionKey: %s\noutputFormat: %s\n", database, partitionKey, outputFormat)

	item := struct {
		ID           string `json:"id"`
		MetaDataUUID string `json:"metadata.uuid"`
	}{
		ID:           "7fd55dcf-ef05-4c99-9b6d-040fd666f018",
		MetaDataUUID: "52314b92-cecd-4b11-aef8-f0cda6d3bb98",
	}

	err := readManifest("ec-provisioner", "manifests", item.MetaDataUUID, item.ID)
	if err != nil {
		log.Printf("readItem failed: %s\n", err)
	}

	// new function
	queryResponse, err := queryManifest("ec-provisioner", "manifests", item.MetaDataUUID)
	if err != nil {
		log.Printf("readItem failed: %s\n", err)
	}

	fmt.Println(queryResponse)
}

func readManifest(databaseName string, containerName string, partitionKey string, itemId string) error {

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

func queryManifest(databaseName string, containerName string, partitionKey string) ([]Manifest, error) {
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
