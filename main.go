package main

import (
	"os"

	"github.com/kyoshidajp/dep-doctor/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
