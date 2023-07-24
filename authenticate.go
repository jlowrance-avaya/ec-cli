package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Credentials struct {
	Username                  string `json:"username"`
	Password                  string `json:"password"`
	ProvisionerAPIEndpoint    string `json:"provisioner_api_endpoint"`
	ProvisionerAPIAccessToken string `json:"provisioner_api_access_token"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func readUsername(prompt string) (string, error) {
	fmt.Print(prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePassword)), nil
}

func createCredsFile(username string, password string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	dir := filepath.Join(usr.HomeDir, ".ec-cli")
	err = os.MkdirAll(dir, 0755) // Create directory if it doesn't exist
	if err != nil {
		return err
	}

	credsFile := filepath.Join(dir, "creds")
	file, err := os.OpenFile(credsFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("\ncredentials stored in ~/.ec-cli/creds")

	_, err = fmt.Fprintf(file, "username: \"%s\"\npassword: \"%s\"\n", username, password)
	return err
}

func getAccessToken(apiEndpoint string) (TokenResponse, error) {
	var tokenResp TokenResponse

	usr, err := user.Current()
	if err != nil {
		return tokenResp, err
	}

	credsFile := filepath.Join(usr.HomeDir, ".ec-cli", "creds")

	data, err := ioutil.ReadFile(credsFile)
	if err != nil {
		return tokenResp, err
	}

	creds := Credentials{}
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return tokenResp, err
	}

	jsonData, err := json.Marshal(creds)
	if err != nil {
		return tokenResp, err
	}

	apiEndpointAuthzUri := fmt.Sprintf("//%s/bearer_token/", apiEndpoint)

	resp, err := http.Post(apiEndpointAuthzUri, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return tokenResp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return tokenResp, err
	}

	return tokenResp, nil
}

func authenticate(ProvisionerApiEndpoint string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	credsFile := filepath.Join(usr.HomeDir, ".ec-cli", "creds")

	data, err := ioutil.ReadFile(credsFile)
	if err != nil {
		return err
	}

	creds := Credentials{}
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return err
	}

	tokenResp, err := getAccessToken(ProvisionerApiEndpoint)
	if err != nil {
		return err
	}

	creds.ProvisionerAPIEndpoint = "provisioner-api.shsrv-nonprod.private.ec.avayacloud.com"
	creds.ProvisionerAPIAccessToken = tokenResp.AccessToken

	fmt.Printf(creds.ProvisionerAPIEndpoint)
	fmt.Printf(creds.ProvisionerAPIAccessToken)

	jsonCreds, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(credsFile, jsonCreds, 0600)
	if err != nil {
		return err
	}

	return nil
}

func oldAuthenticate() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	credsFile := filepath.Join(usr.HomeDir, ".ec-cli", "creds")

	data, err := ioutil.ReadFile(credsFile)
	if err != nil {
		return err
	}

	creds := Credentials{}
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	apiEndpoint := GetEnvWithDefault("PROVISIONER_API_ENDPOINT", "provisioner-api.shsrv-nonprod.private.ec.avayacloud.com")
	apiEndpointAuthzUri := fmt.Sprintf("//%s/bearer_token/", apiEndpoint)

	resp, err := http.Post(apiEndpointAuthzUri,
		"application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	tokenResp := TokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return err
	}

	creds.ProvisionerAPIEndpoint = apiEndpoint
	creds.ProvisionerAPIAccessToken = tokenResp.AccessToken

	err = os.Setenv("PROVISIONER_API_ENDPOINT", apiEndpoint)
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
	}

	err = os.Setenv("PROVISIONER_API_ACCESS_TOKEN", tokenResp.AccessToken)
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
	}

	if err != nil {
	} else {
		fmt.Println("Set PROVISIONER_API_ENDPOINT to", apiEndpoint)
	}

	updatedCreds, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(credsFile, updatedCreds, 0600)
	if err != nil {
		return err
	}

	return nil
}
