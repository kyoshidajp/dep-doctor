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
	httpmock.RegisterResponder("GET", "https://repo.packagist.org/p2/not-found/not-found.json",
		httpmock.NewStringResponder(404, heredoc.Doc(`
		{}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://repo.packagist.org/p2/unmarshal-error/unmarshal-error.json",
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
				Name: "laravel/laravel",
			},
			wantURL: "https://github.com/laravel/laravel.git",
			wantErr: false,
		},
		{
			name: "404 not found",
			lib: types.Library{
				Name: "not-found/not-found",
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
				Name: "unmarshal-error/unmarshal-error",
			},
			wantURL: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Packagist{lib: tt.lib}
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
