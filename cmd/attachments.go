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
  gmro attachments list 18abc123def456
  gmro attachments download 18abc123def456 --all
  gmro attachments download 18abc123def456 --filename report.pdf`,
}
