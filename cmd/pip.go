package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/python/pip"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
)

// https://warehouse.pypa.io/api-reference/json.html
const PYPI_REGISTRY_API = "https://pypi.org/pypi/%s/json"

type PipRegistryResponse struct {
	Info struct {
		ProjectUrls struct {
			SourceCode string `json:"Source Code"`
			Source     string `json:"Source"`
		} `json:"project_urls"`
	} `json:"info"`
}

type PipDoctor struct {
}

func NewPipDoctor() *PipDoctor {
	return &PipDoctor{}
}

func (d *PipDoctor) fetchURLFromRepository(name string) (string, error) {
	url := fmt.Sprintf(PYPI_REGISTRY_API, name)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	client := new(http.Client)
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)

	var PipRegistryResponse PipRegistryResponse
	err := json.Unmarshal(body, &PipRegistryResponse)
	if err != nil {
		return "", nil
	}

	if PipRegistryResponse.Info.ProjectUrls.SourceCode != "" {
		return PipRegistryResponse.Info.ProjectUrls.SourceCode, nil
	}

	return PipRegistryResponse.Info.ProjectUrls.Source, nil
}

func (d *PipDoctor) Diagnose(r parser_io.ReadSeekerAt, year int) map[string]Diagnosis {
	diagnoses := make(map[string]Diagnosis)
	slicedNameWithOwners := [][]github.NameWithOwner{}
	nameWithOwners := d.NameWithOwners(r)
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

func (d *PipDoctor) NameWithOwners(r parser_io.ReadSeekerAt) []github.NameWithOwner {
	var nameWithOwners []github.NameWithOwner
	libs, _, _ := pip.NewParser().Parse(r)

	for _, lib := range libs {
		fmt.Printf("%s\n", lib.Name)

		githubUrl, err := d.fetchURLFromRepository(lib.Name)
		if err != nil {
			nameWithOwners = append(nameWithOwners,
				github.NameWithOwner{
					PackageName: lib.Name,
					CanSearch:   false,
				},
			)
			continue
		}

		repo, err := github.ParseGitHubUrl(githubUrl)
		if err != nil {
			nameWithOwners = append(nameWithOwners,
				github.NameWithOwner{
					PackageName: lib.Name,
					CanSearch:   false,
				},
			)
			continue
		}

		nameWithOwners = append(nameWithOwners,
			github.NameWithOwner{
				Repo:        repo.Repo,
				Owner:       repo.Owner,
				PackageName: lib.Name,
				CanSearch:   true,
			},
		)
	}

	return nameWithOwners
}
