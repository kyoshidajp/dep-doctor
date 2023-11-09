package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jarcoal/httpmock"
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
			Active:    true,
		},
		"concurrent-ruby": {
			Name:      "concurrent-ruby",
			URL:       "https://github.com/ruby-concurrency/concurrent-ruby",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			Active:    true,
		},
		"i18n": {
			Name:      "i18n",
			URL:       "https://github.com/ruby-i18n/i18n",
			Archived:  false,
			Ignored:   true,
			Diagnosed: true,
			Active:    true,
		},
		"method_source": {
			Name:      "method_source",
			URL:       "https://github.com/banister/method_source",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			Active:    true,
		},
		"dotenv": {
			Name:      "dotenv",
			URL:       "https://github.com/bkeepers/dotenv",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			Active:    true,
		},
	}

	t.Run("test", func(t *testing.T) {
		f, err := os.Open("ruby/bundler/testdata/Gemfile.lock")
		require.NoError(t, err)
		defer f.Close()

		doctor := ruby.NewBundlerDoctor()
		cache := map[string]string{}
		o := DiagnoseOption{
			ignores: "i18n",
			year:    5,
		}
		diagnoses := Diagnose(doctor, f, cache, o)
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
	if testing.Short() {
		t.SkipNow()
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// rubygems
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/concurrent-ruby.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
			{
				"name": "concurrent-ruby",
				"source_code_uri": "https://github.com/concurrent-ruby/concurrent-ruby"
			}
			`)),
	)
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/dotenv.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
			{
				"name": "minitest",
				"homepage_uri": "https://github.com/bkeepers/dotenv"
			}
			`)),
	)
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/faker.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
			{
				"name": "faker",
				"source_code_uri": "https://github.com/faker-ruby/faker"
			}
			`)),
	)
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/i18n.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
			{
				"name": "i18n",
				"source_code_uri": "https://github.com/ruby-i18n/i18n"
			}
			`)),
	)
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/method_source.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
			{
				"name": "i18n",
				"source_code_uri": "https://github.com/banister/method_source"
			}
			`)),
	)
	httpmock.RegisterResponder("POST", "https://api.github.com/graphql",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"data": {
			  "search": {
				"repositoryCount": 4,
				"nodes": [
				  {
					"isArchived": false,
					"url": "https://github.com/faker-ruby/faker",
					"name": "faker",
					"defaultBranchRef": {
					  "target": {
						"history": {
						  "edges": [
							{
							  "node": {
								"committedDate": "2023-11-03T21:10:51Z"
							  }
							}
						  ]
						}
					  }
					}
				  },
				  {
					"isArchived": false,
					"url": "https://github.com/bkeepers/dotenv",
					"name": "dotenv",
					"defaultBranchRef": {
					  "target": {
						"history": {
						  "edges": [
							{
							  "node": {
								"committedDate": "2022-07-27T14:37:34Z"
							  }
							}
						  ]
						}
					  }
					}
				  },
				  {
					"isArchived": false,
					"url": "https://github.com/ruby-i18n/i18n",
					"name": "i18n",
					"defaultBranchRef": {
					  "target": {
						"history": {
						  "edges": [
							{
							  "node": {
								"committedDate": "2023-06-21T10:33:08Z"
							  }
							}
						  ]
						}
					  }
					}
				  }
				]
			  }
			}
		  }
			`)),
	)

	// yarn
	httpmock.RegisterResponder("GET", "https://registry.npmjs.org/asap",
		httpmock.NewStringResponder(200, heredoc.Doc(`
			{
				"repository": {
					"url": "https://github.com/facebook/asap"
				}
			}
			`)),
	)
	httpmock.RegisterResponder("GET", "https://registry.npmjs.org/jquery",
		httpmock.NewStringResponder(200, heredoc.Doc(`
			{
				"repository": {
					"url": "https://github.com/jquery/jquery"
				}
			}
			`)),
	)
	httpmock.RegisterResponder("GET", "https://registry.npmjs.org/promise",
		httpmock.NewStringResponder(200, heredoc.Doc(`
			{
				"repository": {
					"url": "https://github.com/then/promise"
				}
			}
			`)),
	)

	tests := []struct {
		name          string
		command       string
		wantOutWriter string
		wantErrWriter string
		wantErr       bool
	}{
		/*
				{
					name:          "bundler with no problems",
					command:       "--disable-cache --package bundler --file ruby/bundler/testdata/Gemfile.lock",
					wantOutWriter: "",
					wantErrWriter: "",
					wantErr:       false,
				},
			{
				name:          "yarn",
				command:       "--disable-cache --package yarn --file nodejs/yarn/testdata/yarn.lock",
				wantOutWriter: "",
				wantErrWriter: "",
				wantErr:       false,
			},
		*/
		{
			name:          "npm",
			command:       "--disable-cache --package npm --file nodejs/npm/testdata/package-lock.json",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       false,
		},
<<<<<<< Updated upstream
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
			wantErrWriter: "",
			wantErr:       true,
		},
		{
			name:          "golang",
			command:       "--package golang --file golang/mod/testdata/go.mod",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       true,
		},
		{
			name:          "cargo",
			command:       "--package cargo --file rust/cargo/testdata/cargo.lock",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       true,
		},
		{
			name:          "cocoapods",
			command:       "--package cocoapods --file swift/cocoapods/testdata/Podfile.lock",
			wantOutWriter: "",
			wantErrWriter: "",
			wantErr:       true,
		},
		{
			name:          "pub",
			command:       "--package pub --file dart/pub/testdata/podspec.lock",
			wantOutWriter: "",
			wantErrWriter: "",
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
			wantErrWriter: "",
			wantErr:       true,
		},
		{
			name:          "unknown package manager",
			command:       "--package unknown --file ruby/bundler/testdata/Gemfile.lock",
			wantOutWriter: "",
			wantErrWriter: "Unknown package manager: unknown. You can choose from [bundler, cargo, cocoapods, composer, golang, mix, npm, pip, pipenv, poetry, pub, yarn]",
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
=======
		/*
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
				wantErrWriter: "",
				wantErr:       true,
			},
			{
				name:          "golang",
				command:       "--package golang --file golang/mod/testdata/go.mod",
				wantOutWriter: "",
				wantErrWriter: "",
				wantErr:       true,
			},
			{
				name:          "cargo",
				command:       "--package cargo --file rust/cargo/testdata/cargo.lock",
				wantOutWriter: "",
				wantErrWriter: "",
				wantErr:       true,
			},
			{
				name:          "cocoapods",
				command:       "--package cocoapods --file swift/cocoapods/testdata/Podfile.lock",
				wantOutWriter: "",
				wantErrWriter: "",
				wantErr:       true,
			},
			{
				name:          "pub",
				command:       "--package pub --file dart/pub/testdata/podspec.lock",
				wantOutWriter: "",
				wantErrWriter: "",
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
				wantErrWriter: "",
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
		*/
>>>>>>> Stashed changes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outWriter := &bytes.Buffer{}
			errWriter := &bytes.Buffer{}

			o := &DiagnoseOption{
				Out:    outWriter,
				ErrOut: errWriter,
			}

			cmd := newDiagnoseCmd(*o)
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
