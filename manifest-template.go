package main

import "fmt"

func getManifestTemplates(database string, outputFormat string) {
	var partitionKey = "/" + database

	fmt.Println("Executing getManifestTemplate with the following parameters:")
	fmt.Printf("---\ndatabase: %s\npartitionKey: %s\noutputFormat: %s\n", database, partitionKey, outputFormat)
}
