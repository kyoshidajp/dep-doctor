package github

import (
	"fmt"
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
		name           string
		URL            string
		wantRepository GitHubRepository
		wantErr        error
	}{
		{
			name: "starts github.com",
			URL:  "github.com/bvaughn/highlight-words-core.git",
			wantRepository: GitHubRepository{
				Owner: "bvaughn",
				Repo:  "highlight-words-core",
				URL:   "github.com/bvaughn/highlight-words-core.git",
			},
			wantErr: nil,
		},
		{
			name: "starts https://",
			URL:  "https://github.com/rails/thor/tree/v1.3.0",
			wantRepository: GitHubRepository{
				Owner: "rails",
				Repo:  "thor",
				URL:   "https://github.com/rails/thor/tree/v1.3.0",
			},
			wantErr: nil,
		},
		{
			name: "starts git@",
			URL:  "git@github.com:rails/thor.git",
			wantRepository: GitHubRepository{
				Owner: "rails",
				Repo:  "thor",
				URL:   "git@github.com:rails/thor.git",
			},
			wantErr: nil,
		},
		{
			name: "starts git+https://",
			URL:  "git+https://github.com/then/promise.git",
			wantRepository: GitHubRepository{
				Owner: "then",
				Repo:  "promise",
				URL:   "git+https://github.com/then/promise.git",
			},
			wantErr: nil,
		},
		{
			name: "starts git://",
			URL:  "git://github.com/es-shims/typedarray.git",
			wantRepository: GitHubRepository{
				Owner: "es-shims",
				Repo:  "typedarray",
				URL:   "git://github.com/es-shims/typedarray.git",
			},
			wantErr: nil,
		},
		{
			name: "starts git+ssh://",
			URL:  "git+ssh://git@github.com/DABH/colors.js.git",
			wantRepository: GitHubRepository{
				Owner: "DABH",
				Repo:  "colors.js",
				URL:   "git+ssh://git@github.com/DABH/colors.js.git",
			},
			wantErr: nil,
		},
		{
			name: "blank",
			URL:  "",
			wantRepository: GitHubRepository{
				Owner: "",
				Repo:  "",
				URL:   "",
			},
			wantErr: fmt.Errorf("source code URL is blank"),
		},
		{
			name: "invalid URL",
			URL:  "https://example.com",
			wantRepository: GitHubRepository{
				Owner: "",
				Repo:  "",
				URL:   "",
			},
			wantErr: fmt.Errorf("unknown source code URL: https://example.com"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGitHubURL(tt.URL)
			wantRepository := tt.wantRepository
			assert.Equal(t, wantRepository.Owner, got.Owner)
			assert.Equal(t, wantRepository.Repo, got.Repo)
			assert.Equal(t, tt.wantErr, err)
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
			Owner:    "rails",
			Repo:     "rails",
			URL:      "https://github.com/rails/rails",
			Archived: false,
		},
		{
			Name:     "strong_parameters",
			Owner:    "rails",
			Repo:     "strong_parameters",
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
