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

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "starts github.com",
			url:  "github.com/bvaughn/highlight-words-core.git",
		},
		{
			name: "starts https://",
			url:  "https://github.com/rails/thor/tree/v1.3.0",
		},
		{
			name: "starts git@",
			url:  "git@github.com:rails/thor.git",
		},
		{
			name: "starts git+https://",
			url:  "git+https://github.com/then/promise.git",
		},
		{
			name: "starts git://",
			url:  "git://github.com/es-shims/typedarray.git",
		},
		{
			name: "starts git+ssh://",
			url:  "git+ssh://git@github.com/DABH/colors.js.git",
		},
	}

	expects := map[string]GitHubRepository{
		"starts github.com": {
			Owner: "bvaughn",
			Repo:  "highlight-words-core",
		},
		"starts https://": {
			Owner: "rails",
			Repo:  "thor",
		},
		"starts git@": {
			Owner: "rails",
			Repo:  "thor",
		},
		"starts git+https://": {
			Owner: "then",
			Repo:  "promise",
		},
		"starts git://": {
			Owner: "es-shims",
			Repo:  "typedarray",
		},
		"starts git+ssh://": {
			Owner: "DABH",
			Repo:  "colors.js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := ParseGitHubURL(tt.url)
			expect := expects[tt.name]
			assert.Equal(t, expect.Owner, r.Owner)
			assert.Equal(t, expect.Repo, r.Repo)
		})
	}
}

func TestFetchFromGitHub(t *testing.T) {
	tests := []struct {
		name   string
		params []FetchRepositoryParam
	}{
		{
			name: "active repository",
			params: []FetchRepositoryParam{
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
			URL:      "https://github.com/rails/rails",
			Archived: false,
		},
		{
			Name:     "strong_parameters",
			Repo:     "rails/strong_parameters",
			URL:      "https://github.com/rails/strong_parameters",
			Archived: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FetchFromGitHub(tt.params)
			if d := cmp.Diff(got, expect, cmpopts.IgnoreFields(GitHubRepository{}, "LastCommittedAt")); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}
