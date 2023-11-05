package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/kyoshidajp/dep-doctor/cmd/ruby"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	expect := map[string]Diagnosis{
		"faker": {
			Name:      "faker",
			URL:       "https://github.com/faker-ruby/faker",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"concurrent-ruby": {
			Name:      "concurrent-ruby",
			URL:       "https://github.com/ruby-concurrency/concurrent-ruby",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"i18n": {
			Name:      "i18n",
			URL:       "https://github.com/ruby-i18n/i18n",
			Archived:  false,
			Ignored:   true,
			Diagnosed: true,
			IsActive:  true,
		},
		"method_source": {
			Name:      "method_source",
			URL:       "https://github.com/banister/method_source",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"dotenv": {
			Name:      "dotenv",
			URL:       "https://github.com/bkeepers/dotenv",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
	}

	t.Run("test", func(t *testing.T) {
		f, err := os.Open("ruby/bundler/testdata/Gemfile.lock")
		require.NoError(t, err)
		defer f.Close()

		doctor := ruby.NewBundlerDoctor()
		ignores := []string{"i18n"}
		cache := map[string]string{}
		diagnoses := Diagnose(doctor, f, 2, ignores, cache, false)
		assert.Equal(t, expect, diagnoses)
	})
}

func TestDiagnosis_ErrorMessage(t *testing.T) {
	tests := []struct {
		name      string
		diagnosis Diagnosis
	}{
		{
			name: "has Error",
			diagnosis: Diagnosis{
				Error: errors.New("unknown error"),
			},
		},
		{
			name: "has no Error",
			diagnosis: Diagnosis{
				Error: nil,
			},
		},
	}
	expects := []struct {
		name    string
		message string
	}{
		{
			name:    "has Error",
			message: "unknown error",
		},
		{
			name:    "has no Error",
			message: "",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.diagnosis.ErrorMessage()
			expect := expects[i].message
			assert.Equal(t, expect, actual)
		})
	}
}

func TestOptions_Ignores(t *testing.T) {
	tests := []struct {
		name   string
		option DiagnoseOption
	}{
		{
			name: "ignores",
			option: DiagnoseOption{
				ignores: "package1 package2 package3",
			},
		},
	}
	expects := []struct {
		name    string
		ignores []string
	}{
		{
			name:    "ignores",
			ignores: []string{"package1", "package2", "package3"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.option.Ignores()
			expect := expects[i].ignores
			assert.Equal(t, expect, actual)
		})
	}
}

func TestDoctors_PackageManagers(t *testing.T) {
	tests := []struct {
		name    string
		doctors Doctors
	}{
		{
			name: "doctors",
			doctors: Doctors{
				"package1": nil,
				"package2": nil,
				"package3": nil,
			},
		},
	}
	expects := []struct {
		name     string
		packages []string
	}{
		{
			name:     "doctors",
			packages: []string{"package1", "package2", "package3"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.doctors.PackageManagers()
			expect := expects[i].packages
			assert.Equal(t, expect, actual)
		})
	}
}

func TestDoctors_UnknownErrorMessage(t *testing.T) {
	tests := []struct {
		name    string
		doctors Doctors
	}{
		{
			name: "doctors",
			doctors: Doctors{
				"p1": nil,
				"p2": nil,
				"p3": nil,
			},
		},
	}
	expects := []struct {
		name    string
		message string
	}{
		{
			name:    "doctors",
			message: "Unknown package manager: xxx. You can choose from [p1, p2, p3]",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.doctors.UnknownErrorMessage("xxx")
			expect := expects[i].message
			assert.Equal(t, expect, actual)
		})
	}
}

func TestDiagnose_newDiagnoseCmd(t *testing.T) {
	tests := []struct {
		name          string
		command       string
		wantOutWriter string
		wantErrWriter string
		wantErr       bool
	}{
		{
			name:          "bundler with no problems",
			command:       "--package bundler --file ruby/bundler/testdata/Gemfile.lock",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "yarn",
			command:       "--package yarn --file nodejs/yarn/testdata/yarn.lock",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "npm",
			command:       "--package npm --file nodejs/npm/testdata/package-lock.json",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "pip",
			command:       "--package pip --file python/pip/testdata/requirements.txt",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "pipenv",
			command:       "--package pipenv --file python/pipenv/testdata/Pipfile.lock",
			wantOutWriter: "",
			wantErrWriter: "has error",
			wantErr:       true,
		},
		{
			name:          "golang",
			command:       "--package golang --file golang/mod/testdata/go.mod",
			wantOutWriter: "",
			wantErrWriter: "has error",
			wantErr:       true,
		},
		{
			name:          "cargo",
			command:       "--package cargo --file rust/cargo/testdata/cargo.lock",
			wantOutWriter: "",
			wantErrWriter: "has error",
			wantErr:       true,
		},
		{
			name:          "cocoapods",
			command:       "--package cocoapods --file swift/cocoapods/testdata/Podfile.lock",
			wantOutWriter: "",
			wantErrWriter: "has error",
			wantErr:       true,
		},
		{
			name:          "pub",
			command:       "--package pub --file dart/pub/testdata/podspec.lock",
			wantOutWriter: "",
			wantErrWriter: "has error",
			wantErr:       true,
		},
		{
			name:          "mix",
			command:       "--package mix --file erlang_elixir/hex/testdata/mix.lock",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
		{
			name:          "has error",
			command:       "--package bundler --file ruby/bundler/testdata/Gemfile_error.lock",
			wantOutWriter: "",
			wantErrWriter: "has error",
			wantErr:       true,
		},
		{
			name:          "unknown package manager",
			command:       "--package unknown --file ruby/bundler/testdata/Gemfile.lock",
			wantOutWriter: "",
			wantErrWriter: "Unknown package manager: unknown. You can choose from [bundler, cargo, cocoapods, composer, golang, mix, npm, pip, pipenv, pub, yarn]",
			wantErr:       true,
		},
		{
			name:          "can't open file",
			command:       "--package bundler --file cant_open_file.txt",
			wantOutWriter: "",
			wantErrWriter: "Can't open: cant_open_file.txt",
			wantErr:       true,
		},
		{
			name:          "no package option",
			command:       "--file ruby/bundler/testdata/Gemfile.lock",
			wantOutWriter: "",
			wantErrWriter: "required flag(s) \"package\" not set",
			wantErr:       true,
		},
		{
			name:          "no file option",
			command:       "--package bundler",
			wantOutWriter: "",
			wantErrWriter: "required flag(s) \"file\" not set",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outWriter := &bytes.Buffer{}
			errWriter := &bytes.Buffer{}

			o := &DiagnoseOption{
				Out:    outWriter,
				ErrOut: errWriter,
			}

			cmd := newDiagnoseCmd(o)
			if tt.command != "" {
				cmd.SetArgs(strings.Split(tt.command, " "))
			}

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if err.Error() != tt.wantErrWriter {
					t.Errorf("get() = %v, want %v", err.Error(), tt.wantErrWriter)
				}
			} else {
				if gotOutWriter := outWriter.String(); gotOutWriter != tt.wantOutWriter {
					t.Errorf("get() = %v, want %v", gotOutWriter, tt.wantOutWriter)
				}
			}
		})
	}
}
