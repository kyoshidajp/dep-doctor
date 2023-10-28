package cmd

import (
	"fmt"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
	"golang.org/x/exp/slices"
)

type BundlerDoctor struct {
}

func NewBundlerDoctor() *BundlerDoctor {
	return &BundlerDoctor{}
}

func (d *BundlerDoctor) Diagnose(r parser_io.ReadSeekerAt, year int, ignores []string) map[string]Diagnosis {
	diagnoses := make(map[string]Diagnosis)
	slicedParams := [][]github.FetchRepositoryParam{}
	fetchRepositoryParams := d.FetchRepositoryParams(r)
	sliceSize := len(fetchRepositoryParams)

	for i := 0; i < sliceSize; i += github.SEARCH_REPOS_PER_ONCE {
		end := i + github.SEARCH_REPOS_PER_ONCE
		if sliceSize < end {
			end = sliceSize
		}
		slicedParams = append(slicedParams, fetchRepositoryParams[i:end])
	}

	for _, param := range slicedParams {
		repos := github.FetchFromGitHub(param)
		for _, r := range repos {
			isIgnore := slices.Contains(ignores, r.Name)
			diagnosis := Diagnosis{
				Name:      r.Name,
				Url:       r.Url,
				Archived:  r.Archived,
				Ignored:   isIgnore,
				Diagnosed: true,
				IsActive:  r.IsActive(year),
			}
			diagnoses[r.Name] = diagnosis
		}
	}

	for _, fetchRepositoryParam := range fetchRepositoryParams {
		if fetchRepositoryParam.CanSearch {
			continue
		}

		diagnosis := Diagnosis{
			Name:      fetchRepositoryParam.PackageName,
			Diagnosed: false,
		}
		diagnoses[fetchRepositoryParam.PackageName] = diagnosis
	}
	return diagnoses
}

func (d *BundlerDoctor) FetchRepositoryParams(r parser_io.ReadSeekerAt) []github.FetchRepositoryParam {
	var params []github.FetchRepositoryParam
	p := &bundler.Parser{}
	libs, _, _ := p.Parse(r)

	rubyGems := RubyGems{}
	for _, lib := range libs {
		fmt.Printf("%s\n", lib.Name)

		githubUrl, err := rubyGems.fetchURLFromRegistry(lib.Name)
		if err != nil {
			continue
		}

		repo, err := github.ParseGitHubUrl(githubUrl)
		if err != nil {
			params = append(params,
				github.FetchRepositoryParam{
					PackageName: lib.Name,
					CanSearch:   false,
				},
			)
			continue
		}

		params = append(params,
			github.FetchRepositoryParam{
				Repo:        repo.Repo,
				Owner:       repo.Owner,
				PackageName: lib.Name,
				CanSearch:   true,
			},
		)
	}

	return params
}
