package github

import (
	"context"
	"fmt"
	"log"
	net_url "net/url"
	"strings"

	"github.com/google/go-github/github"
)

type GitHubRepository struct {
	Owner    string
	Repo     string
	Url      string
	Archived bool
}

func ParseGitHubUrl(url string) (GitHubRepository, error) {
	u, err := net_url.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	paths := strings.Split(u.Path, "/")
	if len(paths) < 3 {
		return GitHubRepository{}, fmt.Errorf("error: Unknown URL")
	}
	return GitHubRepository{
		Owner: paths[1],
		Repo:  paths[2],
		Url:   url,
	}, nil
}

func FetchFromGitHub(owner string, repo string) GitHubRepository {
	client := github.NewClient(nil)
	repository, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	return GitHubRepository{
		Owner:    owner,
		Repo:     repo,
		Archived: repository.GetArchived(),
	}
}
