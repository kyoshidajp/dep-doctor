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

		r := bundler.Diagnose(f)
		println(r)
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)
}
