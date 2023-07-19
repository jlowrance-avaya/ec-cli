package main

import "fmt"

type DeploymentManifest struct {
}

func getDeploymentManifest(api *API, id string) {
	fmt.Println("getDeploymentManifest...")
}

func getDeploymentManifests(api *API) {
	fmt.Println("getDeploymentManifests...")
}

func editDeploymentManifest(api *API, id string) {
	fmt.Println("editDeploymentManifest...")
}

func deleteDeploymentManifest(api *API, id string) {
	fmt.Println("deleteDeploymentManifest...")
}

func createDeploymentManifest(api *API, id string) {
	fmt.Println("createDeploymentManifest...")
}
