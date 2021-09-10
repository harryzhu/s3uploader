package cmd

import (
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "tag the s3 object",
	Long:  `tagging the object(--object=object-name-in-s3) in s3's bucket(--bucket=s3-bucket-name) with KV pairs(--kv="key1:value1,key2:value2,key3:value3") `,
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
