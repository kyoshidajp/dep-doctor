package github

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"
	"github.com/tenntenn/testtime"
	giturls "github.com/whilp/git-urls"
	"golang.org/x/oauth2"
)

const QUERY_SEPARATOR = " "
const SEARCH_REPOS_PER_ONCE = 20
const TOKEN_NAME = "GITHUB_TOKEN"

type GitHubRepository struct {
	Name            string
	Owner           string
	Repo            string
	Url             string
	Archived        bool
	LastCommittedAt time.Time
}

func (r GitHubRepository) IsActive(year int) bool {
	now := testtime.Now()
	targetDate := r.LastCommittedAt.AddDate(year, 0, 0)
	return targetDate.After(now)
}

type FetchRepositoryParam struct {
	PackageName string
	Repo        string
	Owner       string
	CanSearch   bool
}

func (n FetchRepositoryParam) QueryWord() string {
	return fmt.Sprintf("repo:%s/%s", n.Owner, n.Repo)
}

type GitHubURL struct {
	URL string
}

func (githuburl GitHubURL) Parse() (string, string, error) {
	if githuburl.URL == "" {
		return "", "", errors.New("error: Blank URL")
	}

	u, err := giturls.Parse(githuburl.URL)
	if err != nil {
		return "", "", errors.New("error: Unknown URL")
	}

	var owner, repo string
	paths := strings.Split(u.Path, "/")
	if u.Scheme == "ssh" {
		owner = paths[0]
		repo = strings.Replace(paths[1], ".git", "", 1)
	} else if u.Scheme == "https" || u.Scheme == "http" {
		if len(paths) < 3 {
			return "", "", errors.New("error: Unknown URL")
		}
		owner = paths[1]
		repo = paths[2]
	} else if u.Scheme == "file" {
		if paths[0] == "github.com" {
			owner = paths[1]
			repo = strings.Replace(paths[2], ".git", "", 1)
		} else {
			owner = paths[3]
			repo = strings.Replace(paths[4], ".git", "", 1)
		}
	}
	return owner, repo, nil
}

func ParseGitHubUrl(url string) (GitHubRepository, error) {
	githubURL := GitHubURL{
		URL: url,
	}
	owner, repo, err := githubURL.Parse()
	if err != nil {
		return GitHubRepository{}, err
	}
	return GitHubRepository{
		Owner: owner,
		Repo:  repo,
		Url:   url,
	}, nil
}

func FetchFromGitHub(params []FetchRepositoryParam) []GitHubRepository {
	token := os.Getenv(TOKEN_NAME)
	if len(token) == 0 {
		m := fmt.Sprintf("Environment variable `%s` is not found.", TOKEN_NAME)
		log.Fatal(m)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	var query struct {
		Search struct {
			RepositoryCount githubv4.Int
			Nodes           []struct {
				Repository struct {
					IsArchived       githubv4.Boolean
					NameWithOwner    githubv4.String
					IsMirror         githubv4.Boolean
					Url              githubv4.String
					Name             githubv4.String
					DefaultBranchRef struct {
						Target struct {
							Commit struct {
								History struct {
									Edges []struct {
										Node struct {
											CommittedDate githubv4.DateTime
										}
									}
								} `graphql:"history(first:1)"`
							} `graphql:"... on Commit"`
						}
					}
				} `graphql:"... on Repository"`
			}
		} `graphql:"search(query:$query, first:$count, type:REPOSITORY)"`
	}

	names := make([]string, len(params))
	for i, param := range params {
		names[i] = param.QueryWord()
	}
	q := strings.Join(names, QUERY_SEPARATOR)
	variables := map[string]interface{}{
		"query": githubv4.String(q),
		"count": githubv4.Int(len(names)),
	}

	repos := []GitHubRepository{}
	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		return repos
	}

	for _, node := range query.Search.Nodes {
		nodeRepo := node.Repository
		lastCommit := nodeRepo.DefaultBranchRef.Target.Commit.History.Edges[0].Node
		repos = append(repos, GitHubRepository{
			Repo:            string(nodeRepo.NameWithOwner),
			Archived:        bool(nodeRepo.IsArchived),
			Url:             string(nodeRepo.Url),
			Name:            string(nodeRepo.Name),
			LastCommittedAt: time.Time(lastCommit.CommittedDate.Time),
		})
	}

	fmt.Printf("%s", strings.Repeat(".", len(repos)))

	return repos
}
