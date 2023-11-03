package ruby

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRubyGems_fetchURLFromRegistry(t *testing.T) {
	tests := []struct {
		name     string
		gem_name string
	}{
		{
			name:     "source_code_uri exists",
			gem_name: "rails",
		},
		{
			name:     "no source_code_uri, but homepage_uri exists",
			gem_name: "minitest",
		},
	}
	expects := []struct {
		name string
		url  string
	}{
		{
			name: "source_code_uri exists",
			url:  "https://github.com/rails/rails",
		},
		{
			name: "no source_code_uri, but homepage_uri exists",
			url:  "https://github.com/minitest/minitest",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := RubyGems{name: tt.gem_name}
			r, _ := g.fetchURLFromRegistry(http.Client{})
			expect := expects[i]
			assert.Equal(t, true, strings.HasPrefix(r, expect.url))
		})
	}
}
