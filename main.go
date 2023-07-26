package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

		_, err = getBearerToken()
		if err != nil {
			fmt.Println("Unable to get Bearer token:", err)
			return
		}

	case get.FullCommand():
		switch *getResource {
		case "deploymentManifests":
			getObjects("deployment_manifest")
		case "deploymentManifestTemplates":
			getObjects("deployment_template")
		case "deploymentManifest":
			// getObject("deployment_manifest")
		case "deploymentManifestTemplate":
			// getObject("deployment_template")
		default:
			fmt.Println("Invalid resource")
		}

	case edit.FullCommand():
		switch *editResource {
		case "deploymentManifestTemplate":
			// editObject("deployment_template")
		default:
			fmt.Println("Invalid resource")
		}

	case delete.FullCommand():
		switch *deleteResource {
		case "deploymentManifestTemplate":
			// deleteObject("deployment_template")
		default:
			fmt.Println("Invalid resource")
		}

	case create.FullCommand():
		switch *createResource {
		case "deploymentManifestTemplate":
			// createObject("deployment_template")
		default:
			fmt.Println("Invalid resource")
		}
	}

	// flags
	// outputFormat := getCommand.Flag("output", "Output format").Default("yaml").Enum("json", "yaml")

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getObjects(objectType string) {
	apiEndpoint, apiBearerToken := getCreds()

	url := fmt.Sprintf("https://%s/api/%s/all?Page%%20Number=1&Size=10", apiEndpoint, objectType)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiBearerToken))

	// Create a transport to skip SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to make the request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}
