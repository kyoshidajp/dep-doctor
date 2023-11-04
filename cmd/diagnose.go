package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/aquasecurity/go-dep-parser/pkg/io"
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

func Diagnose(d Doctor, r io.ReadSeekCloserAt, year int, ignores []string, cache map[string]string, disableCache bool) map[string]Diagnosis {
	diagnoses := make(map[string]Diagnosis)
	slicedParams := [][]github.FetchRepositoryParam{}
	libs := d.Libraries(r)
	fetchRepositoryParams := FetchRepositoryParams(libs, d, cache, disableCache)
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

func Report(diagnoses map[string]Diagnosis, strict_mode bool) error {
	reporter := NewStdoutReporter(diagnoses, strict_mode)
	return reporter.Report()
}

type Options struct {
	packageManager string
	filePath       string
	ignores        string
	year           int
	strict         bool
	disableCache   bool
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
		if err := Prepare(); err != nil {
			log.Fatal(err)
		}

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

		cacheStore := BuildCacheStore()
		cache := cacheStore.URLbyPackageManager(o.packageManager)
		diagnoses := Diagnose(doctor, f, o.year, o.Ignores(), cache, o.disableCache)
		if err := SaveCache(diagnoses, cacheStore, o.packageManager); err != nil {
			log.Fatal(err)
		}

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
	diagnoseCmd.PersistentFlags().BoolVarP(&o.disableCache, "disable-cache", "", false, "without using cache")

	if err := diagnoseCmd.MarkFlagRequired("package"); err != nil {
		fmt.Println(err.Error())
	}
	if err := diagnoseCmd.MarkFlagRequired("file"); err != nil {
		fmt.Println(err.Error())
	}
}
