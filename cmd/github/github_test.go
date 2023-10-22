package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGitHubUrl(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "active repository",
			url:  "https://github.com/rails/thor/tree/v1.3.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := ParseGitHubUrl(tt.url)
			assert.Equal(t, "rails", r.Owner)
			assert.Equal(t, "thor", r.Repo)
		})
	}
}

func TestFetchFromGitHub(t *testing.T) {
	tests := []struct {
		name           string
		nameWithOwners []NameWithOwner
	}{
		{
			name: "active repository",
			nameWithOwners: []NameWithOwner{
				{
					Owner: "rails",
					Repo:  "rails",
				},
				{
					Owner: "rails",
					Repo:  "strong_parameters",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := FetchFromGitHub(tt.nameWithOwners)
			assert.Equal(t, "test", r)
		})
	}
}
