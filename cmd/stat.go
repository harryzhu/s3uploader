package cmd

import (
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat s3-file(--object=) and save the result to local-file(--file=)",
	Short: "stat object info",
	Long:  `--bucket=s3-bucket-name --object=object-name-in-s3 --file=local-file-path`,
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {
		Stat()
	},

	PostRun: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(statCmd)
}
