package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func getManifests(database string, outputFormat string) {
	var partitionKey = "/" + database

	fmt.Println("Executing getManifest with the following parameters:")
	fmt.Printf("---\ndatabase: %s\npartitionKey: %s\noutputFormat: %s\n", database, partitionKey, outputFormat)

	// Remember to replace with your own values
	// response, err := queryCosmosDB("ec-provisioner", "manifests")
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }

	// fmt.Println("Response: ", response)

	// working function
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
