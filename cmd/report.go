package cmd

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/fatih/color"
)

type Reporter interface {
	Report() error
}

type StdoutReporter struct {
	diagnoses   map[string]Diagnosis
	strict_mode bool
}

func NewStdoutReporter(diagnosis map[string]Diagnosis, strict_mode bool) *StdoutReporter {
	return &StdoutReporter{
		diagnoses:   diagnosis,
		strict_mode: strict_mode,
	}
}

func (r *StdoutReporter) Report() error {
	errMessages, warnMessages, ignoredMessages := []string{}, []string{}, []string{}
	errCount, warnCount, infoCount := 0, 0, 0
	unDiagnosedCount, ignoredCount := 0, 0

	lib_names := make([]string, 0, len(r.diagnoses))
	for key := range r.diagnoses {
		lib_names = append(lib_names, key)
	}
	sort.Strings(lib_names)

	for _, lib_name := range lib_names {
		diagnosis := r.diagnoses[lib_name]
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
	if ignoredCount > 0 {
		fmt.Println(strings.Join(ignoredMessages, "\n"))
	}
	if warnCount > 0 {
		color.Yellow(strings.Join(warnMessages, "\n"))
	}
	if errCount > 0 {
		color.Red(strings.Join(errMessages, "\n"))
	}

	color.Green(heredoc.Docf(`
		Diagnosis completed! %d libraries.
		%d error, %d warn (%d unknown), %d info (%d ignored)`,
		len(r.diagnoses),
		errCount,
		warnCount, unDiagnosedCount,
		infoCount, ignoredCount),
	)

	if errCount > 0 || r.strict_mode && warnCount > 0 {
		return errors.New("has error")
	}

	return nil
}
