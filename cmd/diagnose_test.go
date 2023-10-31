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
			URL:       "https://github.com/faker-ruby/faker",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"concurrent-ruby": {
			Name:      "concurrent-ruby",
			URL:       "https://github.com/ruby-concurrency/concurrent-ruby",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"i18n": {
			Name:      "i18n",
			URL:       "https://github.com/ruby-i18n/i18n",
			Archived:  false,
			Ignored:   true,
			Diagnosed: true,
			IsActive:  true,
		},
		"method_source": {
			Name:      "method_source",
			URL:       "https://github.com/banister/method_source",
			Archived:  false,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  true,
		},
		"paperclip": {
			Name:      "paperclip",
			URL:       "https://github.com/thoughtbot/paperclip",
			Archived:  true,
			Ignored:   false,
			Diagnosed: true,
			IsActive:  false,
		},
		"dotenv": {
			Name:      "dotenv",
			URL:       "https://github.com/bkeepers/dotenv",
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

		doctor := NewBundlerDoctor()
		ignores := []string{"i18n"}
		diagnoses := Diagnose(doctor, f, 2, ignores)
		assert.Equal(t, expect, diagnoses)
	})
}
