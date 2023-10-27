package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	expect := map[string]Diagnosis{
		"faker": {
			Name:      "faker",
			Url:       "https://github.com/faker-ruby/faker",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"concurrent-ruby": {
			Name:      "concurrent-ruby",
			Url:       "https://github.com/ruby-concurrency/concurrent-ruby",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"i18n": {
			Name:      "i18n",
			Url:       "https://github.com/ruby-i18n/i18n",
			Archived:  false,
			Ignored:   true,
			Diagnosed: true,
			IsActive:  true,
		},
		"method_source": {
			Name:      "method_source",
			Url:       "https://github.com/banister/method_source",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"paperclip": {
			Name:      "paperclip",
			Url:       "https://github.com/thoughtbot/paperclip",
			Archived:  true,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  false,
		},
		"dotenv": {
			Name:      "dotenv",
			Url:       "https://github.com/bkeepers/dotenv",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
	}

	t.Run("test", func(t *testing.T) {
		f, err := os.Open("ruby/bundler/testdata/Gemfile.lock")
		require.NoError(t, err)
		defer f.Close()

		doctor := NewDepartment(NewBundlerDoctor())
		ignores := []string{"i18n"}
		diagnoses := doctor.Diagnose(f, 2, ignores)
		assert.Equal(t, expect, diagnoses)
	})
}
