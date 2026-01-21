package cmd

import (
	"context"

	"github.com/piekstra/gmail-ro/internal/gmail"
	"github.com/spf13/cobra"
)

var readJSONOutput bool

func init() {
	rootCmd.AddCommand(readCmd)
	readCmd.Flags().BoolVarP(&readJSONOutput, "json", "j", false, "Output result as JSON")
}

var readCmd = &cobra.Command{
	Use:   "read <message-id>",
	Short: "Read a single message",
	Long: `Read the full content of a Gmail message by its ID.

The message ID can be obtained from the search command output.

Examples:
  gmail-ro read 18abc123def456
  gmail-ro read 18abc123def456 --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := gmail.NewClient(ctx)
		if err != nil {
			return err
		}

		msg, err := client.GetMessage(args[0], true)
		if err != nil {
			return err
		}

		if readJSONOutput {
			return printJSON(msg)
		}

		printMessageHeader(msg, MessagePrintOptions{
			IncludeTo:   true,
			IncludeBody: true,
		})

		return nil
	},
}
