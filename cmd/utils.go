package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/spf13/viper"
)

func ShowVars() {
	logger.Debug("settings",
		zap.String("Bucket", BucketName),
		zap.String("Object", ObjectName),
		zap.String("FilePath", FilePath),
		zap.String("Mime", Mime),
		zap.String("Endpoint", Endpoint),
		zap.Bool("UseSSL", UseSSL),
		zap.String("Hostname", Hostname),
		zap.String("LogFile", LogFile),
	)
}

func MergeEnvFlags() {
	Endpoint = viper.GetString("Endpoint")
	Username = viper.GetString("Username")
	Password = viper.GetString("Password")
	LogFile = viper.GetString("LogFile")
	UseSSL = viper.GetBool("UseSSL")

	Username = strings.TrimSpace(Username)
	Password = strings.TrimSpace(Password)
	Endpoint = strings.TrimSpace(Endpoint)
	LogFile = strings.TrimSpace(LogFile)

	BucketName = strings.TrimSpace(BucketName)
	ObjectName = strings.TrimSpace(ObjectName)
	FilePath = strings.TrimSpace(FilePath)
	Mime = strings.TrimSpace(Mime)

	if Username == "" || Password == "" || Endpoint == "" {
		fmt.Println("Error: endpoint, username, password can NOT be nil, pls set the env variables first.")
		os.Exit(1)
	}

}

func GetS3Client() *minio.Client {

	S3Client, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(Username, Password, ""),
		Secure: UseSSL,
	})

	if err != nil {
		logger.Fatal("GetS3Client", zap.Error(err))
		return nil
	}

	_, err = S3Client.BucketExists(Ctx, BucketName)
	if err != nil {
		logger.Fatal("GetS3Client", zap.Error(err))
		return nil
	}

	return S3Client
}

func SetLogger() {
	var err error
	cfgProduction := zap.NewProductionConfig()
	cfgProduction.DisableStacktrace = true
	cfgProduction.DisableCaller = true
	if Debug {
		cfgProduction.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		cfgProduction.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	}

	if LogFile != "" {
		cfgProduction.OutputPaths = []string{LogFile, "stdout"}
		cfgProduction.ErrorOutputPaths = []string{LogFile, "stderr"}
	} else {
		cfgProduction.OutputPaths = []string{"stdout"}
		cfgProduction.ErrorOutputPaths = []string{"stderr"}
	}

	cfgProduction.InitialFields = map[string]interface{}{"host": Hostname, "app": "s3uploader"}

	logger, err = cfgProduction.Build()
	defer logger.Sync()

	if err != nil {
		panic(fmt.Sprintf("zap logger init failed: %v", err))
	}

}

func FPut() {
	user_tags := map[string]string{
		"runat": Hostname,
		"file":  FilePath,
		"mime":  Mime,
	}

	opts := minio.PutObjectOptions{}

	opts.ContentType = Mime
	opts.UserTags = user_tags

	info, err := S3Client.FPutObject(Ctx, BucketName, ObjectName, FilePath, opts)

	if err != nil {
		logger.Error("FPut",
			zap.String("Bucket", BucketName),
			zap.String("Object", ObjectName),
			zap.String("File", FilePath),
			zap.Error(err),
		)
	} else {
		logger.Info("FPut",
			zap.String("Bucket", BucketName),
			zap.String("Object", ObjectName),
			zap.String("File", FilePath),
			zap.Int64("File", info.Size),
		)
	}
}

func FGet() {
	err := S3Client.FGetObject(Ctx, BucketName, ObjectName, FilePath, minio.GetObjectOptions{})
	if err != nil {
		logger.Error("FGet",
			zap.String("Bucket", BucketName),
			zap.String("Object", ObjectName),
			zap.String("File", FilePath),
			zap.Error(err),
		)
	} else {
		logger.Info("FGet",
			zap.String("Bucket", BucketName),
			zap.String("Object", ObjectName),
			zap.String("File", FilePath),
		)
	}
}

func Stat() error {
	var result = map[string]interface{}{"code": 200, "message": ""}

	objInfo, err := S3Client.StatObject(Ctx, BucketName, ObjectName, minio.StatObjectOptions{})
	if err != nil {
		result["code"] = 404
		result["message"] = err.Error()
		logger.Error("Stat", zap.Error(err))
	} else {
		result["code"] = 0
		result["message"] = objInfo
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		jsonResult = []byte("{\"code\":500,\"message\":\"jsonfy failed\"}")
	}

	fmt.Println(string(jsonResult))

	if FilePath != "" {
		if strings.HasSuffix(FilePath, ".json") == false {
			FilePath = fmt.Sprintf("%s.json", FilePath)
		}
		SaveFile(FilePath, jsonResult)
	}
	return nil
}

func ShowFooter() {
	t_stop := time.Now().Unix() - timer_start
	if t_stop < 1 {
		t_stop = 1
	}
	logger.Debug("ShowFooter", zap.Int64("duration_second", t_stop))
}

func SaveFile(fp string, cnt []byte) error {
	if fp != "" && cnt != nil {
		if err := ioutil.WriteFile(fp, cnt, 0755); err != nil {
			logger.Error("SaveFile", zap.Error(err))
			return err
		}
	} else {
		err := errors.New("path and content cannot be nil")
		logger.Error("SaveFile", zap.Error(err))
		return err
	}

	return nil
}
