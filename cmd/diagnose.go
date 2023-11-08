package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strings"
	"sync"

	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
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

func (params RepositoryParams) SearchableParams() []github.FetchRepositoryParam {
	uniqParams := map[string]github.FetchRepositoryParam{}
	for _, param := range params {
		uniqKey := param.RepoOwner()
		uniqParams[uniqKey] = param
	}

	searchableParams := []github.FetchRepositoryParam{}
	for _, param := range uniqParams {
		if param.Searchable {
			searchableParams = append(searchableParams, param)
		}
	}
	sort.SliceStable(searchableParams, func(i, j int) bool { return searchableParams[i].PackageName < searchableParams[j].PackageName })

	return searchableParams
}

func (params RepositoryParams) SlicedParams() [][]github.FetchRepositoryParam {
	slicedParams := [][]github.FetchRepositoryParam{}
	searchableParams := params.SearchableParams()
	sliceSize := len(searchableParams)

	for i := 0; i < sliceSize; i += github.SEARCH_REPOS_PER_ONCE {
		end := i + github.SEARCH_REPOS_PER_ONCE
		if sliceSize < end {
			end = sliceSize
		}
		slicedParams = append(slicedParams, searchableParams[i:end])
	}

	return slicedParams
}

func (params RepositoryParams) diagnoses(searchedRepos []github.GitHubRepository, ignores []string, year int) map[string]Diagnosis {
	repoByName := make(map[string]github.GitHubRepository)
	for _, r := range searchedRepos {
		uniqKey := r.RepoOwner()
		repoByName[uniqKey] = r
	}

	diagnoses := make(map[string]Diagnosis)
	for _, param := range params {
		uniqKey := param.RepoOwner()
		diagnosis := Diagnosis{}
		repo, ok := repoByName[uniqKey]
		if ok {
			willIgnore := slices.Contains(ignores, repo.Name)
			diagnosis = Diagnosis{
				Name:      param.PackageName,
				URL:       repo.URL,
				Archived:  repo.Archived,
				Ignored:   willIgnore,
				Diagnosed: true,
				IsActive:  repo.IsActive(year),
				Error:     repo.Error,
			}
			diagnoses[repo.Name] = diagnosis
		} else {
			diagnosis = Diagnosis{
				Name:      param.PackageName,
				Diagnosed: false,
				Error:     param.Error,
			}
		}
		diagnoses[param.PackageName] = diagnosis
	}
	return diagnoses
}

type Libraries []types.Library

func (libs Libraries) Uniq() Libraries {
	nameWithLib := map[string]types.Library{}
	for _, lib := range libs {
		nameWithLib[lib.Name] = lib
	}

	var uniqLibs Libraries
	for _, lib := range nameWithLib {
		uniqLibs = append(uniqLibs, lib)
	}
	sort.SliceStable(uniqLibs, func(i, j int) bool { return uniqLibs[i].Name < uniqLibs[j].Name })

	return uniqLibs
}

func NewLibraries(libs []types.Library) Libraries {
	var newLibs Libraries
	for _, lib := range libs {
		newLibs = append(newLibs, lib)
	}
	sort.SliceStable(newLibs, func(i, j int) bool { return newLibs[i].Name < newLibs[j].Name })

	return newLibs
}

func Prepare() error {
	token := os.Getenv(github.TOKEN_NAME)
	if len(token) == 0 {
		m := fmt.Sprintf("The Environment variable `%s` is not set. It must be set before execution. For example, please refer to https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens", github.TOKEN_NAME)
		return errors.New(m)
	}
	return nil
}

func FetchRepositoryParams(libs []types.Library, d Doctor, cache map[string]string, disableCache bool) RepositoryParams {
	var params []github.FetchRepositoryParam
	var wg sync.WaitGroup
	sem := make(chan struct{}, FETCH_REPOS_PER_ONCE)

	for _, lib := range libs {
		wg.Add(1)
		sem <- struct{}{}
		go func(lib types.Library) {
			defer wg.Done()
			defer func() { <-sem }()

			var url string
			url, ok := cache[lib.Name]
			if !disableCache && ok {
				fmt.Printf("%s (from source URL cache)\n", lib.Name)
			} else {
				fmt.Printf("%s\n", lib.Name)
				var err error
				url, err = d.SourceCodeURL(lib)
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

func Diagnose(d Doctor, r parser_io.ReadSeekCloserAt, year int, ignores []string, cache map[string]string, disableCache bool) map[string]Diagnosis {
	libs := NewLibraries(d.Libraries(r)).Uniq()
	searchParams := FetchRepositoryParams(libs, d, cache, disableCache)
	slicedSearchParams := searchParams.SlicedParams()
	searchedRepos := []github.GitHubRepository{}

	var wg sync.WaitGroup
	sem := make(chan struct{}, FETCH_REPOS_PER_ONCE)
	for _, params := range slicedSearchParams {
		wg.Add(1)
		sem <- struct{}{}
		go func(params []github.FetchRepositoryParam) {
			defer wg.Done()
			defer func() { <-sem }()

			tmpRepos := github.FetchFromGitHub(params)
			searchedRepos = append(searchedRepos, tmpRepos...)
		}(params)
	}
	wg.Wait()

	diagnoses := searchParams.diagnoses(searchedRepos, ignores, year)
	return diagnoses
}

func Report(diagnoses map[string]Diagnosis, strict_mode bool) error {
	reporter := NewStdoutReporter(diagnoses, strict_mode)
	return reporter.Report()
}

type DiagnoseOption struct {
	packageManager string
	filePath       string
	ignores        string
	year           int
	strict         bool
	disableCache   bool

	Out    io.Writer
	ErrOut io.Writer
}

func (o *DiagnoseOption) Ignores() []string {
	return strings.Split(o.ignores, " ")
}

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
	"poetry":    python.NewPoetryDoctor(),
	"pub":       dart.NewPubDoctor(),
	"yarn":      nodejs.NewYarnDoctor(),
}

func newDiagnoseCmd(o *DiagnoseOption) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Diagnose dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Prepare(); err != nil {
				return err
			}

			doctor, ok := doctors[o.packageManager]
			if !ok {
				m := doctors.UnknownErrorMessage(o.packageManager)
				return fmt.Errorf(m)
			}

			filePath := o.filePath
			f, err := os.Open(filePath)
			defer func() {
				_ = f.Close()
			}()
			if err != nil {
				m := fmt.Sprintf("Can't open: %s", o.filePath)
				return fmt.Errorf(m)
			}

			cacheStore := BuildCacheStore()
			cache := cacheStore.URLbyPackageManager(o.packageManager)
			diagnoses := Diagnose(doctor, f, o.year, o.Ignores(), cache, o.disableCache)
			if err := SaveCache(diagnoses, cacheStore, o.packageManager); err != nil {
				return err
			}

			if err := Report(diagnoses, o.strict); err != nil {
				return errors.New("")
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&o.packageManager, "package", "p", "", "package manager")
	cmd.Flags().StringVarP(&o.filePath, "file", "f", "", "dependencies file path")
	cmd.Flags().StringVarP(&o.ignores, "ignores", "i", "", "ignore dependencies (separated by a space)")
	cmd.Flags().IntVarP(&o.year, "year", "y", MAX_YEAR_TO_BE_BLANK, "max years of inactivity")
	cmd.PersistentFlags().BoolVarP(&o.strict, "strict", "", false, "exit with non-zero if warnings exist")
	cmd.PersistentFlags().BoolVarP(&o.disableCache, "disable-cache", "", false, "without using cache")

	cmd.SetOut(o.Out)
	cmd.SetErr(o.ErrOut)

	if err := cmd.MarkFlagRequired("package"); err != nil {
		fmt.Println(err.Error())
	}
	if err := cmd.MarkFlagRequired("file"); err != nil {
		fmt.Println(err.Error())
	}

	return cmd
}
