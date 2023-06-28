package main

import "fmt"

func getManifests(database string, outputFormat string) {
	var partitionKey = "/" + database

	fmt.Println("Executing getManifest with the following parameters:")
	fmt.Printf("---\ndatabase: %s\npartitionKey: %s\noutputFormat: %s\n", database, partitionKey, outputFormat)
}
