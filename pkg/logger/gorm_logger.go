package logger

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormLogrusLogger struct {
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

func NewGormLogrusLogger() *GormLogrusLogger {
	return &GormLogrusLogger{
		LogLevel:                  gormlogger.Warn,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

func (l *GormLogrusLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *GormLogrusLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		WithCtx(ctx, "repository", "database").Infof(msg, data...)
	}
}

func (l *GormLogrusLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		WithCtx(ctx, "repository", "database").Warnf(msg, data...)
	}
}

func (l *GormLogrusLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		WithCtx(ctx, "repository", "database").Errorf(msg, data...)
	}
}

func (l *GormLogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	logEntry := WithCtx(ctx, "repository", "database").WithFields(logrus.Fields{
		"duration": elapsed.String(),
		"rows":     rows,
		"sql":      sql,
	})

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		logEntry.WithError(err).Error("gorm query failed")
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		// Slow query log
		logEntry.Warn("gorm slow query")
	case l.LogLevel == gormlogger.Info:
		logEntry.Info("gorm query")
	}
}
