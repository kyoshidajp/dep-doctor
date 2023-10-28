package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodejs_fetchURLFromRegistry(t *testing.T) {
	tests := []struct {
		name     string
		dep_name string
	}{
		{
			name:     "source_code_uri exists",
			dep_name: "react",
		},
	}
	expects := []struct {
		name string
		url  string
	}{
		{
			name: "source_code_uri exists",
			url:  "git+https://github.com/facebook/react",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Nodejs{name: tt.dep_name}
			r, _ := n.fetchURLFromRegistry()
			expect := expects[i]
			assert.Equal(t, true, strings.HasPrefix(r, expect.url))
		})
	}
}
