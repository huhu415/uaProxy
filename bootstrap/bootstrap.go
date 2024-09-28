package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// 版本信息version information
var (
	BuildDate string
	GitCommit string
	Version   string
)

var C Config

type Config struct {
	RedirPort int    `mapstructure:"redir-port"`
	UA        string `mapstructure:"User-Agent"`
}

func LoadConfig() error {
	parseFlag()
	initLog()
	logrus.SetReportCaller(true)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}
	viper.AddConfigPath(filepath.Dir(pathAbs))

	viper.AddConfigPath("/Users/hello/Projects/UA2F-go")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	if err = viper.Unmarshal(&C); err != nil {
		return err
	}

	return nil
}

func parseFlag() {
	pflag.StringP("configFile", "c", "", "config file")
	pflag.StringP("redir-port", "l", "12345", "listen address")
	pflag.BoolP("version", "v", false, "version information")
	pflag.BoolP("help", "h", false, "display help information")
	pflag.Bool("debug", false, "debug mode")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}
	if viper.GetBool("version") {
		versionInfo()
		os.Exit(0)
	}
}

func initLog() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func versionInfo() {
	logrus.Infof("version:\033[1;34m%s\033[0m", Version)
	logrus.Infof("buildDate:\033[1;34m%s\033[0m", BuildDate)
	logrus.Infof("gitCommit:\033[1;34m%s\033[0m", GitCommit)
}
