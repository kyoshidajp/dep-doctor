package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/nodejs/yarn"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
)

const NODEJS_REGISTRY_API = "https://registry.npmjs.org/%s"

type NodejsRegistryResponse struct {
	Repository struct {
		Url string `json:"url"`
	}
}

func FetchFromNodejs(name string) string {
	url := fmt.Sprintf(NODEJS_REGISTRY_API, name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var NodejsRegistryResponse NodejsRegistryResponse
	json.Unmarshal(body, &NodejsRegistryResponse)

	return NodejsRegistryResponse.Repository.Url
}

type YarnStrategy struct {
}

func NewYarnStrategy() *YarnStrategy {
	return &YarnStrategy{}
}

func (s *YarnStrategy) Diagnose(r parser_io.ReadSeekerAt) map[string]Diagnosis {
	diagnoses := make(map[string]Diagnosis)
	slicedNameWithOwners := [][]github.NameWithOwner{}
	nameWithOwners := s.getNameWithOwners(r)
	sliceSize := len(nameWithOwners)

	for i := 0; i < sliceSize; i += GITHUB_SEARCH_REPO_COUNT_PER_ONCE {
		end := i + GITHUB_SEARCH_REPO_COUNT_PER_ONCE
		if sliceSize < end {
			end = sliceSize
		}
		slicedNameWithOwners = append(slicedNameWithOwners, nameWithOwners[i:end])
	}

	for _, nameWithOwners := range slicedNameWithOwners {
		repos := github.FetchFromGitHub(nameWithOwners)
		for _, r := range repos {
			diagnosis := Diagnosis{
				Name:      r.Name,
				Url:       r.Url,
				Archived:  r.Archived,
				Diagnosed: true,
			}
			diagnoses[r.Name] = diagnosis
		}
	}

	for _, nameWithOwner := range nameWithOwners {
		if nameWithOwner.CanSearch {
			continue
		}

		diagnosis := Diagnosis{
			Name:      nameWithOwner.PackageName,
			Diagnosed: false,
		}
		diagnoses[nameWithOwner.PackageName] = diagnosis
	}
	return diagnoses
}

func (s *YarnStrategy) getNameWithOwners(r parser_io.ReadSeekerAt) []github.NameWithOwner {
	var nameWithOwners []github.NameWithOwner
	libs, _, _ := yarn.NewParser().Parse(r)

	for _, lib := range libs {
		fmt.Printf("%s\n", lib.Name)

		githubUrl := FetchFromNodejs(lib.Name)
		repo, err := github.ParseGitHubUrl(githubUrl)

		if err != nil {
			nameWithOwners = append(nameWithOwners,
				github.NameWithOwner{
					PackageName: lib.Name,
					CanSearch:   false,
				},
			)
		} else {
			nameWithOwners = append(nameWithOwners,
				github.NameWithOwner{
					Repo:        repo.Repo,
					Owner:       repo.Owner,
					PackageName: lib.Name,
					CanSearch:   true,
				},
			)
		}
	}

	return nameWithOwners
}
