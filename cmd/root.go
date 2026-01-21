package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "gmail-ro",
	Short: "A read-only Gmail CLI tool",
	Long: `gmail-ro is a command-line interface for reading Gmail messages.

It provides read-only access to your Gmail account, allowing you to:
- Search messages by query
- Read individual messages
- View full conversation threads

This tool uses OAuth2 for authentication and only requests read-only
permissions (gmail.readonly scope).`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("gmail-ro %s\n", Version)
		return nil
	},
}
