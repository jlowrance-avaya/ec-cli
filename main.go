package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/kingpin"
)

var (
	app    = kingpin.New("ec", "Enterprise Cloud CLI")
	get    = app.Command("get", "Get a resource")
	edit   = app.Command("edit", "Edit a resource")
	delete = app.Command("delete", "Delete a resource")
	create = app.Command("create", "Create a resource")

	objects = []string{
		"deploymentManifest",
		"deploymentManifests",
		"deploymentManifestTemplate",
		"deploymentManifestTemplates",
	}

	getResource    = get.Arg("resource", "Resource to operate on").Enum(objects...)
	editResource   = edit.Arg("resource", "Resource to operate on").Enum(objects...)
	deleteResource = delete.Arg("resource", "Resource to operate on").Enum(objects...)
	createResource = create.Arg("resource", "Resource to operate on").Enum(objects...)
)

type API struct {
	Endpoint string
	Token    string
}

func main() {

	checkEnvVars([]string{
		"PROVISIONER_API_ENDPOINT",
		"PROVISIONER_API_TOKEN",
	})

	api := &API{
		Endpoint: os.Getenv("PROVISIONER_API_ENDPOINT"),
		Token:    os.Getenv("PROVISIONER_API_TOKEN"),
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case get.FullCommand():
		switch *getResource {
		case "deploymentManifest":
			getDeploymentManifest(api, "asdf")
		case "deploymentManifestTemplate":
			getDeploymentManifestTemplate()

		case "deploymentManifests":
			getDeploymentManifests(api)
		case "deploymentManifestTemplates":
			getDeploymentManifestTemplates()
		// you can add more cases here as needed
		default:
			fmt.Println("Invalid resource")
		}

	case edit.FullCommand():
		switch *editResource {
		case "deploymentManifest":
			editDeploymentManifest(api, "asdf")
		case "deploymentManifestTemplate":
			editDeploymentManifestTemplate()
		default:
			fmt.Println("Invalid resource")
		}

	case delete.FullCommand():
		switch *deleteResource {
		case "deploymentManifest":
			deleteDeploymentManifest(api, "asdf")
		case "deploymentManifestTemplate":
			deleteDeploymentManifestTemplate()
		default:
			fmt.Println("Invalid resource")
		}

	case create.FullCommand():
		switch *createResource {
		case "deploymentManifest":
			createDeploymentManifest(api, "asdf")
		case "deploymentManifestTemplate":
			createDeploymentManifestTemplate()
		default:
			fmt.Println("Invalid resource")
		}
	}

	// flags
	// outputFormat := getCommand.Flag("output", "Output format").Default("yaml").Enum("json", "yaml")

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

func handle(err error) {
	if err != nil {
		// This will print the error and stop the program.
		log.Fatal(err)
	}
}
