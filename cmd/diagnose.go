package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/MakeNowJust/heredoc"
	"github.com/aquasecurity/go-dep-parser/pkg/io"
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/fatih/color"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

const MAX_YEAR_TO_BE_BLANK = 5

// referenced as the number of goroutine parallels
// should be optimized?
const FETCH_REPOS_PER_ONCE = 20

type Diagnosis struct {
	Name      string
	Url       string
	Archived  bool
	Ignored   bool
	Diagnosed bool
	IsActive  bool
}

type MedicalTechnician interface {
	Deps(r parser_io.ReadSeekerAt) []types.Library
	SourceCodeURL(name string) (string, error)
}

func FetchRepositoryParams(libs []types.Library, g MedicalTechnician) []github.FetchRepositoryParam {
	var params []github.FetchRepositoryParam
	maxConcurrency := FETCH_REPOS_PER_ONCE
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)

	for _, lib := range libs {
		wg.Add(1)
		sem <- struct{}{}
		go func(lib types.Library) {
			defer wg.Done()
			defer func() { <-sem }()

			fmt.Printf("%s\n", lib.Name)

			githubUrl, err := g.SourceCodeURL(lib.Name)
			if err != nil {
				return
			}

			repo, err := github.ParseGitHubUrl(githubUrl)
			if err != nil {
				params = append(params,
					github.FetchRepositoryParam{
						PackageName: lib.Name,
						CanSearch:   false,
					},
				)
				return
			}

			params = append(params,
				github.FetchRepositoryParam{
					Repo:        repo.Repo,
					Owner:       repo.Owner,
					PackageName: lib.Name,
					CanSearch:   true,
				},
			)
		}(lib)

		wg.Wait()
	}

	return params
}

func Diagnose(d MedicalTechnician, r io.ReadSeekCloserAt, year int, ignores []string) map[string]Diagnosis {
	diagnoses := make(map[string]Diagnosis)
	slicedParams := [][]github.FetchRepositoryParam{}
	deps := d.Deps(r)
	fetchRepositoryParams := FetchRepositoryParams(deps, d)
	sliceSize := len(fetchRepositoryParams)

	for i := 0; i < sliceSize; i += github.SEARCH_REPOS_PER_ONCE {
		end := i + github.SEARCH_REPOS_PER_ONCE
		if sliceSize < end {
			end = sliceSize
		}
		slicedParams = append(slicedParams, fetchRepositoryParams[i:end])
	}

	maxConcurrency := FETCH_REPOS_PER_ONCE
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrency)
	for _, param := range slicedParams {
		wg.Add(1)
		sem <- struct{}{}
		go func(param []github.FetchRepositoryParam) {
			defer wg.Done()
			defer func() { <-sem }()

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
		}(param)
	}

	wg.Wait()

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

type Options struct {
	packageManager string
	lockFilePath   string
	ignores        string
	year           int
}

func (o *Options) Ignores() []string {
	return strings.Split(o.ignores, " ")
}

var (
	o = &Options{}
)

var doctors = map[string]MedicalTechnician{
	"bundler": NewBundlerDoctor(),
	"yarn":    NewYarnDoctor(),
	"pip":     NewPipDoctor(),
	"npm":     NewNPMDoctor(),
}

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Diagnose dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		doctor, ok := doctors[o.packageManager]
		if !ok {
			packages := []string{}
			for p := range doctors {
				packages = append(packages, p)
			}
			m := fmt.Sprintf("Unknown package manager: %s. You can choose from [%s]", o.packageManager, strings.Join(packages, ", "))
			log.Fatal(m)
		}

		lockFilePath := o.lockFilePath
		f, err := os.Open(lockFilePath)
		defer func() {
			_ = f.Close()
		}()
		if err != nil {
			m := fmt.Sprintf("Can't open: %s.", o.lockFilePath)
			log.Fatal(m)
		}

		diagnoses := Diagnose(doctor, f, o.year, o.Ignores())
		if err := Report(diagnoses); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
	diagnoseCmd.Flags().StringVarP(&o.packageManager, "package", "p", "bundler", "package manager")
	diagnoseCmd.Flags().StringVarP(&o.lockFilePath, "lock_file", "f", "Gemfile.lock", "lock file path")
	diagnoseCmd.Flags().StringVarP(&o.ignores, "ignores", "i", "", "ignore dependencies")
	diagnoseCmd.Flags().IntVarP(&o.year, "year", "y", MAX_YEAR_TO_BE_BLANK, "max years of inactivity")
}

func Report(diagnoses map[string]Diagnosis) error {
	errMessages := []string{}
	warnMessages := []string{}
	ignoredMessages := []string{}
	errCount := 0
	unDiagnosedCount := 0
	ignoredCount := 0
	for _, diagnosis := range diagnoses {
		if diagnosis.Ignored {
			ignoredMessages = append(ignoredMessages, fmt.Sprintf("[info] %s (ignored):", diagnosis.Name))
			ignoredCount += 1
			continue
		}

		if !diagnosis.Diagnosed {
			warnMessages = append(warnMessages, fmt.Sprintf("[warn] %s (unknown):", diagnosis.Name))
			unDiagnosedCount += 1
			continue
		}
		if diagnosis.Archived {
			errMessages = append(errMessages, fmt.Sprintf("[error] %s (archived): %s", diagnosis.Name, diagnosis.Url))
			errCount += 1
		}
		if !diagnosis.IsActive {
			errMessages = append(errMessages, fmt.Sprintf("[error] %s (not-maintained): %s", diagnosis.Name, diagnosis.Url))
			errCount += 1
		}
	}

	fmt.Printf("\n")
	if len(ignoredMessages) > 0 {
		fmt.Println(strings.Join(ignoredMessages, "\n"))
	}
	if len(warnMessages) > 0 {
		color.Yellow(strings.Join(warnMessages, "\n"))
	}
	if len(errMessages) > 0 {
		color.Red(strings.Join(errMessages, "\n"))
	}

	color.Green(heredoc.Docf(`
		Diagnose complete! %d dependencies.
		%d error, %d unknown, %d ignored`,
		len(diagnoses),
		errCount,
		unDiagnosedCount,
		ignoredCount),
	)

	if len(errMessages) > 0 {
		return errors.New("has error")
	}

	return nil
}
