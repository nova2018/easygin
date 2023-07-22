package logger

import (
	"context"
	"github.com/nova2018/gologger/logger"
	"go.uber.org/zap"
	"sync"
)

var (
	_proxyLogger logger.Logger
	_proxyLock   = &sync.Mutex{}
)

func proxy(ctx context.Context) logger.Logger {
	if ctx != nil {
		lg := getLoggerWithContext(ctx)
		if lg != nil {
			return lg
		}
	}
	if _proxyLogger == nil {
		_proxyLock.Lock()
		if _proxyLogger == nil {
			_proxyLogger = Default().WithOptions(zap.AddCallerSkip(1))
		}
		_proxyLock.Unlock()
	}
	return _proxyLogger
}

func Named(name string) logger.Logger {
	return proxy(nil).Named(name)
}

func WithOptions(opts ...zap.Option) logger.Logger {
	return proxy(nil).WithOptions(opts...)
}

func Debugf(ctx context.Context, message string, fmtArgs ...interface{}) {
	proxy(ctx).Debugf(ctx, message, fmtArgs...)
}

func Infof(ctx context.Context, message string, fmtArgs ...interface{}) {
	proxy(ctx).Infof(ctx, message, fmtArgs...)
}

func Warnf(ctx context.Context, message string, fmtArgs ...interface{}) {
	proxy(ctx).Warnf(ctx, message, fmtArgs...)
}

func Warningf(ctx context.Context, message string, fmtArgs ...interface{}) {
	proxy(ctx).Warningf(ctx, message, fmtArgs...)
}

func Errorf(ctx context.Context, message string, fmtArgs ...interface{}) {
	proxy(ctx).Errorf(ctx, message, fmtArgs...)
}

func Fatalf(ctx context.Context, message string, fmtArgs ...interface{}) {
	proxy(ctx).Fatalf(ctx, message, fmtArgs...)
}

func DPanicf(ctx context.Context, message string, fmtArgs ...interface{}) {
	proxy(ctx).DPanicf(ctx, message, fmtArgs...)
}

func Panicf(ctx context.Context, message string, fmtArgs ...interface{}) {
	proxy(ctx).Panicf(ctx, message, fmtArgs...)
}

func Debugz(ctx context.Context, message string, fields ...zap.Field) {
	proxy(ctx).Debugz(ctx, message, fields...)
}

func Infoz(ctx context.Context, message string, fields ...zap.Field) {
	proxy(ctx).Infoz(ctx, message, fields...)
}

func Warnz(ctx context.Context, message string, fields ...zap.Field) {
	proxy(ctx).Warnz(ctx, message, fields...)
}

func Warningz(ctx context.Context, message string, fields ...zap.Field) {
	proxy(ctx).Warningz(ctx, message, fields...)
}

func Errorz(ctx context.Context, message string, fields ...zap.Field) {
	proxy(ctx).Errorz(ctx, message, fields...)
}

func Fatalz(ctx context.Context, message string, fields ...zap.Field) {
	proxy(ctx).Fatalz(ctx, message, fields...)
}

func DPanicz(ctx context.Context, message string, fields ...zap.Field) {
	proxy(ctx).DPanicz(ctx, message, fields...)
}

func Panicz(ctx context.Context, message string, fields ...zap.Field) {
	proxy(ctx).Panicz(ctx, message, fields...)
}

func Debug(ctx context.Context, message string, context logger.Fields) {
	proxy(ctx).Debug(ctx, message, context)
}

func Info(ctx context.Context, message string, context logger.Fields) {
	proxy(ctx).Info(ctx, message, context)
}

func Warn(ctx context.Context, message string, context logger.Fields) {
	proxy(ctx).Warn(ctx, message, context)
}

func Warning(ctx context.Context, message string, context logger.Fields) {
	proxy(ctx).Warning(ctx, message, context)
}

func Error(ctx context.Context, message string, context logger.Fields) {
	proxy(ctx).Error(ctx, message, context)
}

func Fatal(ctx context.Context, message string, context logger.Fields) {
	proxy(ctx).Fatal(ctx, message, context)
}

func DPanic(ctx context.Context, message string, context logger.Fields) {
	proxy(ctx).DPanic(ctx, message, context)
}

func Panic(ctx context.Context, message string, context logger.Fields) {
	proxy(ctx).Panic(ctx, message, context)
}

func Debugw(ctx context.Context, message string, keysAndValues ...interface{}) {
	proxy(ctx).Debugw(ctx, message, keysAndValues...)
}

func Infow(ctx context.Context, message string, keysAndValues ...interface{}) {
	proxy(ctx).Infow(ctx, message, keysAndValues...)
}

func Warnw(ctx context.Context, message string, keysAndValues ...interface{}) {
	proxy(ctx).Warnw(ctx, message, keysAndValues...)
}

func Warningw(ctx context.Context, message string, keysAndValues ...interface{}) {
	proxy(ctx).Warningw(ctx, message, keysAndValues...)
}

func Errorw(ctx context.Context, message string, keysAndValues ...interface{}) {
	proxy(ctx).Errorw(ctx, message, keysAndValues...)
}

func Fatalw(ctx context.Context, message string, keysAndValues ...interface{}) {
	proxy(ctx).Fatalw(ctx, message, keysAndValues...)
}

func DPanicw(ctx context.Context, message string, keysAndValues ...interface{}) {
	proxy(ctx).DPanicw(ctx, message, keysAndValues...)
}

func Panicw(ctx context.Context, message string, keysAndValues ...interface{}) {
	proxy(ctx).Panicw(ctx, message, keysAndValues...)
}
