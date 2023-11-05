package ruby

import (
	"os"
	"reflect"
	"testing"
)

func TestBundlerDoctor_Libraries(t *testing.T) {
	t.Run("read bundler lock file", func(t *testing.T) {
		d := NewBundlerDoctor()
		f, _ := os.Open("bundler/testdata/Gemfile.lock")
		libs := d.Libraries(f)
		var libNames []string
		for _, v := range libs {
			libNames = append(libNames, v.Name)
		}

		expect := []string{"concurrent-ruby", "dotenv", "faker", "i18n", "method_source"}
		if !reflect.DeepEqual(libNames, expect) {
			t.Errorf("get: %v, want: %v", libNames, expect)
		}
	})
}
