package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd(out, errOut io.Writer) (*cobra.Command, error) {
	o := &DiagnoseOption{}
	o.Out = out
	o.ErrOut = errOut

	cmd := &cobra.Command{
		Use:           "dep-doctor",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.AddCommand(
		newDiagnoseCmd(o),
	)
	cmd.AddCommand(
		newVersionCmd(),
	)
	cmd.SetOut(out)
	cmd.SetErr(errOut)

	return cmd, nil
}

func Execute() int {
	o := os.Stdout
	e := os.Stderr

	rootCmd, err := newRootCmd(o, e)
	if err != nil {
		fmt.Fprintln(e, err)
		return 1
	}

	if err = rootCmd.Execute(); err != nil {
		fmt.Fprintln(e, err)
		return 1
	}

	return 0
}
