package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/alecthomas/kingpin"
)

var (
	app = kingpin.New("ec", "enterprise cloud CLI")
)

type Query struct {
	Query string `json:"query"`
}

type CosmosDBResponse struct {
	Documents []map[string]interface{} `json:"Documents"`
}

func main() {

	checkEnvVars([]string{"AZURE_COSMOS_ENDPOINT", "AZURE_COSMOS_KEY"})

	app := kingpin.New("ec", "A command-line app")

	// Set up "get" command
	getCommand := app.Command("get", "Get operation")
	database := getCommand.Flag("database", "Name of the database").Required().String()
	outputFormat := getCommand.Flag("output", "Output format").Default("yaml").Enum("json", "yaml")

	getCommand.Command("manifests", "Get multiple manifests").Action(func(c *kingpin.ParseContext) error {
		getManifests(*database, *outputFormat)
		return nil
	})

	getCommand.Command("manifestTemplates", "Get multiple manifestTemplates").Action(func(c *kingpin.ParseContext) error {
		getManifestTemplates(*database, *outputFormat)
		return nil
	})

	getCommand.Command("deployments", "Get multiple deployments").Action(func(c *kingpin.ParseContext) error {
		getDeployments(*database, *outputFormat)
		return nil
	})

	kingpin.MustParse(app.Parse(os.Args[1:]))
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
