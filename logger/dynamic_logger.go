package logger

import (
	"context"
	"github.com/nova2018/gologger/logger"
	"go.uber.org/zap"
)

type dynamicLogger struct {
	opts []func(logger.Logger) logger.Logger
}

func Dynamic() logger.Logger {
	return &dynamicLogger{}
}

func (d *dynamicLogger) Sync() error {
	return nil
}

func (d *dynamicLogger) Named(s string) logger.Logger {
	d.opts = append(d.opts, func(lg logger.Logger) logger.Logger {
		return lg.Named(s)
	})
	return d
}

func (d *dynamicLogger) WithOptions(opts ...zap.Option) logger.Logger {
	d.opts = append(d.opts, func(l logger.Logger) logger.Logger {
		return l.WithOptions(opts...)
	})
	return d
}

func (d *dynamicLogger) logger(ctx context.Context) logger.Logger {
	lg := proxy(ctx)
	for _, fn := range d.opts {
		lg = fn(lg)
	}
	return lg
}

func (d *dynamicLogger) Debugf(ctx context.Context, message string, fmtArgs ...interface{}) {
	d.logger(ctx).Debugf(ctx, message, fmtArgs...)
}

func (d *dynamicLogger) Infof(ctx context.Context, message string, fmtArgs ...interface{}) {
	d.logger(ctx).Infof(ctx, message, fmtArgs...)
}

func (d *dynamicLogger) Warnf(ctx context.Context, message string, fmtArgs ...interface{}) {
	d.logger(ctx).Warnf(ctx, message, fmtArgs...)
}

func (d *dynamicLogger) Warningf(ctx context.Context, message string, fmtArgs ...interface{}) {
	d.logger(ctx).Warningf(ctx, message, fmtArgs...)
}

func (d *dynamicLogger) Errorf(ctx context.Context, message string, fmtArgs ...interface{}) {
	d.logger(ctx).Errorf(ctx, message, fmtArgs...)
}

func (d *dynamicLogger) Fatalf(ctx context.Context, message string, fmtArgs ...interface{}) {
	d.logger(ctx).Fatalf(ctx, message, fmtArgs...)
}

func (d *dynamicLogger) DPanicf(ctx context.Context, message string, fmtArgs ...interface{}) {
	d.logger(ctx).DPanicf(ctx, message, fmtArgs...)
}

func (d *dynamicLogger) Panicf(ctx context.Context, message string, fmtArgs ...interface{}) {
	d.logger(ctx).Panicf(ctx, message, fmtArgs...)
}

func (d *dynamicLogger) Debugz(ctx context.Context, message string, fields ...zap.Field) {
	d.logger(ctx).Debugz(ctx, message, fields...)
}

func (d *dynamicLogger) Infoz(ctx context.Context, message string, fields ...zap.Field) {
	d.logger(ctx).Infoz(ctx, message, fields...)
}

func (d *dynamicLogger) Warnz(ctx context.Context, message string, fields ...zap.Field) {
	d.logger(ctx).Warnz(ctx, message, fields...)
}

func (d *dynamicLogger) Warningz(ctx context.Context, message string, fields ...zap.Field) {
	d.logger(ctx).Warningz(ctx, message, fields...)
}

func (d *dynamicLogger) Errorz(ctx context.Context, message string, fields ...zap.Field) {
	d.logger(ctx).Errorz(ctx, message, fields...)
}

func (d *dynamicLogger) Fatalz(ctx context.Context, message string, fields ...zap.Field) {
	d.logger(ctx).Fatalz(ctx, message, fields...)
}

func (d *dynamicLogger) DPanicz(ctx context.Context, message string, fields ...zap.Field) {
	d.logger(ctx).DPanicz(ctx, message, fields...)
}

func (d *dynamicLogger) Panicz(ctx context.Context, message string, fields ...zap.Field) {
	d.logger(ctx).Panicz(ctx, message, fields...)
}

func (d *dynamicLogger) Debug(ctx context.Context, message string, context logger.Fields) {
	d.logger(ctx).Debug(ctx, message, context)
}

func (d *dynamicLogger) Info(ctx context.Context, message string, context logger.Fields) {
	d.logger(ctx).Info(ctx, message, context)
}

func (d *dynamicLogger) Warn(ctx context.Context, message string, context logger.Fields) {
	d.logger(ctx).Warn(ctx, message, context)
}

func (d *dynamicLogger) Warning(ctx context.Context, message string, context logger.Fields) {
	d.logger(ctx).Warning(ctx, message, context)
}

func (d *dynamicLogger) Error(ctx context.Context, message string, context logger.Fields) {
	d.logger(ctx).Error(ctx, message, context)
}

func (d *dynamicLogger) Fatal(ctx context.Context, message string, context logger.Fields) {
	d.logger(ctx).Fatal(ctx, message, context)
}

func (d *dynamicLogger) DPanic(ctx context.Context, message string, context logger.Fields) {
	d.logger(ctx).DPanic(ctx, message, context)
}

func (d *dynamicLogger) Panic(ctx context.Context, message string, context logger.Fields) {
	d.logger(ctx).Panic(ctx, message, context)
}

func (d *dynamicLogger) Debugw(ctx context.Context, message string, keysAndValues ...interface{}) {
	d.logger(ctx).Debugw(ctx, message, keysAndValues...)
}

func (d *dynamicLogger) Infow(ctx context.Context, message string, keysAndValues ...interface{}) {
	d.logger(ctx).Infow(ctx, message, keysAndValues...)
}

func (d *dynamicLogger) Warnw(ctx context.Context, message string, keysAndValues ...interface{}) {
	d.logger(ctx).Warnw(ctx, message, keysAndValues...)
}

func (d *dynamicLogger) Warningw(ctx context.Context, message string, keysAndValues ...interface{}) {
	d.logger(ctx).Warningw(ctx, message, keysAndValues...)
}

func (d *dynamicLogger) Errorw(ctx context.Context, message string, keysAndValues ...interface{}) {
	d.logger(ctx).Errorw(ctx, message, keysAndValues...)
}

func (d *dynamicLogger) Fatalw(ctx context.Context, message string, keysAndValues ...interface{}) {
	d.logger(ctx).Fatalw(ctx, message, keysAndValues...)
}

func (d *dynamicLogger) DPanicw(ctx context.Context, message string, keysAndValues ...interface{}) {
	d.logger(ctx).DPanicw(ctx, message, keysAndValues...)
}

func (d *dynamicLogger) Panicw(ctx context.Context, message string, keysAndValues ...interface{}) {
	d.logger(ctx).Panicw(ctx, message, keysAndValues...)
}
