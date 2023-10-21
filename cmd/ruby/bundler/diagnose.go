package bundler

import (
	"fmt"

	"github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
)

type Diagnosis struct {
	Name     string
	Url      string
	Archived bool
}

func Diagnose(r io.ReadSeekerAt) []Diagnosis {
	var diagnoses []Diagnosis
	p := &bundler.Parser{}
	libs, _, _ := p.Parse(r)

	for _, lib := range libs {
		githubUrl := FetchFromRubyGems(lib.Name)
		repo, err := github.ParseGitHubUrl(githubUrl)
		if err != nil {
			fmt.Println(err)
			continue
		}

		repo2 := github.FetchFromGitHub(repo.Owner, repo.Repo)
		diagnosis := Diagnosis{
			Name:     repo2.Repo,
			Url:      githubUrl,
			Archived: repo2.Archived,
		}
		diagnoses = append(diagnoses, diagnosis)
	}

	return diagnoses
}
