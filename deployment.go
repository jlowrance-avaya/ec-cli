package main

import "fmt"

func getDeployments(database string, outputFormat string) {
	var partitionKey = "/" + database

	fmt.Println("Executing getDeployment with the following parameters:")
	fmt.Printf("---\ndatabase: %s\npartitionKey: %s\noutputFormat: %s\n", database, partitionKey, outputFormat)
}
