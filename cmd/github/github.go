package github

import (
	"context"
	"errors"
	"fmt"
	net_url "net/url"
	"os"
	"strings"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GitHubRepository struct {
	Owner    string
	Repo     string
	Url      string
	Archived bool
}

type NameWithOwner struct {
	Repo  string
	Owner string
}

func (n NameWithOwner) getName() string {
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
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
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
				} `graphql:"... on Repository"`
			}
		} `graphql:"search(query:$query, first:$count, type:REPOSITORY)"`
	}

	names := make([]string, len(nameWithOwners))
	for i, n := range nameWithOwners {
		names[i] = n.getName()
	}
	q := strings.Join(names, " ")
	variables := map[string]interface{}{
		"query": githubv4.String(q),
		"count": githubv4.NewInt(2),
	}

	err := client.Query(context.Background(), &query, variables)
	if err != nil {
		// Handle error.
	}

	repos := []GitHubRepository{}
	for _, node := range query.Search.Nodes {
		repos = append(repos, GitHubRepository{
			Owner:    string(node.Repository.NameWithOwner),
			Repo:     string(node.Repository.NameWithOwner),
			Archived: false,
		})
	}

	return repos
}
