package main

import (
	"fmt"
	"time"
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

	getManifestCmd = getCmd.Command("manifest", "Get a manifest")
	// Add flags for 'get manifest'
	manifestName = getManifestCmd.Flag("name", "Name of the manifest").Required().String()
	outputFormat = getManifestCmd.Flag("output", "Output format").Default("yaml").Enum("json", "yaml")

	getManifestsCmd = getCmd.Command("manifests", "Get all manifests")
	// Add flags for 'get manifests'
	pageNumber = getManifestsCmd.Flag("page", "Page number for manifest listing").Int()
	pageSize   = getManifestsCmd.Flag("size", "Size of each page for manifest listing").Int()
)

func handleGetCommand(cmd string) {
	switch cmd {
	case getManifestCmd.FullCommand():
		fmt.Printf("Executing 'get manifest' command with name '%s' in '%s' format\n", *manifestName, *outputFormat)
		// Add your logic here
	case getManifestsCmd.FullCommand():
		fmt.Printf("Executing 'get manifests' command with page number '%d' and page size '%d'\n", *pageNumber, *pageSize)
		// Add your logic here
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
