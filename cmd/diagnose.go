package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/MakeNowJust/heredoc"
	"github.com/aquasecurity/go-dep-parser/pkg/io"
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/fatih/color"
	dart "github.com/kyoshidajp/dep-doctor/cmd/dart/pub"
	erlang_elixir "github.com/kyoshidajp/dep-doctor/cmd/erlang_elixir/hex"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
	"github.com/kyoshidajp/dep-doctor/cmd/golang"
	"github.com/kyoshidajp/dep-doctor/cmd/nodejs"
	"github.com/kyoshidajp/dep-doctor/cmd/php"
	"github.com/kyoshidajp/dep-doctor/cmd/python"
	"github.com/kyoshidajp/dep-doctor/cmd/ruby"
	"github.com/kyoshidajp/dep-doctor/cmd/rust"
	"github.com/kyoshidajp/dep-doctor/cmd/swift"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

const MAX_YEAR_TO_BE_BLANK = 5

// referenced as the number of goroutine parallels
// should be optimized?
const FETCH_REPOS_PER_ONCE = 20

type Diagnosis struct {
	Name      string
	URL       string
	Archived  bool
	Ignored   bool
	Diagnosed bool
	IsActive  bool
	Error     error
}

func (d *Diagnosis) ErrorMessage() string {
	if d.Error == nil {
		return ""
	}
	return fmt.Sprintf("%s", d.Error)
}

type Doctor interface {
	Libraries(r parser_io.ReadSeekerAt) []types.Library
	SourceCodeURL(lib types.Library) (string, error)
}

type RepositoryParams []github.FetchRepositoryParam

func (p RepositoryParams) SearchableParams() []github.FetchRepositoryParam {
	params := []github.FetchRepositoryParam{}
	for _, param := range p {
		if param.Searchable {
			params = append(params, param)
		}
	}
	return params
}

func FetchRepositoryParams(libs []types.Library, d Doctor) RepositoryParams {
	var params []github.FetchRepositoryParam
	var wg sync.WaitGroup
	sem := make(chan struct{}, FETCH_REPOS_PER_ONCE)

	for _, lib := range libs {
		wg.Add(1)
		sem <- struct{}{}
		go func(lib types.Library) {
			defer wg.Done()
			defer func() { <-sem }()

			fmt.Printf("%s\n", lib.Name)

			url, err := d.SourceCodeURL(lib)
			if err != nil {
				params = append(params,
					github.FetchRepositoryParam{
						PackageName: lib.Name,
						Searchable:  false,
						Error:       err,
					},
				)
				return
			}

			repo, err := github.ParseGitHubURL(url)
			if err != nil {
				params = append(params,
					github.FetchRepositoryParam{
						PackageName: lib.Name,
						Searchable:  false,
						Error:       err,
					},
				)
				return
			}

			params = append(params,
				github.FetchRepositoryParam{
					Repo:        repo.Repo,
					Owner:       repo.Owner,
					PackageName: lib.Name,
					Searchable:  true,
				},
			)
		}(lib)

		wg.Wait()
	}

	return params
}

func Diagnose(d Doctor, r io.ReadSeekCloserAt, year int, ignores []string) map[string]Diagnosis {
	diagnoses := make(map[string]Diagnosis)
	slicedParams := [][]github.FetchRepositoryParam{}
	libs := d.Libraries(r)
	fetchRepositoryParams := FetchRepositoryParams(libs, d)
	searchableRepositoryParams := fetchRepositoryParams.SearchableParams()
	sliceSize := len(searchableRepositoryParams)

	for i := 0; i < sliceSize; i += github.SEARCH_REPOS_PER_ONCE {
		end := i + github.SEARCH_REPOS_PER_ONCE
		if sliceSize < end {
			end = sliceSize
		}
		slicedParams = append(slicedParams, searchableRepositoryParams[i:end])
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, FETCH_REPOS_PER_ONCE)
	for _, params := range slicedParams {
		wg.Add(1)
		sem <- struct{}{}
		go func(params []github.FetchRepositoryParam) {
			defer wg.Done()
			defer func() { <-sem }()

			repos := github.FetchFromGitHub(params)
			for _, r := range repos {
				isIgnore := slices.Contains(ignores, r.Name)
				diagnosis := Diagnosis{
					Name:      r.Name,
					URL:       r.URL,
					Archived:  r.Archived,
					Ignored:   isIgnore,
					Diagnosed: true,
					IsActive:  r.IsActive(year),
					Error:     r.Error,
				}
				diagnoses[r.Name] = diagnosis
			}
		}(params)
	}

	wg.Wait()

	for _, fetchRepositoryParam := range fetchRepositoryParams {
		if fetchRepositoryParam.Searchable {
			continue
		}

		diagnosis := Diagnosis{
			Name:      fetchRepositoryParam.PackageName,
			Diagnosed: false,
			Error:     fetchRepositoryParam.Error,
		}
		diagnoses[fetchRepositoryParam.PackageName] = diagnosis
	}
	return diagnoses
}

type Options struct {
	packageManager string
	filePath       string
	ignores        string
	year           int
	strict         bool
}

func (o *Options) Ignores() []string {
	return strings.Split(o.ignores, " ")
}

var (
	o = &Options{}
)

type Doctors map[string]Doctor

func (d Doctors) PackageManagers() []string {
	packages := []string{}
	for p := range d {
		packages = append(packages, p)
	}
	sort.Strings(packages)
	return packages
}

func (d Doctors) UnknownErrorMessage(packageManager string) string {
	return fmt.Sprintf("Unknown package manager: %s. You can choose from [%s]",
		packageManager,
		strings.Join(d.PackageManagers(), ", "))
}

var doctors = Doctors{
	"bundler":   ruby.NewBundlerDoctor(),
	"cargo":     rust.NewCargoDoctor(),
	"cocoapods": swift.NewCococaPodsDoctor(),
	"composer":  php.NewComposerDoctor(),
	"golang":    golang.NewGolangDoctor(),
	"npm":       nodejs.NewNPMDoctor(),
	"mix":       erlang_elixir.NewMixDoctor(),
	"pip":       python.NewPipDoctor(),
	"pipenv":    python.NewPipenvDoctor(),
	"pub":       dart.NewPubDoctor(),
	"yarn":      nodejs.NewYarnDoctor(),
}

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Diagnose dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		doctor, ok := doctors[o.packageManager]
		if !ok {
			m := doctors.UnknownErrorMessage(o.packageManager)
			log.Fatal(m)
		}

		filePath := o.filePath
		f, err := os.Open(filePath)
		defer func() {
			_ = f.Close()
		}()
		if err != nil {
			m := fmt.Sprintf("Can't open: %s.", o.filePath)
			log.Fatal(m)
		}

		diagnoses := Diagnose(doctor, f, o.year, o.Ignores())
		if err := Report(diagnoses, o.strict); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
	diagnoseCmd.Flags().StringVarP(&o.packageManager, "package", "p", "", "package manager")
	diagnoseCmd.Flags().StringVarP(&o.filePath, "file", "f", "", "dependencies file path")
	diagnoseCmd.Flags().StringVarP(&o.ignores, "ignores", "i", "", "ignore dependencies (separated by a space)")
	diagnoseCmd.Flags().IntVarP(&o.year, "year", "y", MAX_YEAR_TO_BE_BLANK, "max years of inactivity")
	diagnoseCmd.PersistentFlags().BoolVarP(&o.strict, "strict", "", false, "exit with non-zero if warnings exist")

	if err := diagnoseCmd.MarkFlagRequired("package"); err != nil {
		fmt.Println(err.Error())
	}
	if err := diagnoseCmd.MarkFlagRequired("file"); err != nil {
		fmt.Println(err.Error())
	}
}

func Report(diagnoses map[string]Diagnosis, strict_mode bool) error {
	errMessages, warnMessages, ignoredMessages := []string{}, []string{}, []string{}
	errCount, warnCount, infoCount := 0, 0, 0
	unDiagnosedCount, ignoredCount := 0, 0

	lib_names := make([]string, 0, len(diagnoses))
	for key := range diagnoses {
		lib_names = append(lib_names, key)
	}
	sort.Strings(lib_names)

	for _, lib_name := range lib_names {
		diagnosis := diagnoses[lib_name]
		if diagnosis.Ignored {
			ignoredMessages = append(ignoredMessages, fmt.Sprintf("[info] %s (ignored):", diagnosis.Name))
			ignoredCount += 1
			infoCount += 1
			continue
		}

		if diagnosis.Error != nil {
			errMessages = append(errMessages, fmt.Sprintf("[error] %s: %s", diagnosis.Name, diagnosis.Error))
			errCount += 1
			continue
		}

		if !diagnosis.Diagnosed {
			warnMessages = append(warnMessages, fmt.Sprintf("[warn] %s (unknown): %s", diagnosis.Name, diagnosis.ErrorMessage()))
			unDiagnosedCount += 1
			warnCount += 1
			continue
		}
		if diagnosis.Archived {
			errMessages = append(errMessages, fmt.Sprintf("[error] %s (archived): %s", diagnosis.Name, diagnosis.URL))
			errCount += 1
		}
		if !diagnosis.IsActive {
			warnMessages = append(warnMessages, fmt.Sprintf("[warn] %s (not-maintained): %s", diagnosis.Name, diagnosis.URL))
			warnCount += 1
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
		Diagnosis completed! %d libraries.
		%d error, %d warn (%d unknown), %d info (%d ignored)`,
		len(diagnoses),
		errCount,
		warnCount, unDiagnosedCount,
		infoCount, ignoredCount),
	)

	if len(errMessages) > 0 || strict_mode && warnCount > 0 {
		return errors.New("has error")
	}

	return nil
}
