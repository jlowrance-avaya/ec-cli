package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

type Credentials struct {
	Username                  string `yaml:"username"`
	Password                  string `yaml:"password"`
	ProvisionerAPIEndpoint    string `yaml:"apiEndpoint"`
	ProvisionerAPIBearerToken string `yaml:"apiBearerToken"`
}

type BearerTokenResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
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

func createCredsFile(username string, password string, apiEndpoint string) error {
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

	_, err = fmt.Fprintf(file, "username: \"%s\"\npassword: \"%s\"\napiEndpoint: \"%s\"\napiBearerToken: \"\"\n", username, password, apiEndpoint)
	return err
}

func getBearerToken() (string, error) {
	url := fmt.Sprintf("https://%s/bearer_token/", *apiEndpointFlag)
	payload := map[string]string{
		"username": "testuser1",
		"password": "testpassword",
	}

	// Convert payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %s", err)
	}

	// Disable SSL certificate verification (equivalent to --insecure in cURL)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create a new HTTP client with the custom transport
	client := &http.Client{Transport: tr}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %s", err)
	}

	// Set request headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %s", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %s", err)
	}

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result BearerTokenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling JSON response: %s", err)
	}

	// Update the creds file with the new access token
	err = updateCredsFile(result.AccessToken)
	if err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func updateCredsFile(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home directory: %s", err)
	}

	credsFile := homeDir + "/.ec-cli/creds"

	// Read the creds file
	content, err := ioutil.ReadFile(credsFile)
	if err != nil {
		return fmt.Errorf("error reading creds file: %s", err)
	}

	// Update the apiBearerToken value
	newContent := strings.ReplaceAll(string(content), "apiBearerToken: \"\"", "apiBearerToken: \""+token+"\"")

	// Write the updated content back to the file
	err = ioutil.WriteFile(credsFile, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to creds file: %s", err)
	}

	return nil
}

func getCreds() (apiEndpoint string, apiBearerToken string) {
	// Replace with the actual path to your file
	filePath := os.Getenv("HOME") + "/.ec-cli/creds"

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed reading data from file: %s", err)
	}

	// Parse the YAML data
	var creds Credentials
	err = yaml.Unmarshal(data, &creds)
	if err != nil {
		log.Fatalf("error occurred during unmarshalling. %s", err)
	}

	apiEndpoint, apiBearerToken = creds.ProvisionerAPIEndpoint, creds.ProvisionerAPIBearerToken
	return
}
