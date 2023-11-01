package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCocoaPod_PodspecPath(t *testing.T) {
	tests := []struct {
		name string
		pod  CocoaPod
	}{
		{
			name: "KeychanAccess",
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
