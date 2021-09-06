package cmd

import (
	"github.com/spf13/cobra"
)

var fgetCmd = &cobra.Command{
	Use:   "fget",
	Short: "download s3 object and save as local file",
	Long:  `read object from s3 remotely(--object=object-name-in-s3) and save it into local disk(--file=local-file-save-path)`,
	PreRun: func(cmd *cobra.Command, args []string) {
		ShowVars()
	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Start Downloading...")
		FGet()

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		ShowFooter()
	},
}

func init() {
	rootCmd.MarkPersistentFlagRequired("file")
	rootCmd.AddCommand(fgetCmd)

}
