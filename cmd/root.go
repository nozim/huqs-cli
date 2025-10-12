package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "huq",
	Short: "Huq CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.SetArgs([]string{"version"})
		cmd.Execute()
		//stat, _ := os.Stdin.Stat()
		//if (stat.Mode() & os.ModeCharDevice) == 0 {
		//	data, _ := os.ReadFile("/dev/stdin")
		//	fmt.Println("Received stdin input:")
		//	fmt.Println(string(data))
		//} else {
		//	fmt.Println("Huq CLI. Use `huq init` or `huq build`")
		//}
	},
}

func Execute() {
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(playCmd)
	rootCmd.AddCommand(deployCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
