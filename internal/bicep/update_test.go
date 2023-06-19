package bicep

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/christosgalano/bruh/internal/types"
)

const (
	updateAzureDeployBicepFile      = "testdata/update/azure.deploy.bicep"
	updateAzureDeployParametersFile = "testdata/update/azure.deploy.parameters.json"
)

func TestUpdateFile(t *testing.T) {
	type args struct {
		filename       string
		resources      []types.ResourceInfo
		inPlace        bool
		includePreview bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid-file",
			args: args{
				filename: updateAzureDeployBicepFile,
				resources: []types.ResourceInfo{
					{
						ID:                "Microsoft.Resources/resourceGroups",
						Name:              "resourceGroups",
						Namespace:         "Microsoft.Resources",
						CurrentAPIVersion: "2022-09-01",
						AvailableAPIVersions: []string{
							"2022-09-01",
							"2021-04-01",
							"2021-01-01",
							"2020-10-01",
						},
					},
				},
				inPlace:        true,
				includePreview: false,
			},
			wantErr: false,
		},
		{
			name: "invalid-file",
			args: args{
				filename:       updateAzureDeployParametersFile,
				resources:      []types.ResourceInfo{},
				includePreview: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateFile(tt.args.filename, tt.args.resources, tt.args.inPlace, tt.args.includePreview); (err != nil) != tt.wantErr {
				t.Errorf("UpdateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateDirectory(t *testing.T) {
	type args struct {
		files          map[string][]types.ResourceInfo
		inPlace        bool
		includePreview bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "testdata-update",
			args: args{
				files: map[string][]types.ResourceInfo{
					updateAzureDeployBicepFile: {
						{
							ID:                "Microsoft.Resources/resourceGroups",
							Name:              "resourceGroups",
							Namespace:         "Microsoft.Resources",
							CurrentAPIVersion: "2022-09-01",
							AvailableAPIVersions: []string{
								"2022-09-01",
								"2021-04-01",
								"2021-01-01",
								"2020-10-01",
							},
						},
					},
					"testdata/update/modules/compute.bicep": {
						{
							ID:                "Microsoft.Web/serverfarms",
							Name:              "serverfarms",
							Namespace:         "Microsoft.Web",
							CurrentAPIVersion: "2021-03-01",
							AvailableAPIVersions: []string{
								"2022-03-01",
								"2021-03-01",
								"2021-02-01",
								"2021-01-15",
								"2021-01-01",
								"2020-12-01",
								"2020-10-01",
							},
						},
						{
							ID:                "Microsoft.Web/sites",
							Name:              "sites",
							Namespace:         "Microsoft.Web",
							CurrentAPIVersion: "2021-02-01",
							AvailableAPIVersions: []string{
								"2022-03-01",
								"2021-03-01",
								"2021-02-01",
								"2021-01-15",
								"2021-01-01",
								"2020-12-01",
								"2020-10-01",
							},
						},
					},
					"testdata/update/modules/identity.bicep": {
						{
							ID:                "Microsoft.ManagedIdentity/userAssignedIdentities",
							Name:              "userAssignedIdentities",
							Namespace:         "Microsoft.ManagedIdentity",
							CurrentAPIVersion: "2022-01-31-preview",
							AvailableAPIVersions: []string{
								"2023-01-31",
								"2022-01-31-preview",
								"2021-09-30-preview",
								"2018-11-30",
								"2015-08-31-preview",
							},
						},
					},
				},
				inPlace:        false,
				includePreview: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateDirectory(tt.args.files, tt.args.inPlace, tt.args.includePreview); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}

			newFiles := []string{}
			for filename := range tt.args.files {
				newFiles = append(newFiles, strings.Replace(filename, ".bicep", "_updated.bicep", 1))
			}

			for _, filename := range newFiles {
				exists, err := fileExists(filename)
				if err != nil {
					t.Errorf("fileExists() error = %v", err)
				}
				if !exists {
					t.Errorf("UpdateDirectory() error = %v", fmt.Errorf("file %s does not exist", filename))
				}

				err = deleteFile(filename)
				if err != nil {
					t.Errorf("deleteFile() error = %v", err)
				}
			}
		})
	}
}

// deleteFile deletes the given file.
func deleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("failed to delete file: %s", err)
	}
	return nil
}

// fileExists checks if the given file exists.
func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check if file exists: %s", err)
	}
	return true, nil
}
