package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(attachmentsCmd)
}

var attachmentsCmd = &cobra.Command{
	Use:   "attachments",
	Short: "Manage message attachments",
	Long: `List and download attachments from Gmail messages.

This command group provides read-only access to message attachments.
Use 'list' to view attachment metadata and 'download' to save files locally.

Examples:
  gmail-ro attachments list 18abc123def456
  gmail-ro attachments download 18abc123def456 --all
  gmail-ro attachments download 18abc123def456 --filename report.pdf`,
}
