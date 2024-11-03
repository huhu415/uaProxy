package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const UA = "User-Agent"

// Config holds the configuration options
type Config struct {
	RedirPort   int               `mapstructure:"redir-port"`
	Headers     map[string]string `mapstructure:"headers"`
	Debug       bool              `mapstructure:"debug"`
	Stats       bool              `mapstructure:"stats"`
	StatsConfig string            `mapstructure:"stats-config"`
}

// Global variable to hold the configuration
var C Config

func LoadConfig() {
	parseFlags()
	initLog()
	logrus.SetReportCaller(true)
}

func parseFlags() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	csvPath := filepath.Join(filepath.Dir(exePath), "stats-config.csv")

	pflag.String("redir-port", "12345", "listen address")
	pflag.StringToString("headers", map[string]string{
		UA: "MicroMessenger Client",
	}, "custom headers")
	pflag.Bool("debug", false, "debug mode")
	pflag.Bool("stats", false, "enable"+UA+"statistics collection")
	pflag.String("stats-config", csvPath, "configuration file")
	pflag.BoolP("version", "v", false, "version information")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logrus.Errorf("bind flags error: %+v\n", err)
		panic(err)
	}

	if err = viper.Unmarshal(&C); err != nil {
		logrus.Errorf("unmarshal config file error: %+v\n", err)
		return
	}

	if err := checkValid(); err != nil {
		logrus.Errorf("invalid configuration: %v", err)
		os.Exit(1)
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

func checkValid() error {
	if C.Headers[UA] == "" && C.Stats {
		return fmt.Errorf("User-Agent is required when stats is enabled")
	}
	return nil
}
