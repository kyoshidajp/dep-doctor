package python

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPyPi_fetchURLFromRegistry(t *testing.T) {
	tests := []struct {
		name     string
		lib_name string
	}{
		{
			name:     "info.project_urls.Source_Code exists",
			lib_name: "pip",
		},
		{
			name:     "not found",
			lib_name: "not-found-xxxx",
		},
	}
	expects := []struct {
		name string
		url  string
	}{
		{
			name: "info.project_urls.Source_Code exists",
			url:  "https://github.com/pypa/pip",
		},
		{
			name: "not found",
			url:  "",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pypi{name: tt.lib_name}
			r, _ := p.fetchURLFromRegistry(http.Client{})
			expect := expects[i]
			assert.Equal(t, expect.url, r)
		})
	}
}
