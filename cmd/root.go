package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	minio "github.com/minio/minio-go/v7"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Ctx        context.Context
	S3Client   *minio.Client
	Prefix     string
	Hostname   string
	Username   string
	Password   string
	Endpoint   string
	UseSSL     bool
	BucketName string
	ObjectName string
	MaxSize_MB int64
	FilePath   string
	Mime       string
	KV         string
	LogFile    string
	Debug      bool
)
var (
	timer_start int64
)

var (
	logger *zap.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "s3uploader fput|fget|stat|tag",
	Short: "s3uploader",
	Long: `Usage:
1)set env variables first: 
	export S3UPLOADER_USERNAME=YOUR_API_KEY
	export S3UPLOADER_PASSWORD=YOUR_API_SECRET
	export S3UPLOADER_ENDPOINT=MINIO_S3_SERVICE_DOMAIN_NAME
	export S3UPLOADER_USESSL=false or true(depends on your server side);
	export S3UPLOADER_LOGFILE=/Users/harryzhu/logs/s3uploader.log;
	
  or if you set parameter --prefix=ALPHA, env variables should be like:
	export ALPHA_S3UPLOADER_USERNAME=YOUR_API_KEY
	export ALPHA_S3UPLOADER_PASSWORD=YOUR_API_SECRET
	export ALPHA_S3UPLOADER_ENDPOINT=MINIO_S3_SERVICE_DOMAIN_NAME
	export ALPHA_S3UPLOADER_USESSL=false or true(depends on your server side);
	export ALPHA_S3UPLOADER_LOGFILE=/Users/harryzhu/logs/s3uploader_alpha.log;
	
2)upload file: ./s3uploader fput --bucket=test --object=s3/object/path_or_name.png --file=local/path/of/your/file.png --mime="image/png" --debug

3)download file: ./s3uploader fget --bucket=test --object=s3/object/path_or_name.png --file=local/path/of/your/file.png

4)tag object: ./s3uploader tag --bucket=test --object=s3/object/path_or_name.png --kv="project:BeChangedByTheWorld,owner:harryzhu,year:2020"

5)stat object: ./s3uploader stat --bucket=test --object=s3/object/path_or_name.png
  stat object: ./s3uploader stat --bucket=test --object=s3/object/path_or_name.png --file=local/path/of/exporting/result/for/save.json
	`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		MergeEnvFlags()
		SetLogger()
		Ctx = context.Background()
		S3Client = GetS3Client()

	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Thank you for choosing s3uploader, if any bug, pls contact harryzhu")

	},

	PersistentPostRun: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&Prefix, "prefix", "", "env variables prefix, format: {prefix_}S3UPLOADER_{USERNAME|PASSWORD|ENDPOINT|USESSL}=YOUR_VALUE")
	rootCmd.PersistentFlags().StringVar(&BucketName, "bucket", "", "s3 bucket name (required)")
	rootCmd.PersistentFlags().StringVar(&ObjectName, "object", "", "s3 object name (required)")
	rootCmd.PersistentFlags().StringVar(&FilePath, "file", "", "local file path (required)")
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "is debugging: true/false")

	rootCmd.MarkPersistentFlagRequired("bucket")
	rootCmd.MarkPersistentFlagRequired("object")

	timer_start = time.Now().Unix()

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	Hostname, err = os.Hostname()
	if err != nil {
		Hostname = "unknown"
		fmt.Println("cannot get the hostname.")
	}

	if Prefix != "" {
		viper.SetEnvPrefix(fmt.Sprintf("%s_%s", Prefix, "S3UPLOADER"))
	} else {
		viper.SetEnvPrefix("S3UPLOADER")
	}

	viper.BindEnv("USERNAME")
	viper.BindEnv("PASSWORD")
	viper.BindEnv("ENDPOINT")
	viper.BindEnv("USESSL")
	viper.BindEnv("LOGFILE")
	viper.BindEnv("MAXSIZE_MB")

	viper.AutomaticEnv() // read in environment variables that match

}
