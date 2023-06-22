package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/christosgalano/bruh/internal/types"
	"github.com/olekukonko/tablewriter"
)

// printFileNormal prints the file's information in normal format.
func printFileNormal(bicepFile *types.BicepFile, filename string, outdated bool, mode types.Mode) {
	fmt.Printf("\n%s:\n", filename)
	for _, resource := range bicepFile.Resources {
		latestAPIVersion := resource.AvailableAPIVersions[0]
		if mode == types.ModeScan {
			if resource.CurrentAPIVersion != latestAPIVersion {
				fmt.Printf("  - %s is using %s while the latest version is %s\n", resource.ID, resource.CurrentAPIVersion, latestAPIVersion)
			} else if !outdated {
				fmt.Printf("  + %s is using the latest version %s\n", resource.ID, resource.CurrentAPIVersion)
			}
		} else {
			fmt.Printf("  + Updated %s from version %s to %s\n", resource.ID, resource.CurrentAPIVersion, latestAPIVersion)
		}
	}
	fmt.Println()
}

// printFileTable prints the file's information in tabular format.
func printFileTable(bicepFile *types.BicepFile, outdated bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Resource", "Current API Version", "Latest API Version"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

	fmt.Printf("\n%s:\n", bicepFile.Name)
	for _, resource := range bicepFile.Resources {
		if outdated && (resource.CurrentAPIVersion == resource.AvailableAPIVersions[0]) {
			continue
		}
		table.Append([]string{resource.ID, resource.CurrentAPIVersion, resource.AvailableAPIVersions[0]})
	}
	table.Render()
	fmt.Println()
}

// printDirectoryNormal prints the directory's information in normal format.
func printDirectoryNormal(bicepDirectory *types.BicepDirectory, outdated bool, mode types.Mode) {
	absolutePath, err := filepath.Abs(bicepDirectory.Name)
	if err != nil {
		absolutePath = bicepDirectory.Name
	}
	fmt.Printf("\n%s:\n", absolutePath)
	for _, file := range bicepDirectory.Files {
		filename, err := filepath.Rel(bicepDirectory.Name, file.Name)
		if err != nil {
			filename = file.Name
		}
		printFileNormal(&file, filename, outdated, mode)
	}
}

// printDirectoryTable prints the directory's information in tabular format.
func printDirectoryTable(bicepDirectory *types.BicepDirectory, outdated bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File", "Resource", "Current API Version", "Latest API Version"})
	table.SetBorders(tablewriter.Border{Left: true, Top: true, Right: true, Bottom: true})
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)

	absolutePath, err := filepath.Abs(bicepDirectory.Name)
	if err != nil {
		absolutePath = bicepDirectory.Name
	}
	fmt.Printf("\n%s:\n", absolutePath)
	for _, file := range bicepDirectory.Files {
		for _, resource := range file.Resources {
			filename, err := filepath.Rel(bicepDirectory.Name, file.Name)
			if err != nil {
				filename = file.Name
			}
			if outdated && (resource.CurrentAPIVersion == resource.AvailableAPIVersions[0]) {
				continue
			}
			table.Append([]string{filename, resource.ID, resource.CurrentAPIVersion, resource.AvailableAPIVersions[0]})
		}
	}
	table.Render()
	fmt.Println()
}