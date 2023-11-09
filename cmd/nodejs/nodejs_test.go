package nodejs

import (
	"net/http"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/jarcoal/httpmock"
)

func TestNodejs_fetchURLFromRegistry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://registry.npmjs.org/react",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"repository": {
				"url": "git+https://github.com/facebook/react.git"
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://registry.npmjs.org/not-found",
		httpmock.NewStringResponder(404, heredoc.Doc(`
		{}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://registry.npmjs.org/unmarshal-error",
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
			name: "source_code_uri exists",
			lib: types.Library{
				Name: "react",
			},
			wantURL: "git+https://github.com/facebook/react.git",
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
			n := Nodejs{lib: tt.lib}
			got, _ := n.fetchURLFromRegistry(http.Client{})
			expect := tt.wantURL
			if got != expect {
				t.Errorf("get() = %v, want %v", got, expect)
			}
		})
	}
}
