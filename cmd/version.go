package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version = "1.2.0"
)

func getVersion() string {
	return fmt.Sprintf(`Version: %s
OS: %s
Arch: %s`, Version, runtime.GOOS, runtime.GOARCH)
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion())
		},
		Short: "Show version info",
	}
	return cmd
}
