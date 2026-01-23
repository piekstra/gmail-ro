package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-cli-collective/gmail-ro/internal/gmail"
	ziputil "github.com/open-cli-collective/gmail-ro/internal/zip"
	"github.com/spf13/cobra"
)

var (
	downloadFilename string
	downloadDir      string
	downloadExtract  bool
	downloadAll      bool
)

func init() {
	attachmentsCmd.AddCommand(downloadAttachmentsCmd)
	downloadAttachmentsCmd.Flags().StringVarP(&downloadFilename, "filename", "f", "",
		"Download only attachment with this filename")
	downloadAttachmentsCmd.Flags().StringVarP(&downloadDir, "output", "o", ".",
		"Directory to save attachments")
	downloadAttachmentsCmd.Flags().BoolVarP(&downloadExtract, "extract", "e", false,
		"Extract zip files after download")
	downloadAttachmentsCmd.Flags().BoolVarP(&downloadAll, "all", "a", false,
		"Download all attachments (required if no --filename specified)")
}

var downloadAttachmentsCmd = &cobra.Command{
	Use:   "download <message-id>",
	Short: "Download attachments from a message",
	Long: `Download attachments from a Gmail message to local disk.

By default, requires --filename to specify which attachment to download,
or --all to download all attachments.

Zip files can be automatically extracted with --extract flag.

Examples:
  gmro attachments download 18abc123def456 --filename report.pdf
  gmro attachments download 18abc123def456 --all
  gmro attachments download 18abc123def456 --all --output ~/Downloads
  gmro attachments download 18abc123def456 --filename archive.zip --extract`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if downloadFilename == "" && !downloadAll {
			return fmt.Errorf("must specify --filename or --all")
		}

		client, err := newGmailClient()
		if err != nil {
			return err
		}

		messageID := args[0]
		attachments, err := client.GetAttachments(messageID)
		if err != nil {
			return err
		}

		if len(attachments) == 0 {
			fmt.Println("No attachments found for message.")
			return nil
		}

		// Filter by filename if specified
		var toDownload []*gmail.Attachment
		for _, att := range attachments {
			if downloadFilename == "" || att.Filename == downloadFilename {
				toDownload = append(toDownload, att)
			}
		}

		if len(toDownload) == 0 {
			return fmt.Errorf("attachment not found: %s", downloadFilename)
		}

		// Create output directory if needed
		if err := os.MkdirAll(downloadDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		// Download each attachment
		for _, att := range toDownload {
			data, err := downloadAttachment(client, messageID, att)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", att.Filename, err)
				continue
			}

			outputPath := filepath.Join(downloadDir, att.Filename)
			if err := saveAttachment(outputPath, data); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving %s: %v\n", att.Filename, err)
				continue
			}

			fmt.Printf("Downloaded: %s (%s)\n", outputPath, formatSize(int64(len(data))))

			// Extract if zip and --extract flag
			if downloadExtract && isZipFile(att.Filename, att.MimeType) {
				extractDir := filepath.Join(downloadDir,
					strings.TrimSuffix(att.Filename, filepath.Ext(att.Filename)))
				if err := ziputil.Extract(outputPath, extractDir, ziputil.DefaultOptions()); err != nil {
					fmt.Fprintf(os.Stderr, "Error extracting %s: %v\n", att.Filename, err)
				} else {
					fmt.Printf("Extracted to: %s\n", extractDir)
				}
			}
		}

		return nil
	},
}

func downloadAttachment(client *gmail.Client, messageID string, att *gmail.Attachment) ([]byte, error) {
	if att.AttachmentID != "" {
		return client.DownloadAttachment(messageID, att.AttachmentID)
	}
	return client.DownloadInlineAttachment(messageID, att.PartID)
}

func saveAttachment(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func isZipFile(filename, mimeType string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".zip" ||
		mimeType == "application/zip" ||
		mimeType == "application/x-zip-compressed"
}
