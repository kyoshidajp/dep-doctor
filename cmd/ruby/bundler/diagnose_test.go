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

func ExampleReport() {
	tests := []struct {
		name      string
		diagnoses []bundler.Diagnosis
	}{
		{
			name: "normal",
			diagnoses: []bundler.Diagnosis{
				{
					Name:      "rails",
					Url:       "https://github.com/rails/rails",
					Archived:  false,
					Diagnosed: true,
				},
				{
					Name:      "paperclip",
					Url:       "https://github.com/thoughtbot/paperclip",
					Archived:  true,
					Diagnosed: true,
				},
				{
					Name:      "dotenv",
					Url:       "",
					Archived:  false,
					Diagnosed: false,
				},
			},
		},
	}

	for _, tt := range tests {
		bundler.Report(tt.diagnoses)
		// Output:
		// aaa
	}
}
