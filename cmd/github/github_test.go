package github

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/tenntenn/testtime"
)

func TestGitHubRepository_IsActive(t *testing.T) {
	cases := []struct {
		name            string
		now             time.Time
		lastCommittedAt time.Time
		year            int
		want            bool
	}{
		{
			"active",
			time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
			5,
			true,
		},
		{
			"not active",
			time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
			5,
			false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			g := &GitHubRepository{
				LastCommittedAt: tt.lastCommittedAt,
			}
			testtime.SetTime(t, tt.now)
			got := g.IsActive(tt.year)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseGitHubUrl(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "github.com repository",
			url:  "github.com/bvaughn/highlight-words-core.git",
		},
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
		"github.com repository": {
			Owner: "bvaughn",
			Repo:  "highlight-words-core",
		},
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
			got := FetchFromGitHub(tt.nameWithOwners)
			if d := cmp.Diff(got, expect, cmpopts.IgnoreFields(GitHubRepository{}, "LastCommittedAt")); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}
