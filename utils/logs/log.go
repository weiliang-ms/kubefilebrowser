package logs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"time"
)

var logger = logrus.New()

// 封装logRus.Fields
type Fields logrus.Fields

const skip = 2

func SetLogLevel(level string) {
	LogLevel, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	logger.Level = LogLevel
}
func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

// gin access log
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 捕抓异常
		defer func() {
			if err := recover(); err != nil {
				Error(err)
			}
		}()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 请求协议
		proto := c.Request.Proto
		// 请求ID
		requestID := c.Request.Header.Get("X-Request-Id")
		// Content-Type
		ContentType := c.Writer.Header().Get("Content-Type")
		// 日志格式
		logger.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUri,
			"proto":        proto,
			"request_id":   requestID,
			"content_type": ContentType,
		}).Info("AccessLog")
	}
}

// Debug
func Debug(msg ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skip)
		entry.Debug(msg)
	}
}

// 带有field的Debug
func DebugWithFields(msg interface{}, field Fields) {
	if logger.Level >= logrus.DebugLevel {
		entry := logger.WithFields(logrus.Fields(field))
		entry.Data["file"] = fileInfo(skip)
		entry.Debug(msg)
	}
}

// Info
func Info(msg ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skip)
		entry.Info(msg...)
	}
}

// 带有field的Info
func InfoWithFields(msg interface{}, field Fields) {
	if logger.Level >= logrus.InfoLevel {
		entry := logger.WithFields(logrus.Fields(field))
		entry.Data["file"] = fileInfo(skip)
		entry.Info(msg)
	}
}

// Warn
func Warn(msg ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skip)
		entry.Warn(msg...)
	}
}

// 带有Field的Warn
func WarnWithFields(msg interface{}, field Fields) {
	if logger.Level >= logrus.WarnLevel {
		entry := logger.WithFields(logrus.Fields(field))
		entry.Data["file"] = fileInfo(skip)
		entry.Warn(msg)
	}
}

// Error
func Error(msg ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skip)
		entry.Error(msg...)
	}
}

// 带有Fields的Error
func ErrorWithFields(msg interface{}, field Fields) {
	if logger.Level >= logrus.ErrorLevel {
		entry := logger.WithFields(logrus.Fields(field))
		entry.Data["file"] = fileInfo(skip)
		entry.Error(msg)
	}
}

// Fatal
func Fatal(msg ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skip)
		entry.Fatal(msg...)
	}
}

// 带有Field的Fatal
func FatalWithFields(msg interface{}, field Fields) {
	if logger.Level >= logrus.FatalLevel {
		entry := logger.WithFields(logrus.Fields(field))
		entry.Data["file"] = fileInfo(skip)
		entry.Fatal(msg)
	}
}

// Panic
func Panic(msg ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(skip)
		entry.Panic(msg...)
	}
}

// 带有Field的Panic
func PanicWithFields(msg interface{}, field Fields) {
	if logger.Level >= logrus.PanicLevel {
		entry := logger.WithFields(logrus.Fields(field))
		entry.Data["file"] = fileInfo(skip)
		entry.Panic(msg)
	}
}
func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
