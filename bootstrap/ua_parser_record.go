package bootstrap

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mileusna/useragent"
	"github.com/sirupsen/logrus"
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
	record    sync.Map // [string]*atomic.Int64
	filepath  string
	startTime time.Time
}

func GiveParserRecord() *ParserRecord {
	return pr
}

func NewParserRecord(ctx context.Context, filePath string) {
	pr = &ParserRecord{
		record:    sync.Map{},
		filepath:  filePath,
		startTime: time.Now(),
	}
	go func() {
		logrus.Debug("NewParserRecord finish")
		ticker := time.NewTicker(1 * time.Minute)
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
	if !C.Stats {
		return
	}
	ua := useragent.Parse(uaString)

	deviceKey := ua.OS + "-" + ua.OSVersion + " " + ua.Device
	val, _ := u.record.LoadOrStore(deviceKey, &atomic.Int64{})
	val.(*atomic.Int64).Add(1)
}

// 使用这个函数定期把统计信息写入日志文件里面
func (u *ParserRecord) WriteLog() {
	// 打开或创建文件
	file, err := os.Create(u.filepath)
	if err != nil {
		logrus.Errorln("Error creating file:", err)
		return
	}
	defer file.Close()

	// 使用 bufio 写入文件，以提高写入效率
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 写入表头，并对齐列宽
	currentTime := time.Now().Format(TIMESTAMPFORMAT)
	startTime := u.startTime.Format(TIMESTAMPFORMAT)
	_, err = writer.WriteString(fmt.Sprintf("%-50s to %-50s\n", startTime, currentTime))
	_, err = writer.WriteString(fmt.Sprintf("%-50s | %-50s\n", "Key", "Value"))
	_, err = writer.WriteString(fmt.Sprintf("%s\n", strings.Repeat("-", 100)))
	if err != nil {
		logrus.Errorln("Error writing to file:", err)
		return
	}

	// 遍历 sync.Map 并写入键值对，每列宽度固定，左对齐
	u.record.Range(func(key, value interface{}) bool {
		k := key.(string)
		v := value.(*atomic.Int64)
		line := fmt.Sprintf("%-50s | %-50d\n", k, v.Load())
		_, err := writer.WriteString(line)
		if err != nil {
			logrus.Errorln("Error writing to file:", err)
			return false // 返回 false 以停止迭代
		}
		return true // 返回 true 以继续迭代
	})

	logrus.Debug("Data successfully recorded to file.")
}
