package swift

import (
	"net/http"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCocoaPod_PodspecPath(t *testing.T) {
	tests := []struct {
		name string
		pod  CocoaPod
	}{
		{
			name: "KeychainAccess",
			pod: CocoaPod{
				name:    "KeychainAccess",
				version: "4.2.2",
			},
		},
		{
			name: "GoogleUtilities/Logger",
			pod: CocoaPod{
				name:    "GoogleUtilities/Logger",
				version: "7.11.0",
			},
		},
	}
	expects := []struct {
		name string
		path string
	}{
		{
			name: "KeychainAccess",
			path: "f/6/3/KeychainAccess/4.2.2/KeychainAccess.podspec.json",
		},
		{
			name: "GoogleUtilities/Logger",
			path: "0/8/4/GoogleUtilities/7.11.0/GoogleUtilities.podspec.json",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.pod.PodspecPath()
			expect := expects[i].path
			assert.Equal(t, expect, actual)
		})
	}
}

func TestCocoapod_fetchURLFromRegistry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://cdn.jsdelivr.net/cocoa/Specs/0/8/4/GoogleUtilities/7.10.0/GoogleUtilities.podspec.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"Source": {
				"Git": "https://github.com/google/GoogleUtilities.git"
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://cdn.jsdelivr.net/cocoa/Specs/f/9/9/AppCenter/4.2.0/AppCenter.podspec.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"homepage": "https://appcenter.ms"
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://cdn.jsdelivr.net/cocoa/Specs/8/0/3/not-found/1.2.3/not-found.podspec.json",
		httpmock.NewStringResponder(404, heredoc.Doc(`
		{}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://cdn.jsdelivr.net/cocoa/Specs/3/5/8/unmarshal-error/1.2.3/unmarshal-error.podspec.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"unmarshal": "xxx",
		}
		`)),
	)

	tests := []struct {
		name    string
		libName string
		version string
		wantURL string
		wantErr bool
	}{
		{
			name:    "has Source.git",
			libName: "GoogleUtilities/Environment",
			version: "7.10.0",
			wantURL: "https://github.com/google/GoogleUtilities.git",
			wantErr: false,
		},
		{
			// have no mock
			name:    "don't have Source.Git",
			libName: "not_exists",
			version: "7.10.0",
			wantURL: "",
			wantErr: true,
		},
		{
			name:    "don't have Source.Git, but have Homepage",
			libName: "AppCenter",
			version: "4.2.0",
			wantURL: "https://appcenter.ms",
			wantErr: false,
		},
		{
			name:    "404 not found",
			libName: "not-found",
			version: "1.2.3",
			wantURL: "",
			wantErr: true,
		},
		{
			name:    "Unmarshal error",
			libName: "unmarshal-error",
			version: "1.2.3",
			wantURL: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := CocoaPod{name: tt.libName, version: tt.version}
			got, err := pod.fetchURLFromRegistry(http.Client{})
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			expect := tt.wantURL
			if got != expect {
				t.Errorf("get() = %v, want starts %v", got, expect)
			}
		})
	}
}
