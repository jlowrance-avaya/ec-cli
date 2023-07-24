package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	app = kingpin.New("ec", "Enterprise Cloud CLI")

	loginCommand    = app.Command("login", "Login to the application.")
	usernameFlag    = loginCommand.Flag("username", "Username for login").String()
	passwordFlag    = loginCommand.Flag("password", "Password for login (usage of this flag is not recommended)").String()
	apiEndpointFlag = loginCommand.Flag("endpoint", "API Endpoint to use").Default("provisioner-api.shsrv-nonprod.private.ec.avayacloud.com").String()

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

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case loginCommand.FullCommand():
		username := *usernameFlag
		if username == "" {
			var err error
			username, err = readUsername("Enter username: ")
			if err != nil {
				fmt.Println("Error reading password:", err)
				return
			}
		}

		password := *passwordFlag
		if password == "" {
			var err error
			password, err = readPassword("Enter password: ")
			if err != nil {
				fmt.Println("Error reading password:", err)
				return
			}
		}

		err := createCredsFile(username, password, *apiEndpointFlag)
		if err != nil {
			fmt.Println("Error creating creds file:", err)
		}

		// token, _ := getAccessToken()
		// fmt.Println(token)

	case get.FullCommand():
		switch *getResource {
		// case "deploymentManifest":
		// 	getDeploymentManifest(apiCredentials, "asdf")
		case "deploymentManifests":
			bearerToken, err := getBearerToken()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			fmt.Println("Bearer Token:", bearerToken)

		case "deploymentManifestTemplate":
			getDeploymentManifestTemplate()
		case "deploymentManifestTemplates":
			getDeploymentManifestTemplates()
		// you can add more cases here as needed
		default:
			fmt.Println("Invalid resource")
		}

	case edit.FullCommand():
		switch *editResource {
		case "deploymentManifest":
			// editDeploymentManifest(apiCredentials, "asdf")
		case "deploymentManifestTemplate":
			editDeploymentManifestTemplate()
		default:
			fmt.Println("Invalid resource")
		}

	case delete.FullCommand():
		switch *deleteResource {
		case "deploymentManifest":
			// deleteDeploymentManifest(apiCredentials, "asdf")
		case "deploymentManifestTemplate":
			deleteDeploymentManifestTemplate()
		default:
			fmt.Println("Invalid resource")
		}

	case create.FullCommand():
		switch *createResource {
		case "deploymentManifest":
			// createDeploymentManifest(apiCredentials, "asdf")
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

// func checkEnvVars(varNames []string) {
// 	missingVars := []string{}
// 	for _, varName := range varNames {
// 		val := os.Getenv(varName)
// 		if val == "" {
// 			missingVars = append(missingVars, varName)
// 		}
// 	}
// 	if len(missingVars) > 0 {
// 		fmt.Printf("Please set the following environment variables:\n%s\n", strings.Join(missingVars, ", "))
// 		os.Exit(1)
// 	}
// }

func handle(err error) {
	if err != nil {
		// This will print the error and stop the program.
		log.Fatal(err)
	}
}
