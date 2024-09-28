package bootstrap

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var C Config

type Config struct {
	RedirPort int    `mapstructure:"redir-port"`
	UA        string `mapstructure:"User-Agent"`
}

func LoadConfig() error {
	logrus.SetLevel(logrus.DebugLevel)
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
