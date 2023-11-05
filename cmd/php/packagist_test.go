package php

import (
	"net/http"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/jarcoal/httpmock"
)

func TestPackagist_fetchURLFromRegistry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://repo.packagist.org/p2/laravel/laravel.json",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"packages": {
				"laravel/laravel": [
					{
						"source": {
							"url": "https://github.com/laravel/laravel.git"
						}
					}
				]
			}
		}
		`)),
	)

	tests := []struct {
		name    string
		lib     types.Library
		wantURL string
	}{
		{
			name: "source_code_uri exists",
			lib: types.Library{
				Name: "laravel/laravel",
			},
			wantURL: "https://github.com/laravel/laravel.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Packagist{lib: tt.lib}
			got, _ := p.fetchURLFromRegistry(http.Client{})
			expect := tt.wantURL
			if got != expect {
				t.Errorf("get() = %v, want %v", got, expect)
			}
		})
	}
}
