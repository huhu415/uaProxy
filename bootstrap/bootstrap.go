package bootstrap

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mileusna/useragent"
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

const TIMESTAMPFORMAT = "2006-01-02 15:04:05"

var pr *ParserRecord

type ParserRecord struct {
	record   map[string]*atomic.Int64
	filepath string
}

func GiveParserRecord() *ParserRecord {
	return pr
}

func NewParserRecord(ctx context.Context, filePath string) {
	pr = &ParserRecord{
		record:   make(map[string]*atomic.Int64),
		filepath: filePath,
	}
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				pr.WriteLog() // 最后一次写入日志
				return
			case <-ticker.C:
				pr.WriteLog()
			}
		}
	}()
}

func (u *ParserRecord) ParserAndRecord(uaString string) {
	ua := useragent.Parse(uaString)

	deviceKey := ua.OS + " " + ua.Device
	if _, exists := u.record[deviceKey]; !exists {
		u.record[deviceKey] = &atomic.Int64{}
	}
	u.record[deviceKey].Add(1)
}

// 使用这个函数定期把统计信息写入日志文件里面
func (u *ParserRecord) WriteLog() {
	// 打开或创建文件
	file, err := os.Create(u.filepath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 使用 bufio 写入文件，以提高写入效率
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 写入表头，并对齐列宽
	currentTime := time.Now().Format(TIMESTAMPFORMAT)
	_, err = writer.WriteString(fmt.Sprintf("%-50s | %-50s\n", "Current Time", currentTime))
	_, err = writer.WriteString(fmt.Sprintf("%-50s | %-50s\n", "Key", "Value"))
	_, err = writer.WriteString(fmt.Sprintf("%s\n", strings.Repeat("-", 100)))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// 遍历 map 并写入键值对，每列宽度固定，左对齐
	for key, value := range u.record {
		line := fmt.Sprintf("%-50s | %-50d\n", key, value.Load())
		_, err := writer.WriteString(line)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	logrus.Debug("Data successfully recorded to file.")
}

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
	pflag.String("User-Agent", "fffffffffffffff", "User-Agent value")
	pflag.Bool("debug", false, "debug mode")
	pflag.Bool("stats", false, "enable statistics collection")
	pflag.String("stats-config", csvPath, "configuration file")
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
		TimestampFormat: TIMESTAMPFORMAT,
	})
	logrus.SetOutput(os.Stdout)
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
}
