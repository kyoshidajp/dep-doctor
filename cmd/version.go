package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Revision = "2"
	Version  = "1"
)

func getVersion() string {
	return fmt.Sprintf(`Version: %s
Revision: %s
OS: %s
Arch: %s`, Version, Revision, runtime.GOOS, runtime.GOARCH)
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
