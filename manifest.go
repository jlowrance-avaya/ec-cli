package main

import (
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
	partitionKey = app.Flag("partitionKey", "Name of the container").Required().String()

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
		test(*database, *container, *partitionKey)

	case getManifestsCmd.FullCommand():
		fmt.Printf("Executing 'get manifests' command with page number '%d' and page size '%d'. Database: %s, Container: %s\n",
			*pageNumber, *pageSize, *database, *container)
		// Add your logic here
	}
}

func test(DatabaseName string, ContainerName string, PartitionKey string) {

	endpoint := os.Getenv("AZURE_COSMOS_ENDPOINT")
	if endpoint == "" {
		log.Fatal("AZURE_COSMOS_ENDPOINT could not be found")
	}

	key := os.Getenv("AZURE_COSMOS_KEY")
	if key == "" {
		log.Fatal("AZURE_COSMOS_KEY could not be found")
	}

	var databaseName = *database
	var containerName = *container
	var partitionKey = *partitionKey

	item := struct {
		ID         string `json:"id"`
		CustomerId string `json:"customerId"`
		// Title        string
		// FirstName    string
		// LastName     string
		// EmailAddress string
		// PhoneNumber  string
		// CreationDate string
	}{
		ID:         "1",
		CustomerId: "1",
		// Title:        "Mr",
		// FirstName:    "Luke",
		// LastName:     "Hayes",
		// EmailAddress: "luke12@adventure-works.com",
		// PhoneNumber:  "879-555-0197",
	}

	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		log.Fatal("Failed to create a credential: ", err)
	}

	// Create a CosmosDB client
	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		log.Fatal("Failed to create Azure Cosmos DB db client: ", err)
	}

	err = createDatabase(client, databaseName)
	if err != nil {
		log.Printf("createDatabase failed: %s\n", err)
	}

	err = createContainer(client, databaseName, containerName, partitionKey)
	if err != nil {
		log.Printf("createContainer failed: %s\n", err)
	}

	err = createItem(client, databaseName, containerName, item.CustomerId, item)
	if err != nil {
		log.Printf("createItem failed: %s\n", err)
	}

	err = readItem(client, databaseName, containerName, item.CustomerId, item.ID)
	if err != nil {
		log.Printf("readItem failed: %s\n", err)
	}

	err = deleteItem(client, databaseName, containerName, item.CustomerId, item.ID)
	if err != nil {
		log.Printf("deleteItem failed: %s\n", err)
	}
}

// func manifestGet() {
// 	fmt.Println("Executing 'manifest get' command")
// 	// Add your logic here
// }

// func manifestList() {
// 	fmt.Println("Executing 'manifest list' command")
// 	// Add your logic here
// }

// ---

// func createItem(client *azcosmos.Client, database, container, partitionKey string, item any) error {
// 	//	database = "adventureworks"
// 	//	container = "customer"
// 	//	partitionKey = "1"

// 	item = struct {
// 		ID           string `json:"id"`
// 		CustomerId   string `json:"customerId"`
// 		Title        string
// 		FirstName    string
// 		LastName     string
// 		EmailAddress string
// 		PhoneNumber  string
// 		CreationDate string
// 	}{
// 		ID:           "1",
// 		CustomerId:   "1",
// 		Title:        "Mr",
// 		FirstName:    "Luke",
// 		LastName:     "Hayes",
// 		EmailAddress: "luke12@adventure-works.com",
// 		PhoneNumber:  "879-555-0197",
// 		CreationDate: "2014-02-25T00:00:00",
// 	}
// 	// create container client
// 	containerClient, err := client.NewContainer(database, container)
// 	if err != nil {
// 		return fmt.Errorf("failed to create a container client: %s", err)
// 	}

// 	// specifies the value of the partiton key
// 	pk := azcosmos.NewPartitionKeyString(partitionKey)

// 	b, err := json.Marshal(item)
// 	if err != nil {
// 		return err
// 	}
// 	// setting the item options upon creating ie. consistency level
// 	itemOptions := azcosmos.ItemOptions{
// 		ConsistencyLevel: azcosmos.ConsistencyLevelSession.ToPtr(),
// 	}

// 	// this is a helper function that swallows 409 errors
// 	errorIs409 := func(err error) bool {
// 		var responseErr *azcore.ResponseError
// 		return err != nil && errors.As(err, &responseErr) && responseErr.StatusCode == 409
// 	}

// 	ctx := context.TODO()
// 	itemResponse, err := containerClient.CreateItem(ctx, pk, b, &itemOptions)

// 	switch {
// 	case errorIs409(err):
// 		log.Printf("Item with partitionkey value %s already exists\n", pk)
// 	case err != nil:
// 		return err
// 	default:
// 		log.Printf("Status %d. Item %v created. ActivityId %s. Consuming %v Request Units.\n", itemResponse.RawResponse.StatusCode, pk, itemResponse.ActivityID, itemResponse.RequestCharge)
// 	}

// 	return nil
// }
