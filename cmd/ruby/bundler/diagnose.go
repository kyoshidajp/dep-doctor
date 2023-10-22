package bundler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/ruby/bundler"
	"github.com/fatih/color"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
)

type Diagnosis struct {
	Name      string
	Url       string
	Archived  bool
	Diagnosed bool
}

func Diagnose(r io.ReadSeekerAt) []Diagnosis {
	var diagnoses []Diagnosis
	p := &bundler.Parser{}
	libs, _, _ := p.Parse(r)

	for _, lib := range libs {
		fmt.Printf("%s\n", lib.Name)

		githubUrl := FetchFromRubyGems(lib.Name)
		repo, err := github.ParseGitHubUrl(githubUrl)
		var (
			url       string
			archived  bool
			diagnosed bool
		)
		if err != nil {
			url = ""
			archived = false
			diagnosed = false
		} else {
			repo2 := github.FetchFromGitHub(repo.Owner, repo.Repo)
			url = githubUrl
			archived = repo2.Archived
			diagnosed = true
		}

		diagnosis := Diagnosis{
			Name:      lib.Name,
			Url:       url,
			Archived:  archived,
			Diagnosed: diagnosed,
		}
		diagnoses = append(diagnoses, diagnosis)
	}

	return diagnoses
}

func Report(diagnoses []Diagnosis) error {
	errMessages := []string{}
	warnMessages := []string{}
	errorCount := 0
	unDiagnosedCount := 0
	for _, diagnosis := range diagnoses {
		if !diagnosis.Diagnosed {
			warnMessages = append(warnMessages, fmt.Sprintf("[warn] %s (unknown):\n", diagnosis.Name))
			unDiagnosedCount += 1
			continue
		}
		if diagnosis.Archived {
			errMessages = append(errMessages, fmt.Sprintf("[error] %s (archived): %s\n", diagnosis.Name, diagnosis.Url))
			errorCount += 1
		}
	}

	color.Red(strings.Join(errMessages, "\n"))
	color.Yellow(strings.Join(warnMessages, "\n"))
	color.Green(heredoc.Docf(`
		Diagnose complete! %d dependencies.
		%d error, %d unknown`,
		len(diagnoses),
		errorCount,
		unDiagnosedCount),
	)

	if len(errMessages) > 0 {
		return errors.New("has error")
	}

	return nil
}
