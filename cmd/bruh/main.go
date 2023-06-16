/*
The bruh command line tool is a simple utility that parses Bicep files and directories
and prints out the resources that are defined in them.
*/
package main

import (
	"fmt"
	"log"

	"github.com/christosgalano/bruh/pkg/parser"
)

func main() {
	printResourceInfo := func(result parser.ResourceInfo) {
		fmt.Printf("%s\n", result.ID)
		fmt.Printf("  Namespace: %s\n", result.Namespace)
		fmt.Printf("  Resource: %s\n", result.Resource)
		fmt.Printf("  Version: %s\n", result.Version)
	}

	file := "/Users/galano/Developer/Christos/Application/Go/go-development/bicep/modules/compute.bicep"
	fileResults, err := parser.ParseFile(file)
	if err != nil {
		log.Fatalf("Parser error: %s", err)
	}
	fmt.Printf("File: %s\n\n", file)
	for _, result := range fileResults {
		printResourceInfo(result)
	}

	fmt.Printf("\n\n")

	dir := "/Users/galano/Developer/Christos/Application/Go/go-development/bicep"
	dirResults, err := parser.ParseDir(dir)
	if err != nil {
		log.Fatalf("Parser error: %s", err)
	}

	fmt.Printf("\nDirectory: %s\n", dir)
	for filename, results := range dirResults {
		fmt.Printf("\nFile: %s\n", filename)
		for _, result := range results {
			printResourceInfo(result)
		}
	}
	fmt.Println()
}