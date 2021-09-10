package cmd

import (
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "tag the s3 object",
	Long:  `read local file(--file=local-file-path), then put it into s3's bucket(--bucket=s3-bucket-name) with name(--object=object-name-in-s3)`,
	PreRun: func(cmd *cobra.Command, args []string) {

		ShowVars()

	},
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Start Tagging...")
		TagObject()
	},

	PostRun: func(cmd *cobra.Command, args []string) {
		ShowFooter()
	},
}

func init() {
	tagCmd.Flags().StringVar(&KV, "kv", "", "Tag Key Value Pair(s): --kv=a:b,c:d,e:f,g:h (required)")

	tagCmd.MarkFlagRequired("kv")

	rootCmd.AddCommand(tagCmd)
}
