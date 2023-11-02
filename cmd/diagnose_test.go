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

func TestDoctors_PackageManagers(t *testing.T) {
	tests := []struct {
		name    string
		doctors Doctors
	}{
		{
			name: "doctors",
			doctors: Doctors{
				"package1": nil,
				"package2": nil,
				"package3": nil,
			},
		},
	}
	expects := []struct {
		name     string
		packages []string
	}{
		{
			name:     "doctors",
			packages: []string{"package1", "package2", "package3"},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.doctors.PackageManagers()
			expect := expects[i].packages
			assert.Equal(t, expect, actual)
		})
	}
}

func TestDoctors_UnknownErrorMessage(t *testing.T) {
	tests := []struct {
		name    string
		doctors Doctors
	}{
		{
			name: "doctors",
			doctors: Doctors{
				"p1": nil,
				"p2": nil,
				"p3": nil,
			},
		},
	}
	expects := []struct {
		name    string
		message string
	}{
		{
			name:    "doctors",
			message: "Unknown package manager: xxx. You can choose from [p1, p2, p3]",
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.doctors.UnknownErrorMessage("xxx")
			expect := expects[i].message
			assert.Equal(t, expect, actual)
		})
	}
}
