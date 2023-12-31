package bicep

import (
	"reflect"
	"strings"
	"testing"

	"github.com/christosgalano/bruh/internal/types"
)

func TestParseUpdateFile(t *testing.T) {
	type args struct {
		filename  string
		resources []types.Resource
		inPlace   bool
	}
	tests := []struct {
		name    string
		args    args
		initial types.BicepFile
		final   types.BicepFile
		wantErr bool
	}{
		{
			name: "azure.deploy.bicep",
			args: args{
				filename: "testdata/parse_update/azure.deploy.bicep",
				resources: []types.Resource{
					{
						AvailableAPIVersions: []string{
							"2022-09-01",
							"2021-04-01",
							"2021-01-01",
							"2020-10-01",
						},
					},
				},
				inPlace: false,
			},
			initial: types.BicepFile{
				Path: "testdata/parse_update/azure.deploy.bicep",
				Resources: []types.Resource{
					{
						ID:                "Microsoft.Resources/resourceGroups",
						Name:              "resourceGroups",
						Namespace:         "Microsoft.Resources",
						CurrentAPIVersion: "2021-01-01",
					},
				},
			},
			final: types.BicepFile{
				Path: "testdata/parse_update/azure.deploy_updated.bicep",
				Resources: []types.Resource{
					{
						ID:                "Microsoft.Resources/resourceGroups",
						Name:              "resourceGroups",
						Namespace:         "Microsoft.Resources",
						CurrentAPIVersion: "2022-09-01",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First parse
			got, err := ParseFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got.Resources, tt.initial.Resources) {
				t.Fatalf("First parse:\nParseFile() = %v, want %v", got.Resources, tt.initial.Resources)
			}

			// Inject available API versions
			for i := range got.Resources {
				got.Resources[i].AvailableAPIVersions = tt.args.resources[i].AvailableAPIVersions
			}

			// Update file
			err = UpdateFile(got, tt.args.inPlace)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UpdateFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Second parse - not in-place
			updatedFile := strings.Replace(tt.args.filename, ".bicep", "_updated.bicep", 1)
			got, err = ParseFile(updatedFile)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got.Resources, tt.final.Resources) {
				t.Fatalf("Second parse:\nParseFile() = %v, want %v", got.Resources, tt.final.Resources)
			}

			// Cleanup
			err = deleteFile(updatedFile)
			if err != nil {
				t.Fatalf("deleteFile() error = %v", err)
			}
		})
	}
}
