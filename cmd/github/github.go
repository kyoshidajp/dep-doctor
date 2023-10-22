package github

import (
	"context"
	"errors"
	"fmt"
	"log"
	net_url "net/url"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const QUERY_SEPARATOR = " "

type GitHubRepository struct {
	Name     string
	Owner    string
	Repo     string
	Url      string
	Archived bool
}

type NameWithOwner struct {
	PackageName string
	Repo        string
	Owner       string
	CanSearch   bool
}

func (n NameWithOwner) GetName() string {
	return fmt.Sprintf("repo:%s/%s", n.Owner, n.Repo)
}

func ParseGitHubUrl(url string) (GitHubRepository, error) {
	u, err := net_url.Parse(url)
	if err != nil {
		return GitHubRepository{}, errors.New("error: Unknown URL")
	}

	paths := strings.Split(u.Path, "/")
	if len(paths) < 3 {
		return GitHubRepository{}, errors.New("error: Unknown URL")
	}
	return GitHubRepository{
		Owner: paths[1],
		Repo:  paths[2],
		Url:   url,
	}, nil
}

func FetchFromGitHub(nameWithOwners []NameWithOwner) []GitHubRepository {
	token := os.Getenv("GITHUB_TOKEN")
	if len(token) == 0 {
		log.Fatal("env var `GITHUB_TOKEN` is not found")
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
					IsArchived    githubv4.Boolean
					NameWithOwner githubv4.String
					IsMirror      githubv4.Boolean
					Url           githubv4.String
					Name          githubv4.String
				} `graphql:"... on Repository"`
			}
		} `graphql:"search(query:$query, first:$count, type:REPOSITORY)"`
	}

	names := make([]string, len(nameWithOwners))
	for i, n := range nameWithOwners {
		names[i] = n.GetName()
	}
	q := strings.Join(names, QUERY_SEPARATOR)
	variables := map[string]interface{}{
		"query": githubv4.String(q),
		"count": githubv4.Int(len(names)),
	}

	client.Query(context.Background(), &query, variables)
	repos := []GitHubRepository{}
	for _, node := range query.Search.Nodes {
		repos = append(repos, GitHubRepository{
			Repo:     string(node.Repository.NameWithOwner),
			Archived: bool(node.Repository.IsArchived),
			Url:      string(node.Repository.Url),
			Name:     string(node.Repository.Name),
		})
	}

	fmt.Printf("%s", strings.Repeat(".", len(repos)))

	return repos
}
