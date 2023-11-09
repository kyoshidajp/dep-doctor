package dart

import (
	"net/http"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/jarcoal/httpmock"
)

func TestPub_fetchURLFromRegistry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://pub.dev/api/packages/uuid",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"latest": {
				"pubspec": {
					"repository": "https://github.com/Daegalus/dart-uuid"
				}
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pub.dev/api/packages/have-homepage",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"latest": {
				"pubspec": {
					"homepage": "https://github.com/Daegalus/have-homepage"
				}
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pub.dev/api/packages/have-issue-tracker",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"latest": {
				"pubspec": {
					"issue_tracker": "https://github.com/Daegalus/have-issue-tracker"
				}
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pub.dev/api/packages/not-found",
		httpmock.NewStringResponder(404, heredoc.Doc(`
		{}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://pub.dev/api/packages/unmarshal-error",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"unmarshal": "xxx",
		}
		`)),
	)

	tests := []struct {
		name    string
		lib     types.Library
		wantURL string
		wantErr bool
	}{
		{
			name: "repository exists",
			lib: types.Library{
				Name: "uuid",
			},
			wantURL: "https://github.com/Daegalus/dart-uuid",
			wantErr: false,
		},
		{
			name: "no repository URL, but homepage URL exists",
			lib: types.Library{
				Name: "have-homepage",
			},
			wantURL: "https://github.com/Daegalus/have-homepage",
			wantErr: false,
		},
		{
			name: "no repository and homepage URL, but issue_tracker exists",
			lib: types.Library{
				Name: "have-issue-tracker",
			},
			wantURL: "https://github.com/Daegalus/have-issue-tracker",
			wantErr: false,
		},
		{
			name: "404 not found",
			lib: types.Library{
				Name: "not-found",
			},
			wantURL: "",
			wantErr: true,
		},
		{
			// have no mock
			name: "request error",
			lib: types.Library{
				Name: "request-error",
			},
			wantURL: "",
			wantErr: true,
		},
		{
			name: "Unmarshal error",
			lib: types.Library{
				Name: "unmarshal-error",
			},
			wantURL: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pub{lib: tt.lib}
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
