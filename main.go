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
	app    = kingpin.New("ec", "Enterprise Cloud CLI")
	get    = app.Command("get", "Get a resource")
	edit   = app.Command("edit", "Edit a resource")
	delete = app.Command("delete", "Delete a resource")
	create = app.Command("create", "Create a resource")

	getResource    = get.Arg("resource", "Resource to operate on").Enum("deploymentManifest", "deploymentManifestTemplate", "deployment")
	editResource   = edit.Arg("resource", "Resource to operate on").Enum("deploymentManifest", "deploymentManifestTemplate", "deployment")
	deleteResource = delete.Arg("resource", "Resource to operate on").Enum("deploymentManifest", "deploymentManifestTemplate", "deployment")
	createResource = create.Arg("resource", "Resource to operate on").Enum("deploymentManifest", "deploymentManifestTemplate", "deployment")
)

type Query struct {
	Query string `json:"query"`
}

type CosmosDBResponse struct {
	Documents []map[string]interface{} `json:"Documents"`
}

func main() {

	checkEnvVars([]string{
		"PROVISIONER_API_ENDPOINT",
		"PROVISIONER_API_TOKEN",
	})

	baseRequestUrl := fmt.Sprintf("https://%s:443/", os.Getenv("PROVISIONER_API_ENDPOINT"))
	authzHeader := fmt.Sprintf("Authorization: Bearer %s", os.Getenv("PROVISIONER_API_TOKEN"))

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case get.FullCommand():
		fmt.Printf("Getting %s\n", *getResource)
	case edit.FullCommand():
		fmt.Printf("Editing %s\n", *editResource)
	case delete.FullCommand():
		fmt.Printf("Deleting %s\n", *deleteResource)
	case create.FullCommand():
		fmt.Printf("Creating %s\n", *createResource)
	}

	// sub-commands
	// "manifest"
	// "manifests"
	// "manifestTemplate"
	// "manifestTemplates"
	// "deployment"
	// "deployments"
	// "job"
	// "jobs"

	// flags
	// outputFormat := getCommand.Flag("output", "Output format").Default("yaml").Enum("json", "yaml")

	// getCommand.Command("subscription", "Get Azure subscription").Action(func(c *kingpin.ParseContext) error {
	// 	getManifests(baseRequestUrl, *outputFormat)
	// 	return nil
	// })

	// getCommand.Command("manifestTemplates", "Get multiple manifestTemplates").Action(func(c *kingpin.ParseContext) error {
	// 	getManifestTemplates(baseRequestUrl, *outputFormat)
	// 	return nil
	// })

	// getCommand.Command("deployments", "Get multiple deployments").Action(func(c *kingpin.ParseContext) error {
	// 	getDeployments(baseRequestUrl, *outputFormat)
	// 	return nil
	// })

	fmt.Println(
		"---\n",
		baseRequestUrl+"\n",
		authzHeader+"\n",
	)

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
		handle(err)
	}

	// Create a CosmosDB client
	client, err := azcosmos.NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		handle(err)
	}

	return client
}

func handle(err error) {
	if err != nil {
		// This will print the error and stop the program.
		log.Fatal(err)
	}
}
