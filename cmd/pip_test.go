package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPyPiDoctor_fetchURLFromRepository(t *testing.T) {
	tests := []struct {
		name     string
		gem_name string
	}{
		{
			name:     "source_code_uri exists",
			gem_name: "pip",
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
			s := PipDoctor{}
			r, _ := s.fetchURLFromRepository(tt.gem_name)
			expect := expects[i]
			assert.Equal(t, expect.url, r)
		})
	}
}
