package swift

import (
	"net/http"
	"testing"

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
	tests := []struct {
		name string
		pod  CocoaPod
	}{
		{
			name: "has Source.git",
			pod: CocoaPod{
				name:    "GoogleUtilities/Environment",
				version: "7.10.0",
			},
		},
		{
			name: "don't have Source.Git",
			pod: CocoaPod{
				name:    "not_exists",
				version: "7.10.0",
			},
		},
		{
			name: "don't have Source.Git, but have Homepage",
			pod: CocoaPod{
				name:    "AppCenter",
				version: "4.2.0",
			},
		},
	}
	expects := []struct {
		name string
		url  string
	}{
		{
			name: "has Source.git",
			url:  "https://github.com/google/GoogleUtilities.git",
		},
		{
			name: "don't have Source.Git",
			url:  "",
		},
		{
			name: "don't have Source.Git, but have Homepage",
			url:  "https://appcenter.ms",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := tt.pod.fetchURLFromRegistry(http.Client{})
			expected := expects[i].url
			assert.Equal(t, expected, actual)
		})
	}
}
