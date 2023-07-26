package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

type DeploymentManifest struct {
}

func getDeploymentManifest() {
}

func getDeploymentManifests() {
	apiEndpoint, apiBearerToken := getCreds()
	fmt.Printf("API Endpoint: %s\n", apiEndpoint)
	fmt.Printf("API Bearer Token: %s\n", apiBearerToken)

	url := fmt.Sprintf("https://%s/api/deployment_manifest/all?Page%%20Number=1&Size=10", apiEndpoint)

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

func editDeploymentManifest() {
	fmt.Println("editDeploymentManifest...")
}

func deleteDeploymentManifest() {
	fmt.Println("deleteDeploymentManifest...")
}

func createDeploymentManifest() {
	fmt.Println("createDeploymentManifest...")
}
