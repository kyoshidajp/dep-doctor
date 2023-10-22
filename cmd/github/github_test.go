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
			name: "http scheme repository",
			url:  "https://github.com/rails/thor/tree/v1.3.0",
		},
		{
			name: "git scheme repository",
			url:  "git@github.com:rails/thor.git",
		},
		{
			name: "git+https scheme repository",
			url:  "git+https://github.com/then/promise.git",
		},
	}

	expects := map[string]GitHubRepository{
		"http scheme repository": {
			Owner: "rails",
			Repo:  "thor",
		},
		"git scheme repository": {
			Owner: "rails",
			Repo:  "thor",
		},
		"git+https scheme repository": {
			Owner: "then",
			Repo:  "promise",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := ParseGitHubUrl(tt.url)
			expect := expects[tt.name]
			assert.Equal(t, expect.Owner, r.Owner)
			assert.Equal(t, expect.Repo, r.Repo)
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

	expect := []GitHubRepository{
		{
			Name:     "rails",
			Repo:     "rails/rails",
			Url:      "https://github.com/rails/rails",
			Archived: false,
		},
		{
			Name:     "strong_parameters",
			Repo:     "rails/strong_parameters",
			Url:      "https://github.com/rails/strong_parameters",
			Archived: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := FetchFromGitHub(tt.nameWithOwners)
			assert.Equal(t, expect, r)
		})
	}
}
