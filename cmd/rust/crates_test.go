package rust

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrate_fetchURLFromRegistry(t *testing.T) {
	tests := []struct {
		name  string
		crate Crate
	}{
		{
			name: "has Crate.Repository",
			crate: Crate{
				name: "libc",
			},
		},
		{
			name: "not found",
			crate: Crate{
				name: "not-found-crate",
			},
		},
	}
	expects := []struct {
		name string
		url  string
	}{
		{
			name: "has Source.git",
			url:  "https://github.com/rust-lang/libc",
		},
		{
			name: "not found",
			url:  "",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := tt.crate.fetchURLFromRegistry(http.Client{})
			expected := expects[i].url
			assert.Equal(t, expected, actual)
		})
	}
}
