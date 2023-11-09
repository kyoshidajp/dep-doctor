package ruby

import (
	"net/http"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jarcoal/httpmock"
)

func TestRubyGems_fetchURLFromRegistry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/rails.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"name": "rails",
			"source_code_uri": "https://github.com/rails/rails/tree/v7.1.1",
			"homepage_uri": "https://rubyonrails.org"
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/minitest.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"name": "minitest",
			"source_code_uri": null,
			"homepage_uri": "https://github.com/minitest/minitest"
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/not-found.json",
		httpmock.NewStringResponder(404, heredoc.Doc(`
		{}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/no-url.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://rubygems.org/api/v1/gems/unmarshal-error.json",
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
			name:    "source_code_uri exists",
			libName: "rails",
			wantURL: "https://github.com/rails/rails",
			wantErr: false,
		},
		{
			name:    "no source_code_uri, but homepage_uri exists",
			libName: "minitest",
			wantURL: "https://github.com/minitest/minitest",
			wantErr: false,
		},
		{
			name:    "404 not found",
			libName: "not-found",
			wantURL: "",
			wantErr: true,
		},
		{
			name:    "no URL",
			libName: "no-url",
			wantURL: "",
			wantErr: true,
		},
		{
			// have no mock
			name:    "request error",
			libName: "request-error",
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
			g := RubyGems{name: tt.libName}
			got, err := g.fetchURLFromRegistry(http.Client{})
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			expect := tt.wantURL
			if !strings.HasPrefix(got, expect) {
				t.Errorf("get() = %v, want starts %v", got, expect)
			}
		})
	}
}
