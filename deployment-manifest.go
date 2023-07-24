package main

import (
	"fmt"
	"net/http"
	"strings"
)

type DeploymentManifest struct {
}

func getDeploymentManifest(api *API, id string) {
	http.HandleFunc("/api/deployment_manifest/", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Here you would validate the token to make sure it's valid. For simplicity,
		// we'll just check that it's a specific value.
		if token != api.Token {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		id := r.URL.Path[len("/api/deployment_manifest/"):]

		// Here you would typically fetch the deployment manifest from your data store
		// based on the id. For simplicity, we'll just return a static string.

		fmt.Fprintf(w, "You requested deployment manifest: %s", id)

	})

	http.ListenAndServe(fmt.Sprintf("%s:8080", api.Endpoint), nil)
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
