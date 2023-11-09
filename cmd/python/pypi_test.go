package python

import (
	"net/http"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jarcoal/httpmock"
)

func TestPyPi_fetchURLFromRegistry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://pypi.org/pypi/code/json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"info": {
				"project_urls": {
					"Code": "https://github.com/pypa/code"
				}
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pypi.org/pypi/github_project/json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"info": {
				"project_urls": {
					"GitHub Project": "https://github.com/pypa/github_project"
				}
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pypi.org/pypi/pip/json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"info": {
				"project_urls": {
					"Source Code": "https://github.com/pypa/pip"
				}
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pypi.org/pypi/source/json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"info": {
				"project_urls": {
					"Source": "https://github.com/source/source"
				}
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pypi.org/pypi/not-found/json",
		httpmock.NewStringResponder(404, heredoc.Doc(`
		{}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pypi.org/pypi/unmarshal-error/json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"unmarshal": "xxx",
		}
		`)),
	)

	tests := []struct {
		name    string
		libName string
		wantURL string
		wantErr bool
	}{
		{
			name:    "info.project_urls.GitHub Project exists",
			libName: "github_project",
			wantURL: "https://github.com/pypa/github_project",
			wantErr: false,
		},
		{
			name:    "info.project_urls.Code exists",
			libName: "code",
			wantURL: "https://github.com/pypa/code",
			wantErr: false,
		},
		{
			name:    "info.project_urls.Source_Code exists",
			libName: "pip",
			wantURL: "https://github.com/pypa/pip",
			wantErr: false,
		},
		{
			name:    "info.project_urls.Source exists",
			libName: "source",
			wantURL: "https://github.com/source/source",
			wantErr: false,
		},
		{
			name:    "not found",
			libName: "not-found-xxxx",
			wantURL: "",
			wantErr: true,
		},
		{
			name:    "404 not found",
			libName: "not-found",
			wantURL: "",
			wantErr: true,
		},
		{
			name:    "Unmarshal error",
			libName: "unmarshal-error",
			wantURL: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pypi{name: tt.libName}
			got, err := p.fetchURLFromRegistry(http.Client{})
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			expect := tt.wantURL
			if got != expect {
				t.Errorf("get() = %v, want starts %v", got, expect)
			}
		})
	}
}
