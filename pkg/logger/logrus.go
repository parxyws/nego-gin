package logger

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/parxyws/nego-gin/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var fieldsLogrusLevelMap = map[string]logrus.Level{
	"trace": logrus.TraceLevel,
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
	"fatal": logrus.FatalLevel,
	"panic": logrus.PanicLevel,
}

func NewLogrusLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	level, exist := fieldsLogrusLevelMap[cfg.Logger.Level]
	if !exist {
		fmt.Printf("unknown logger level: %s\n", cfg.Logger.Level)
		level = logrus.DebugLevel
	}

	if cfg.Logger.Encoding != "text" && cfg.Logger.Encoding != "json" {
		fmt.Printf("unknown logger encoding: %s\n", cfg.Logger.Encoding)
		cfg.Logger.Encoding = "text"
	}

	var formatter logrus.Formatter
	if cfg.Logger.Encoding == "text" {
		formatter = &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			DisableColors:   false,

			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				file = fmt.Sprintf("%s:%d", f.File, f.Line)
				return function, filepath.Base(file)
			},
		}
	} else {
		formatter = &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			PrettyPrint:     true,
			CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
				file = fmt.Sprintf("%s:%d", f.File, f.Line)
				return function, filepath.Base(file)
			},
		}
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/application.log",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     28,
		Compress:   true,
	}

	logger.SetLevel(level)
	logger.SetFormatter(formatter)
	logger.SetReportCaller(cfg.Logger.Caller)
	logger.SetOutput(lumberjackLogger)

	// Configure global logrus standard logger
	logrus.SetLevel(level)
	logrus.SetFormatter(formatter)
	logrus.SetReportCaller(cfg.Logger.Caller)
	logrus.SetOutput(lumberjackLogger)

	return logger
}
