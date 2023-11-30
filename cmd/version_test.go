package cmd

import (
	"fmt"
	"runtime"
	"testing"
)

func TestGetVersion(t *testing.T) {
	got := getVersion()
	want := fmt.Sprintf(`Version: %s
OS: %s
Arch: %s`, Version, runtime.GOOS, runtime.GOARCH)

	if want != got {
		t.Fatalf("unexpected version info. want: %s, got: %s", want, got)
	}
}
