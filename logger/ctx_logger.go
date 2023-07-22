package logger

import (
	"context"
	cctx "github.com/nova2018/easygin/context"
	"github.com/nova2018/gologger/logger"
	"github.com/nova2018/goutils"
	"go.uber.org/zap"
)

type ContextLogger interface {
	Named(string) ContextLogger

	WithOptions(opts ...zap.Option) ContextLogger

	Debugf(message string, fmtArgs ...interface{})

	Infof(message string, fmtArgs ...interface{})

	Warnf(message string, fmtArgs ...interface{})

	Warningf(message string, fmtArgs ...interface{})

	Errorf(message string, fmtArgs ...interface{})

	Fatalf(message string, fmtArgs ...interface{})

	DPanicf(message string, fmtArgs ...interface{})

	Panicf(message string, fmtArgs ...interface{})

	Debugz(message string, fields ...zap.Field)

	Infoz(message string, fields ...zap.Field)

	Warnz(message string, fields ...zap.Field)

	Warningz(message string, fields ...zap.Field)

	Errorz(message string, fields ...zap.Field)

	Fatalz(message string, fields ...zap.Field)

	DPanicz(message string, fields ...zap.Field)

	Panicz(message string, fields ...zap.Field)

	Debug(message string, context logger.Fields)

	Info(message string, context logger.Fields)

	Warn(message string, context logger.Fields)

	Warning(message string, context logger.Fields)

	Error(message string, context logger.Fields)

	Fatal(message string, context logger.Fields)

	DPanic(message string, context logger.Fields)

	Panic(message string, context logger.Fields)

	Debugw(message string, keysAndValues ...interface{})

	Infow(message string, keysAndValues ...interface{})

	Warnw(message string, keysAndValues ...interface{})

	Warningw(message string, keysAndValues ...interface{})

	Errorw(message string, keysAndValues ...interface{})

	Fatalw(message string, keysAndValues ...interface{})

	DPanicw(message string, keysAndValues ...interface{})

	Panicw(message string, keysAndValues ...interface{})
}

type ctxLogger struct {
	ctx    context.Context
	logger logger.Logger
	p      goutils.Pool[*ctxLogger]
}

func (c *ctxLogger) Named(name string) ContextLogger {
	clone := c.clone()
	clone.logger = c.logger.Named(name)
	return clone
}

func (c *ctxLogger) WithOptions(opts ...zap.Option) ContextLogger {
	clone := c.clone()
	clone.logger = c.logger.WithOptions(opts...)
	return clone
}

func (c *ctxLogger) clone() *ctxLogger {
	x := getCtxLogger()
	x.logger = c.logger
	x.ctx = c.ctx
	x.p = c.p
	x.WaitRelease()
	return x
}

func (c *ctxLogger) Reset() {
	c.logger = nil
	c.ctx = nil
	c.p = nil
}

func (c *ctxLogger) SetPool(p goutils.Pool[*ctxLogger]) {
	c.p = p
}

func (c *ctxLogger) Release() {
	if c.p != nil {
		c.p.Put(c)
	}
}

func (c *ctxLogger) WaitRelease() {
	cctx.WaitContextDone(c.ctx, func() {
		c.Release()
	})
}

func (c *ctxLogger) Debugf(message string, fmtArgs ...interface{}) {
	c.logger.Debugf(c.ctx, message, fmtArgs...)
}

func (c *ctxLogger) Infof(message string, fmtArgs ...interface{}) {
	c.logger.Infof(c.ctx, message, fmtArgs...)
}

func (c *ctxLogger) Warnf(message string, fmtArgs ...interface{}) {
	c.logger.Warnf(c.ctx, message, fmtArgs...)
}

func (c *ctxLogger) Warningf(message string, fmtArgs ...interface{}) {
	c.logger.Warningf(c.ctx, message, fmtArgs...)
}

func (c *ctxLogger) Errorf(message string, fmtArgs ...interface{}) {
	c.logger.Errorf(c.ctx, message, fmtArgs...)
}

func (c *ctxLogger) Fatalf(message string, fmtArgs ...interface{}) {
	c.logger.Fatalf(c.ctx, message, fmtArgs...)
}

func (c *ctxLogger) DPanicf(message string, fmtArgs ...interface{}) {
	c.logger.DPanicf(c.ctx, message, fmtArgs...)
}

func (c *ctxLogger) Panicf(message string, fmtArgs ...interface{}) {
	c.logger.Panicf(c.ctx, message, fmtArgs...)
}

func (c *ctxLogger) Debugz(message string, fields ...zap.Field) {
	c.logger.Debugz(c.ctx, message, fields...)
}

func (c *ctxLogger) Infoz(message string, fields ...zap.Field) {
	c.logger.Infoz(c.ctx, message, fields...)
}

func (c *ctxLogger) Warnz(message string, fields ...zap.Field) {
	c.logger.Warnz(c.ctx, message, fields...)
}

func (c *ctxLogger) Warningz(message string, fields ...zap.Field) {
	c.logger.Warningz(c.ctx, message, fields...)
}

func (c *ctxLogger) Errorz(message string, fields ...zap.Field) {
	c.logger.Errorz(c.ctx, message, fields...)
}

func (c *ctxLogger) Fatalz(message string, fields ...zap.Field) {
	c.logger.Fatalz(c.ctx, message, fields...)
}

func (c *ctxLogger) DPanicz(message string, fields ...zap.Field) {
	c.logger.DPanicz(c.ctx, message, fields...)
}

func (c *ctxLogger) Panicz(message string, fields ...zap.Field) {
	c.logger.Panicz(c.ctx, message, fields...)
}

func (c *ctxLogger) Debug(message string, context logger.Fields) {
	c.logger.Debug(c.ctx, message, context)
}

func (c *ctxLogger) Info(message string, context logger.Fields) {
	c.logger.Info(c.ctx, message, context)
}

func (c *ctxLogger) Warn(message string, context logger.Fields) {
	c.logger.Warn(c.ctx, message, context)
}

func (c *ctxLogger) Warning(message string, context logger.Fields) {
	c.logger.Warning(c.ctx, message, context)
}

func (c *ctxLogger) Error(message string, context logger.Fields) {
	c.logger.Error(c.ctx, message, context)
}

func (c *ctxLogger) Fatal(message string, context logger.Fields) {
	c.logger.Fatal(c.ctx, message, context)
}

func (c *ctxLogger) DPanic(message string, context logger.Fields) {
	c.logger.DPanic(c.ctx, message, context)
}

func (c *ctxLogger) Panic(message string, context logger.Fields) {
	c.logger.Panic(c.ctx, message, context)
}

func (c *ctxLogger) Debugw(message string, keysAndValues ...interface{}) {
	c.logger.Debugw(c.ctx, message, keysAndValues...)
}

func (c *ctxLogger) Infow(message string, keysAndValues ...interface{}) {
	c.logger.Infow(c.ctx, message, keysAndValues...)
}

func (c *ctxLogger) Warnw(message string, keysAndValues ...interface{}) {
	c.logger.Warnw(c.ctx, message, keysAndValues...)
}

func (c *ctxLogger) Warningw(message string, keysAndValues ...interface{}) {
	c.logger.Warningw(c.ctx, message, keysAndValues...)
}

func (c *ctxLogger) Errorw(message string, keysAndValues ...interface{}) {
	c.logger.Errorw(c.ctx, message, keysAndValues...)
}

func (c *ctxLogger) Fatalw(message string, keysAndValues ...interface{}) {
	c.logger.Fatalw(c.ctx, message, keysAndValues...)
}

func (c *ctxLogger) DPanicw(message string, keysAndValues ...interface{}) {
	c.logger.DPanicw(c.ctx, message, keysAndValues...)
}

func (c *ctxLogger) Panicw(message string, keysAndValues ...interface{}) {
	c.logger.Panicw(c.ctx, message, keysAndValues...)
}

var (
	_ctxLoggerPool = goutils.NewPool(func() interface{} {
		return &ctxLogger{}
	})
)

func Factory(ctx context.Context, opts ...zap.Option) ContextLogger {
	return FactoryWithLogger(ctx, Default(), opts...)
}

func getCtxLogger() *ctxLogger {
	return _ctxLoggerPool.Get().(*ctxLogger)
}

func FactoryWithLogger(ctx context.Context, logger logger.Logger, opts ...zap.Option) ContextLogger {
	cLogger := getCtxLogger()
	cLogger.ctx = ctx
	lg := logger.WithOptions(zap.AddCallerSkip(1))
	if len(opts) > 0 {
		lg = lg.WithOptions(opts...)
	}
	cLogger.logger = lg
	cLogger.WaitRelease()
	return cLogger
}
