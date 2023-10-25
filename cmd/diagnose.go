package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/aquasecurity/go-dep-parser/pkg/io"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const MAX_YEAR_TO_BE_BLANK = 5

type DiagnoseStrategy interface {
	Diagnose(r io.ReadSeekerAt, year int) map[string]Diagnosis
}

type Diagnosis struct {
	Name      string
	Url       string
	Archived  bool
	Diagnosed bool
	IsActive  bool
}

type Doctor struct {
	strategy DiagnoseStrategy
}

func NewDoctor(strategy DiagnoseStrategy) *Doctor {
	return &Doctor{
		strategy: strategy,
	}
}

func (d *Doctor) Diagnose(r io.ReadSeekCloserAt, year int) map[string]Diagnosis {
	return d.strategy.Diagnose(r, year)
}

var DoctorWithStrategy = map[string]DiagnoseStrategy{
	"bundler": NewBundlerStrategy(),
	"yarn":    NewYarnStrategy(),
}

type Options struct {
	packageManagerName string
	lockFilePath       string
}

var (
	o = &Options{}
)

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Diagnose packages",
	Run: func(cmd *cobra.Command, args []string) {
		lockFilePath := o.lockFilePath
		f, _ := os.Open(lockFilePath)
		defer f.Close()

		packageManagerStrategy, ok := DoctorWithStrategy[o.packageManagerName]
		if !ok {
			log.Fatal("unknown package manager")
		}

		doctor := NewDoctor(packageManagerStrategy)
		diagnoses := doctor.Diagnose(f, MAX_YEAR_TO_BE_BLANK)
		err := Report(diagnoses)
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
	diagnoseCmd.Flags().StringVarP(&o.packageManagerName, "package", "p", "bundler", "package manager")
	diagnoseCmd.Flags().StringVarP(&o.lockFilePath, "lock_file", "f", "Gemfile.lock", "lock file path")
}

func Report(diagnoses map[string]Diagnosis) error {
	errMessages := []string{}
	warnMessages := []string{}
	errorCount := 0
	unDiagnosedCount := 0
	for _, diagnosis := range diagnoses {
		if !diagnosis.Diagnosed {
			warnMessages = append(warnMessages, fmt.Sprintf("[warn] %s (unknown):", diagnosis.Name))
			unDiagnosedCount += 1
			continue
		}
		if diagnosis.Archived {
			errMessages = append(errMessages, fmt.Sprintf("[error] %s (archived): %s", diagnosis.Name, diagnosis.Url))
			errorCount += 1
		}
		if !diagnosis.IsActive {
			errMessages = append(errMessages, fmt.Sprintf("[error] %s (not-maintained): %s", diagnosis.Name, diagnosis.Url))
			errorCount += 1
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
		errorCount,
		unDiagnosedCount),
	)

	if len(errMessages) > 0 {
		return errors.New("has error")
	}

	return nil
}
