package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/nova2018/gologger/logger"
	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type loggerGroup struct {
	list []gormLogger.Interface
}

func newLoggerGroup() *loggerGroup {
	return &loggerGroup{list: make([]gormLogger.Interface, 0, 2)}
}

func (l *loggerGroup) AttachLogger(lg gormLogger.Interface) {
	l.list = append(l.list, lg)
}

func (l *loggerGroup) LogMode(lv gormLogger.LogLevel) gormLogger.Interface {
	newGroup := *l
	newGroup.list = make([]gormLogger.Interface, len(l.list))
	for k, v := range l.list {
		newGroup.list[k] = v.LogMode(lv)
	}
	return &newGroup
}

func (l *loggerGroup) Info(ctx context.Context, s string, i ...interface{}) {
	for _, v := range l.list {
		v.Info(ctx, s, i...)
	}
}

func (l *loggerGroup) Warn(ctx context.Context, s string, i ...interface{}) {
	for _, v := range l.list {
		v.Info(ctx, s, i...)
	}
}
func (l *loggerGroup) Error(ctx context.Context, s string, i ...interface{}) {
	for _, v := range l.list {
		v.Info(ctx, s, i...)
	}
}
func (l *loggerGroup) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	for _, v := range l.list {
		v.Trace(ctx, begin, fc, err)
	}
}

type appLogger struct {
	SlowThreshold time.Duration
	lg            logger.Logger
	LogLevel      gormLogger.LogLevel
}

func newAppLogger(lg logger.Logger) *appLogger {
	return &appLogger{
		lg: lg.WithOptions(zap.AddCallerSkip(1 + 3)),
	}
}

// LogMode log mode
func (l *appLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info print info
func (l *appLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		l.lg.Infof(ctx, msg, data...)
	}
}

// Warn print warn messages
func (l *appLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		l.lg.Warnf(ctx, msg, data...)
	}
}

// Error print error messages
func (l *appLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		l.lg.Errorf(ctx, msg, data...)
	}
}

// Trace print sql message
func (l *appLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	if l.LogLevel <= gormLogger.Silent {
		return
	}
	traceFormat := "[%.3fms] [rows:%v] %s"
	traceWarnFormat := "%s [%.3fms] [rows:%v] %s"
	traceErrFormat := "%s [%.3fms] [rows:%v] %s"

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormLogger.Error && !errors.Is(err, gormLogger.ErrRecordNotFound):
		sql, rows := fc()
		if rows == -1 {
			l.lg.Errorf(ctx, traceErrFormat, err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.lg.Errorf(ctx, traceErrFormat, err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.lg.Warnf(ctx, traceWarnFormat, slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.lg.Warnf(ctx, traceWarnFormat, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == gormLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.lg.Infof(ctx, traceFormat, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.lg.Infof(ctx, traceFormat, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
