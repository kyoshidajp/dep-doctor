package cmd

import (
	"fmt"

	"github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	"github.com/kyoshidajp/dep-doctor/cmd"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
)

const GITHUB_SEARCH_REPO_COUNT_PER_ONCE = 20

func getNameWithOwners(r io.ReadSeekerAt) []github.NameWithOwner {
	var nameWithOwners []github.NameWithOwner
	p := &bundler.Parser{}
	libs, _, _ := p.Parse(r)

	for _, lib := range libs {
		fmt.Printf("%s\n", lib.Name)

		githubUrl := cmd.FetchFromRubyGems(lib.Name)
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

type BundlerDoctor struct {
}

func (d *BundlerDoctor) Diagnose(r io.ReadSeekerAt) map[string]cmd.Diagnosis {
	diagnoses := make(map[string]cmd.Diagnosis)
	slicedNameWithOwners := [][]github.NameWithOwner{}
	nameWithOwners := getNameWithOwners(r)
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
			diagnosis := cmd.Diagnosis{
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

		diagnosis := cmd.Diagnosis{
			Name:      nameWithOwner.PackageName,
			Diagnosed: false,
		}
		diagnoses[nameWithOwner.PackageName] = diagnosis
	}
	return diagnoses
}
