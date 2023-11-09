package github

import (
	"fmt"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jarcoal/httpmock"
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

func TestGitHubRepository_RepoOwner(t *testing.T) {
	tests := []struct {
		name string
		repo GitHubRepository
		want string
	}{
		{
			name: "normal",
			repo: GitHubRepository{
				Owner: "owner",
				Repo:  "repo",
			},
			want: "owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.repo.RepoOwner()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFetchRepositoryParam_RepoOwner(t *testing.T) {
	tests := []struct {
		name string
		repo FetchRepositoryParam
		want string
	}{
		{
			name: "normal",
			repo: FetchRepositoryParam{
				Owner: "owner",
				Repo:  "repo",
			},
			want: "owner/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.repo.RepoOwner()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGitHubURL_Parse(t *testing.T) {
	tests := []struct {
		name      string
		url       GitHubURL
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name: "normal",
			url: GitHubURL{
				URL: "https://github.com/kyoshidajp/dep-doctor",
			},
			wantOwner: "kyoshidajp",
			wantRepo:  "dep-doctor",
			wantErr:   false,
		},
		{
			name: "blank URL",
			url: GitHubURL{
				URL: "",
			},
			wantOwner: "",
			wantRepo:  "",
			wantErr:   true,
		},
		{
			name: "unsupported schema",
			url: GitHubURL{
				URL: "ftp://example.com/test1/test2/test3",
			},
			wantOwner: "",
			wantRepo:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOwner, gotRepo, err := tt.url.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			expectOwner := tt.wantOwner
			if gotOwner != expectOwner {
				t.Errorf("get() = %v, want %v", gotOwner, expectOwner)
			}
			expectRepo := tt.wantRepo
			if gotRepo != expectRepo {
				t.Errorf("get() = %v, want %v", gotRepo, expectRepo)
			}
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
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.github.com/graphql",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"data": {
			  "search": {
				"repositoryCount": 2,
				"nodes": [
				  {
					"isArchived": false,
					"url": "https://github.com/rails/rails",
					"name": "rails",
					"owner": {
					  "login": "rails"
					},
					"defaultBranchRef": {
					  "target": {
						"history": {
						  "edges": [
							{
							  "node": {
								"committedDate": "2023-11-13T18:38:38Z"
							  }
							}
						  ]
						}
					  }
					}
				  },
				  {
					"isArchived": true,
					"url": "https://github.com/rails/strong_parameters",
					"name": "strong_parameters",
					"owner": {
					  "login": "rails"
					},
					"defaultBranchRef": {
					  "target": {
						"history": {
						  "edges": [
							{
							  "node": {
								"committedDate": "2017-08-08T18:36:30Z"
							  }
							}
						  ]
						}
					  }
					}
				  }
				]
			  }
			}
		  }
		`)),
	)

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

func TestFetchFromGitHub_NotFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "https://api.github.com/graphql",
		httpmock.NewStringResponder(404, heredoc.Doc(`
		{}
		`)),
	)

	tests := []struct {
		name   string
		params []FetchRepositoryParam
	}{
		{
			name: "active repository",
			params: []FetchRepositoryParam{
				{
					PackageName: "kyoshidajp/not-found",
					Owner:       "kyoshidajp",
					Repo:        "not-found",
				},
			},
		},
	}

	expect := []GitHubRepository{
		{
			Name: "kyoshidajp/not-found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FetchFromGitHub(tt.params)
			if d := cmp.Diff(got, expect, cmpopts.IgnoreFields(GitHubRepository{}, "LastCommittedAt", "Error")); len(d) != 0 {
				t.Errorf("differs: (-got +want)\n%s", d)
			}
		})
	}
}
