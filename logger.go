package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type levelCode struct {
	Normal    string
	Important string
	Warning   string
}

var customLevel = &levelCode{
	Normal:    "37",
	Important: "35",
	Warning:   "33",
}

type CustomFormatter struct {
}

var rwMutex = &sync.RWMutex{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取调用栈信息
	_, file, line, ok := runtime.Caller(6) // 7 表示调用栈深度，调整为适合你的深度
	if !ok {
		file = "???"
		line = 0
	}
	// 仅保留文件名
	fileParts := strings.Split(file, "/")
	fileName := fileParts[len(fileParts)-1]

	var fileInfo string
	switch entry.Level {
	case logrus.TraceLevel, logrus.DebugLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		fileInfo = fmt.Sprintf("[%s:%d]", fileName, line)
	}

	var format string
	if entry.Level == logrus.ErrorLevel || entry.Level == logrus.FatalLevel || entry.Level == logrus.PanicLevel || entry.Level == logrus.WarnLevel {
		format = "%s \033[1;%dm[%s]%s%s %s\033[1;0m\n"
	} else if entry.Level == logrus.InfoLevel {
		var infoColor string
		rwMutex.Lock()
		if entry.Data != nil {
			if entry.Data["Color"] != nil {
				infoColor = entry.Data["Color"].(string)
			}
		}
		rwMutex.Unlock()

		if infoColor == "" {
			format = "%s \033[1;%dm[%s]%s%s\033[1;0m \033[1;32m%s\033[1;0m\n"
		} else {
			format = "%s \033[1;%dm[%s]%s%s\033[1;0m \033[1;" + infoColor + "m%s\033[1;0m\n"
		}
	} else {
		format = "%s \033[1;%dm[%s]%s%s\033[1;0m %s\033[1;0m\n"
	}

	var traceId string
	if entry.Data != nil {
		if entry.Data["ReqId"] != nil {
			traceId = entry.Data["ReqId"].(string)
			if traceId != "" {
				traceId, _, _ = strings.Cut(traceId, "-")
			}
		}
	}

	// 构建自定义格式
	msg := fmt.Sprintf(format,
		entry.Time.Format("02 15:04:05.000"),
		levelColor(entry.Level),
		strings.ToUpper(entry.Level.String()),
		traceId,
		fileInfo,
		entry.Message,
	)

	return []byte(msg), nil
}

type LoggerEntry struct {
	rwMutex *sync.RWMutex
	entry   *logrus.Entry
}

func (f *LoggerEntry) SetData(key string, value interface{}) {
	f.rwMutex.Lock()
	f.entry.Data[key] = value
	f.rwMutex.Unlock()
}

func (f *LoggerEntry) delData(key string) {
	f.rwMutex.Lock()
	delete(f.entry.Data, key)
	f.rwMutex.Unlock()
}

func (f *LoggerEntry) Info(args ...interface{}) {
	f.entry.Info(args...)
}

func (f *LoggerEntry) Infof(format string, args ...interface{}) {
	f.entry.Infof(format, args...)
}

func (f *LoggerEntry) InfoNormal(args ...interface{}) {
	f.SetData("Color", customLevel.Normal)
	f.entry.Info(args...)
	f.delData("Color")
}

func (f *LoggerEntry) InfoNormalf(format string, args ...interface{}) {
	f.SetData("Color", customLevel.Normal)
	f.entry.Infof(format, args...)
	f.delData("Color")
}

func (f *LoggerEntry) InfoImportant(args ...interface{}) {
	f.SetData("Color", customLevel.Important)
	f.entry.Info(args...)
	f.delData("Color")
}

func (f *LoggerEntry) InfoImportantf(format string, args ...interface{}) {
	f.SetData("Color", customLevel.Important)
	f.entry.Infof(format, args...)
	f.delData("Color")
}

func (f *LoggerEntry) InfoWarning(args ...interface{}) {
	f.SetData("Color", customLevel.Warning)
	f.entry.Info(args...)
	f.delData("Color")
}

func (f *LoggerEntry) InfoWarningf(format string, args ...interface{}) {
	f.SetData("Color", customLevel.Warning)
	f.entry.Infof(format, args...)
	f.delData("Color")
}

func (f *LoggerEntry) Warn(args ...interface{}) {
	f.entry.Warn(args...)
}

func (f *LoggerEntry) Warnf(format string, args ...interface{}) {
	f.entry.Warnf(format, args...)
}

func (f *LoggerEntry) Debug(args ...interface{}) {
	f.entry.Debug(args...)
}

func (f *LoggerEntry) Debugf(format string, args ...interface{}) {
	f.entry.Debugf(format, args...)
}

func (f *LoggerEntry) Trace(args ...interface{}) {
	f.entry.Trace(args...)
}

func (f *LoggerEntry) Tracef(format string, args ...interface{}) {
	f.entry.Tracef(format, args...)
}

func (f *LoggerEntry) Error(args ...interface{}) {
	f.entry.Error(args...)
}

func (f *LoggerEntry) Errorf(format string, args ...interface{}) {
	f.entry.Errorf(format, args...)
}

func (f *LoggerEntry) Fatal(args ...interface{}) {
	f.entry.Fatal(args...)
}

func NewLogger() *LoggerEntry {
	newLogger := logrus.New()
	newLogger.Out = os.Stdout
	newLogger.SetFormatter(&CustomFormatter{})
	newLogger.SetLevel(logrus.TraceLevel)
	custom := &LoggerEntry{
		rwMutex: &sync.RWMutex{},
		entry:   logrus.NewEntry(newLogger),
	}
	return custom
}

// levelColor 返回与日志级别相对应的 ANSI 颜色代码
func levelColor(level logrus.Level) int {
	// 颜色代码对应关系：
	// 30以下 服务器亮绿色
	// 30 服务器暗灰色 MAC暗黑色
	// 31 服务器亮红色 MAC亮红色
	// 32 服务器暗绿色 MAC暗绿色
	// 33 服务器亮橙色 MAC暗黄色
	// 34 服务器暗蓝色 MAC暗蓝色
	// 35 服务器暗红色 MAC紫色
	// 36 服务器亮蓝色 MAC亮蓝色
	// 37 服务器亮白色 MAC亮白色
	// 38/39/40 服务器亮绿色
	switch level {
	case logrus.InfoLevel:
		return 32
	case logrus.DebugLevel:
		return 36
	case logrus.TraceLevel:
		return 35
	case logrus.WarnLevel:
		return 33
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return 31 // 红色
	default:
		return 31 // 默认颜色
	}
}
