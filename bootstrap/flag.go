package bootstrap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds the configuration options
type Config struct {
	RedirPort   int    `mapstructure:"redir-port"`
	UserAgent   string `mapstructure:"User-Agent"`
	Debug       bool   `mapstructure:"debug"`
	Stats       bool   `mapstructure:"stats"`
	StatsConfig string `mapstructure:"stats-config"`
}

// Global variable to hold the configuration
var C Config

func LoadConfig() {
	parseFlag()
	initLog()
	logrus.SetReportCaller(true)
}

func parseFlag() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	csvPath := filepath.Join(filepath.Dir(exePath), "stats-config.csv")

	pflag.String("redir-port", "12345", "listen address")
	pflag.String("User-Agent", "MicroMessenger Client", "User-Agent value")
	pflag.Bool("debug", false, "debug mode")
	pflag.Bool("stats", false, "enable statistics collection")
	pflag.String("stats-config", csvPath, "configuration file")
	pflag.BoolP("version", "v", false, "version information")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	if err = viper.Unmarshal(&C); err != nil {
		log.Printf("unmarshal config file error: %+v\n", err)
		return
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
		TimestampFormat: TIMESTAMPFORMAT,
	})
	logrus.SetOutput(os.Stdout)
	if C.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
