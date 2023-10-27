package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/aquasecurity/go-dep-parser/pkg/io"
	parser_io "github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/fatih/color"
	"github.com/kyoshidajp/dep-doctor/cmd/github"
	"github.com/spf13/cobra"
)

const MAX_YEAR_TO_BE_BLANK = 5

type Doctor interface {
	Diagnose(r io.ReadSeekerAt, year int) map[string]Diagnosis
	fetchURLFromRepository(name string) (string, error)
	NameWithOwners(r parser_io.ReadSeekerAt) []github.NameWithOwner
}

type Diagnosis struct {
	Name      string
	Url       string
	Archived  bool
	Diagnosed bool
	IsActive  bool
}

type Department struct {
	doctor Doctor
}

func NewDepartment(d Doctor) *Department {
	return &Department{
		doctor: d,
	}
}

func (d *Department) Diagnose(r io.ReadSeekCloserAt, year int) map[string]Diagnosis {
	return d.doctor.Diagnose(r, year)
}

type Options struct {
	packageManager string
	lockFilePath   string
}

var (
	o = &Options{}
)

var doctors = map[string]Doctor{
	"bundler": NewBundlerDoctor(),
	"yarn":    NewYarnDoctor(),
	"pip":     NewPipDoctor(),
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

		department := NewDepartment(doctor)
		diagnoses := department.Diagnose(f, MAX_YEAR_TO_BE_BLANK)
		if err := Report(diagnoses); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
	diagnoseCmd.Flags().StringVarP(&o.packageManager, "package", "p", "bundler", "package manager")
	diagnoseCmd.Flags().StringVarP(&o.lockFilePath, "lock_file", "f", "Gemfile.lock", "lock file path")
}

func Report(diagnoses map[string]Diagnosis) error {
	errMessages := []string{}
	warnMessages := []string{}
	errCount := 0
	unDiagnosedCount := 0
	for _, diagnosis := range diagnoses {
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
	if len(warnMessages) > 0 {
		color.Yellow(strings.Join(warnMessages, "\n"))
	}
	if len(errMessages) > 0 {
		color.Red(strings.Join(errMessages, "\n"))
	}

	color.Green(heredoc.Docf(`
		Diagnose complete! %d dependencies.
		%d error, %d unknown`,
		len(diagnoses),
		errCount,
		unDiagnosedCount),
	)

	if len(errMessages) > 0 {
		return errors.New("has error")
	}

	return nil
}
