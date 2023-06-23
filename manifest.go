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
	CreatedAt          int64             `json:"createdAt"`
	UpdatedAt          *time.Time        `json:"updatedAt"`
	ID                 string            `json:"id"`
	Version            string            `json:"version"`
	Kind               string            `json:"kind"`
	Build              string            `json:"build"`
	DeploymentID       string            `json:"deployment_id"`
	ProductSKU         string            `json:"product_sku"`
	OrderID            string            `json:"order_id"`
	CustomerID         string            `json:"customer_id"`
	CustomerName       string            `json:"customer_name"`
	CustomerCode       interface{}       `json:"customer_code"`
	ApplicationAccount string            `json:"application_account"`
	Metadata           ManifestMetadata  `json:"metadata"`
	Variables          ManifestVariables `json:"variables"`
	Jobs               []ManifestJob     `json:"jobs"`
	RID                string            `json:"_rid"`
	Self               string            `json:"_self"`
	ETag               string            `json:"_etag"`
	Attachments        string            `json:"_attachments"`
	Timestamp          int64             `json:"_ts"`
}

type ManifestMetadata struct {
	UUID        string        `json:"uuid"`
	Name        string        `json:"name"`
	Description interface{}   `json:"description"`
	Labels      interface{}   `json:"labels"`
	Tags        []ManifestTag `json:"tags"`
	Selector    interface{}   `json:"selector"`
	Created     string        `json:"created"`
	Creator     interface{}   `json:"creator"`
	Modified    interface{}   `json:"modified"`
	Modifier    interface{}   `json:"modifier"`
}

type ManifestTag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ManifestVariables map[string]string

type ManifestJob struct {
	Kind     string              `json:"kind"`
	Version  string              `json:"version"`
	Metadata ManifestJobMetadata `json:"metadata"`
	Type     string              `json:"type"`
	Subtype  string              `json:"subtype"`
	Payload  ManifestPayload     `json:"payload"`
	Sequence int                 `json:"sequence"`
	Target   string              `json:"target"`
}

type ManifestJobMetadata struct {
	UUID        string      `json:"uuid"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Labels      interface{} `json:"labels"`
	Tags        interface{} `json:"tags"`
	Selector    interface{} `json:"selector"`
	Created     string      `json:"created"`
	Creator     interface{} `json:"creator"`
	Modified    interface{} `json:"modified"`
	Modifier    interface{} `json:"modifier"`
}

type ManifestPayload struct {
	Workspace  string       `json:"workspace"`
	ConfigRepo ManifestRepo `json:"config_repo"`
	Repo       ManifestRepo `json:"repo"`
}

type ManifestRepo struct {
	URL     string `json:"url"`
	Version string `json:"version"`
}

var (
	getCmd = app.Command("get", "Get command")
	// Define global flags
	database     = app.Flag("database", "Name of the database").Required().String()
	container    = app.Flag("container", "Name of the container").Required().String()
	partitionKey = app.Flag("partitonKey", "Name of the container").Required().String()

	getManifestCmd = getCmd.Command("manifest", "Get a manifest")
	// Add flags for 'get manifest'
	name         = getManifestCmd.Flag("name", "Name of the manifest").Required().String()
	outputFormat = getManifestCmd.Flag("output", "Output format").Default("yaml").Enum("json", "yaml")

	getManifestsCmd = getCmd.Command("manifests", "Get all manifests")
	// Add flags for 'get manifests'
	pageNumber = getManifestsCmd.Flag("page", "Page number for manifest listing").Int()
	pageSize   = getManifestsCmd.Flag("size", "Size of each page for manifest listing").Int()
)

func handleGetCommand(cmd string) {
	switch cmd {
	case getManifestCmd.FullCommand():
		fmt.Printf("Executing 'get manifest' command with name '%s' format:'%s' format. Database: %s, Container: %s\n",
			*name, *outputFormat, *database, *container)
		// Add your logic here
		client := createClient(
			os.Getenv("AZURE_COSMOS_ENDPOINT"),
			os.Getenv("AZURE_COSMOS_KEY"),
		)

		// createDatabaseContainerItem

		// item := struct {
		// 	ID         string `json:"id"`
		// 	CustomerId string `json:"customerId"`
		// }{
		// 	ID:         "1",
		// 	CustomerId: "1",
		// }

		item := struct {
			ID           string `json:"id"`
			MetaDataUUID string `json:"metadata.uuid"`
		}{
			ID:           "7fd55dcf-ef05-4c99-9b6d-040fd666f018",
			MetaDataUUID: "52314b92-cecd-4b11-aef8-f0cda6d3bb98",
		}

		err := readManifest(client, *database, *container, item.MetaDataUUID, item.ID)
		if err != nil {
			log.Printf("readItem failed: %s\n", err)
		}

	case getManifestsCmd.FullCommand():
		fmt.Printf("Executing 'get manifests' command with page number '%d' and page size '%d'. Database: %s, Container: %s\n",
			*pageNumber, *pageSize, *database, *container)
	}
}

func readManifest(client *azcosmos.Client, databaseName string, containerName string, partitionKey string, itemId string) error {
	//	databaseName = "adventureworks"
	//	containerName = "customer"
	//	partitionKey = "1"
	//	itemId = "1"

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
