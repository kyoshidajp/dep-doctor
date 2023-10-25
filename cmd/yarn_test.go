package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchFromRubyGems(t *testing.T) {
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
			r, _ := FetchFromRubyGems(tt.gem_name)
			expect := expects[i]
			assert.Equal(t, true, strings.HasPrefix(r, expect.url))
		})
	}
}
