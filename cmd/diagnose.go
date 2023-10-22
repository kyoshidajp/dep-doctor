package cmd

import (
	"os"

	"github.com/kyoshidajp/dep-doctor/cmd/ruby/bundler"
	"github.com/spf13/cobra"
)

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Diagnose packages",
	Run: func(cmd *cobra.Command, args []string) {
		f, _ := os.Open(args[0])
		defer f.Close()

		diagnoses := bundler.Diagnose(f)
		err := bundler.Report(diagnoses)
		if err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
}
