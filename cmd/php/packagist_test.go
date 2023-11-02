package php

import (
	"net/http"
	"strings"
	"testing"

	"github.com/aquasecurity/go-dep-parser/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestPackagist_fetchURLFromRegistry(t *testing.T) {
	tests := []struct {
		name string
		lib  types.Library
	}{
		{
			name: "source_code_uri exists",
			lib: types.Library{
				Name: "laravel/laravel",
			},
		},
	}
	expects := []struct {
		name string
		url  string
	}{
		{
			name: "source_code_uri exists",
			url:  "https://github.com/laravel/laravel",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Packagist{lib: tt.lib}
			r, _ := p.fetchURLFromRegistry(http.Client{})
			expect := expects[i]
			assert.Equal(t, true, strings.HasPrefix(r, expect.url))
		})
	}
}
