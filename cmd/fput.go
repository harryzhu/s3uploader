package cmd

import (
	"github.com/spf13/cobra"
)

var fputCmd = &cobra.Command{
	Use:   "fput",
	Short: "upload local file into s3 objects",
	Long:  `read local file(--file=local-file-path), then put it into s3's bucket(--bucket=s3-bucket-name) with name(--object=object-name-in-s3)`,
	PreRun: func(cmd *cobra.Command, args []string) {

		ShowVars()

	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Start Uploading...")
		FPut()
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		ShowFooter()
	},
}

func init() {
	fputCmd.Flags().StringVar(&Mime, "mime", "application/octet-stream", "content-type of the object")
	rootCmd.MarkPersistentFlagRequired("file")
	rootCmd.AddCommand(fputCmd)
}
