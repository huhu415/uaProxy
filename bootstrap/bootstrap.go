package bootstrap

import (
	"fmt"
	"os"

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

func LoadConfig() {
	parseFlag()
	initLog()
	logrus.SetReportCaller(true)
}

func parseFlag() {
	pflag.String("redir-port", "12345", "listen address")
	pflag.String("User-Agent", "fffffffffffffff", "User-Agent value")
	pflag.Bool("debug", false, "debug mode")
	pflag.BoolP("version", "v", false, "version information")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	if viper.GetBool("version") {
		fmt.Printf("version:\033[1;34m%s\033[0m\n", Version)
		fmt.Printf("buildDate:\033[1;34m%s\033[0m\n", BuildDate)
		fmt.Printf("gitCommit:\033[1;34m%s\033[0m\n", GitCommit)
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
