package cmd

import "github.com/spf13/cobra"

func AddCommand(command *cobra.Command) {
	rootCmd.AddCommand(command)
}
