package main

import "fmt"

type Deployment struct {
}

func getDeployment(api *API, id string) {
	fmt.Println("getDeployment...")
}

func getDeployments(api *API) {
	fmt.Println("getDeployments...")
}

func editDeployment(api *API, id string) {
	fmt.Println("editDeployment...")
}

func deleteDeployment(api *API, id string) {
	fmt.Println("deleteDeployment...")
}

func createDeployment(api *API, id string) {
	fmt.Println("createDeployment...")
}
