package cmd

import (
	"net/http"
	"strings"
	"testing"

	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNodejs_fetchURLFromRegistry(t *testing.T) {
	tests := []struct {
		name string
		lib  types.Library
	}{
		{
			name: "source_code_uri exists",
			lib: types.Library{
				Name: "react",
			},
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
			n := Nodejs{lib: tt.lib}
			r, _ := n.fetchURLFromRegistry(http.Client{})
			expect := expects[i]
			assert.Equal(t, true, strings.HasPrefix(r, expect.url))
		})
	}
}
