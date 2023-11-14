package rust

import (
	"net/http"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/jarcoal/httpmock"
)

func TestCrate_fetchURLFromRegistry(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://crates.io/api/v1/crates/libc",
		httpmock.NewStringResponder(200, heredoc.Doc(`
		{
			"crate": {
				"repository": "https://github.com/rust-lang/libc"
			}
		}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://crates.io/api/v1/crates/not-found",
		httpmock.NewStringResponder(404, heredoc.Doc(`
		{}
		`)),
	)
	httpmock.RegisterResponder("GET", "https://crates.io/api/v1/crates/unmarshal-error",
		httpmock.NewStringResponder(200, ""),
	)

	tests := []struct {
		name    string
		libName string
		wantURL string
		wantErr bool
	}{
		{
			name:    "has Crate.Repository",
			libName: "libc",
			wantURL: "https://github.com/rust-lang/libc",
			wantErr: false,
		},
		{
			name:    "not found",
			libName: "not-found",
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
			c := Crate{name: tt.libName}
			got, err := c.fetchURLFromRegistry(http.Client{})
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
