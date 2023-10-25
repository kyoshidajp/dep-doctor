package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
)

const GITHUB_SEARCH_REPO_COUNT_PER_ONCE = 20

// https://guides.rubygems.org/rubygems-org-api/
const RUBYGEMS_ORG_API = "https://rubygems.org/api/v1/gems/%s.json"

type GemResponse struct {
	Name          string `json:"name"`
	SourceCodeUri string `json:"source_code_uri"`
	HomepageUri   string `json:"homepage_uri"`
}

func FetchFromRubyGems(name string) string {
	url := fmt.Sprintf(RUBYGEMS_ORG_API, name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var Gem GemResponse
	json.Unmarshal(body, &Gem)

	if Gem.SourceCodeUri != "" {
		return Gem.SourceCodeUri
	} else if Gem.HomepageUri != "" {
		return Gem.HomepageUri
	}

	return ""
}

type BundlerStrategy struct {
}

func NewBundlerStrategy() *BundlerStrategy {
	return &BundlerStrategy{}
}

func (s *BundlerStrategy) Diagnose(r parser_io.ReadSeekerAt, year int) map[string]Diagnosis {
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
				IsActive:  r.IsActive(year),
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

func (s *BundlerStrategy) getNameWithOwners(r parser_io.ReadSeekerAt) []github.NameWithOwner {
	var nameWithOwners []github.NameWithOwner
	p := &bundler.Parser{}
	libs, _, _ := p.Parse(r)

	for _, lib := range libs {
		fmt.Printf("%s\n", lib.Name)

		githubUrl := FetchFromRubyGems(lib.Name)
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
