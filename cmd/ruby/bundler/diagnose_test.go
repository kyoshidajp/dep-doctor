package bundler_test

import (
	"os"
	"testing"

	"github.com/kyoshidajp/dep-doctor/cmd/ruby/bundler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{
			name: "normal",
			file: "testdata/Gemfile.lock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.file)
			require.NoError(t, err)
			defer f.Close()

			r := bundler.Diagnose(f)
			assert.Equal(t, "bundler", r)
		})
	}
}
