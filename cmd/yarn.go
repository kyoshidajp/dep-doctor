package cmd

import (
	"fmt"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/nodejs/yarn"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
)

type YarnDoctor struct {
}

func NewYarnDoctor() *YarnDoctor {
	return &YarnDoctor{}
}

func (d *YarnDoctor) Diagnose(r parser_io.ReadSeekerAt, year int, ignores []string) map[string]Diagnosis {
	diagnoses := make(map[string]Diagnosis)
	slicedParams := [][]github.FetchRepositoryParam{}
	params := d.FetchRepositoryParams(r)
	sliceSize := len(slicedParams)

	for i := 0; i < sliceSize; i += github.SEARCH_REPOS_PER_ONCE {
		end := i + github.SEARCH_REPOS_PER_ONCE
		if sliceSize < end {
			end = sliceSize
		}
		slicedParams = append(slicedParams, params[i:end])
	}

	for _, slicedParam := range slicedParams {
		repos := github.FetchFromGitHub(slicedParam)
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

	for _, param := range params {
		if param.CanSearch {
			continue
		}

		diagnosis := Diagnosis{
			Name:      param.PackageName,
			Diagnosed: false,
		}
		diagnoses[param.PackageName] = diagnosis
	}
	return diagnoses
}

func (d *YarnDoctor) FetchRepositoryParams(r parser_io.ReadSeekerAt) []github.FetchRepositoryParam {
	var params []github.FetchRepositoryParam
	libs, _, _ := yarn.NewParser().Parse(r)

	nodejs := Nodejs{}
	for _, lib := range libs {
		fmt.Printf("%s\n", lib.Name)

		githubUrl, err := nodejs.fetchURLFromRegistry(lib.Name)
		if err != nil {
			params = append(params,
				github.FetchRepositoryParam{
					PackageName: lib.Name,
					CanSearch:   false,
				},
			)
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
