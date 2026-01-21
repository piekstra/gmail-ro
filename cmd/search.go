package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/piekstra/gmail-ro/internal/gmail"
	"github.com/spf13/cobra"
)

var (
	searchMaxResults int64
	searchJSONOutput bool
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().Int64VarP(&searchMaxResults, "max", "m", 10, "Maximum number of results to return")
	searchCmd.Flags().BoolVarP(&searchJSONOutput, "json", "j", false, "Output results as JSON")
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for messages",
	Long: `Search for Gmail messages using Gmail's search syntax.

Examples:
  gmail-ro search "from:alice@example.com"
  gmail-ro search "subject:meeting" --max 20
  gmail-ro search "is:unread" --json
  gmail-ro search "after:2024/01/01 before:2024/02/01"

For more query operators, see: https://support.google.com/mail/answer/7190`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := gmail.NewClient(ctx)
		if err != nil {
			return err
		}

		messages, err := client.SearchMessages(args[0], searchMaxResults)
		if err != nil {
			return err
		}

		if len(messages) == 0 {
			fmt.Println("No messages found.")
			return nil
		}

		if searchJSONOutput {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(messages)
		}

		for _, msg := range messages {
			fmt.Printf("ID: %s\n", msg.ID)
			fmt.Printf("ThreadID: %s\n", msg.ThreadID)
			fmt.Printf("From: %s\n", msg.From)
			fmt.Printf("Subject: %s\n", msg.Subject)
			fmt.Printf("Date: %s\n", msg.Date)
			if len(msg.Labels) > 0 {
				fmt.Printf("Labels: %s\n", strings.Join(msg.Labels, ", "))
			}
			if len(msg.Categories) > 0 {
				fmt.Printf("Categories: %s\n", strings.Join(msg.Categories, ", "))
			}
			fmt.Printf("Snippet: %s\n", msg.Snippet)
			fmt.Println("---")
		}

		return nil
	},
}
