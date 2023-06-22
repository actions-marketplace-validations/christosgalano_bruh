/*
Package cli TODO: add description
*/
package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/christosgalano/bruh/internal/apiversions"
	"github.com/christosgalano/bruh/internal/bicep"
	"github.com/christosgalano/bruh/internal/types"
	"github.com/spf13/cobra"
)

var (
	updatePath     string
	inPlace        bool
	includePreview bool
	silent         bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a bicep file or a directory containing bicep files",
	Long: `Update a bicep file or a directory containing bicep files
so that each Azure resource uses the latest API version available.
It is possible to update the files in place or create new files with "_updated.bicep" extension.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Invalid path
		fs, err := os.Stat(updatePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "Error: no such file or directory %q\n", updatePath)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(1)
		}

		// Save stdout and stderr
		originalStdout := os.Stdout
		originalStderr := os.Stderr

		// Silent mode - close stdout and stderr
		if silent {
			err = os.Stdout.Close()
			err = os.Stderr.Close()
		}

		// Update file or directory
		if fs.IsDir() {
			err = updateDirectory(updatePath, inPlace, includePreview)
		} else {
			err = updateFile(updatePath, inPlace, includePreview)
		}

		// Restore stdout and stderr
		if silent {
			os.Stdout = originalStdout
			os.Stderr = originalStderr
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Local flags

	// path - required
	updateCmd.Flags().StringVarP(&updatePath, "path", "p", "", "path to bicep file or directory containing bicep files")
	updateCmd.MarkFlagRequired("path")

	// in-place - optional
	updateCmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "update the bicep files in place (if not set: create new files with \"_updated.bicep\" extension)")

	// include-preview - optional
	updateCmd.Flags().BoolVarP(&includePreview, "include-preview", "r", false, "include preview API versions (if not set: only non-preview versions will be considered)")

	// silent - optional
	updateCmd.Flags().BoolVarP(&silent, "silent", "s", false, "silent mode (no output)")
}

// updateFile parses the given file, fetches the latest API versions for each Azure resource and updates the file.
// If inPlace is true, the file will be updated in place; otherwise, a new file with "_updated.bicep" extension will be created.
// If includePreview is true, preview API versions will be included; otherwise, only non-preview versions will be considered.
func updateFile(path string, inPlace bool, includePreview bool) error {
	bicepFile, err := bicep.ParseFile(path)
	if err != nil {
		return err
	}

	err = apiversions.UpdateBicepFile(bicepFile)
	if err != nil {
		return err
	}

	err = bicep.UpdateFile(bicepFile, inPlace, includePreview)
	if err != nil {
		return err
	}

	printFileNormal(bicepFile, path, outdated, types.ModeUpdate)
	return nil
}

// updateDirectory parses the given directory, fetches the latest API versions for each Azure resource and updates each file.
// If inPlace is true, the files will be updated in place; otherwise, new files with "_updated.bicep" extension will be created.
// If includePreview is true, preview API versions will be included; otherwise, only non-preview versions will be considered.
func updateDirectory(path string, inPlace bool, includePreview bool) error {
	bicepDirectory, err := bicep.ParseDirectory(path)
	if err != nil {
		return err
	}

	err = apiversions.UpdateBicepDirectory(bicepDirectory)
	if err != nil {
		return err
	}

	err = bicep.UpdateDirectory(bicepDirectory, inPlace, includePreview)
	if err != nil {
		return err
	}

	printDirectoryNormal(bicepDirectory, outdated, types.ModeUpdate)
	return nil
}