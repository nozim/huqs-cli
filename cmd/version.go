package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Huq cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Huq verision: %s\n", VERSION)
	},
}
