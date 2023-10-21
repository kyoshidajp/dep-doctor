package bundler_test

import (
	"strings"
	"testing"

	"github.com/kyoshidajp/dep-doctor/cmd/ruby/bundler"
	"github.com/stretchr/testify/assert"
)

func TestFetchFromRubyGems(t *testing.T) {
	tests := []struct {
		name     string
		gem_name string
	}{
		{
			name:     "normal",
			gem_name: "rails",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bundler.FetchFromRubyGems(tt.gem_name)
			assert.Equal(t, true, strings.HasPrefix(r, "https://github.com/rails/rails"))
		})
	}
}
