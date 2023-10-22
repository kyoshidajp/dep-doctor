package bundler_test

import (
	"os"
	"testing"

	"github.com/kyoshidajp/dep-doctor/cmd/ruby/bundler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	expect := map[string]bundler.Diagnosis{
		"faker": {
			Name:      "faker",
			Url:       "https://github.com/faker-ruby/faker",
			Archived:  false,
			Diagnosed: true,
		},
		"concurrent-ruby": {
			Name:      "concurrent-ruby",
			Url:       "https://github.com/ruby-concurrency/concurrent-ruby",
			Archived:  false,
			Diagnosed: true,
		},
		"i18n": {
			Name:      "i18n",
			Url:       "https://github.com/ruby-i18n/i18n",
			Archived:  false,
			Diagnosed: true,
		},
		"method_source": {
			Name:      "method_source",
			Url:       "https://github.com/banister/method_source",
			Archived:  false,
			Diagnosed: true,
		},
		"paperclip": {
			Name:      "paperclip",
			Url:       "https://github.com/thoughtbot/paperclip",
			Archived:  true,
			Diagnosed: true,
		},
		"dotenv": {
			Name:      "dotenv",
			Url:       "",
			Archived:  false,
			Diagnosed: false,
		},
	}

	t.Run("test", func(t *testing.T) {
		f, err := os.Open("testdata/Gemfile.lock")
		require.NoError(t, err)
		defer f.Close()

		r := bundler.Diagnose(f)
		assert.Equal(t, expect, r)
	})
}
