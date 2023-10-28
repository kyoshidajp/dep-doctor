package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPyPi_fetchURLFromRegistry(t *testing.T) {
	tests := []struct {
		name     string
		dep_name string
	}{
		{
			name:     "source_code_uri exists",
			dep_name: "pip",
		},
	}
	expects := []struct {
		name string
		url  string
	}{
		{
			name: "source_code_uri exists",
			url:  "https://github.com/pypa/pip",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pypi{name: tt.dep_name}
			r, _ := p.fetchURLFromRepository()
			expect := expects[i]
			assert.Equal(t, expect.url, r)
		})
	}
}
