package logger

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nova2018/gologger/logger"
	"go.uber.org/zap"
)

const (
	CtxKeyLogger = "easygin.logger.ctxLogger"
)

func WithLogger(ctx context.Context, opts ...zap.Option) context.Context {
	return WithCustomLogger(ctx, loggerDefault, opts...)
}

func WithCustomLogger(ctx context.Context, name string, opts ...zap.Option) context.Context {
	lg := FactoryWithLogger(ctx, Logger(name), opts...)
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set(CtxKeyLogger, lg)
		return ginCtx
	}
	return context.WithValue(ctx, CtxKeyLogger, lg)
}

func getLoggerWithContext(ctx context.Context) logger.Logger {
	v := ctx.Value(CtxKeyLogger)
	if i, ok := v.(*ctxLogger); ok {
		return i.logger
	}
	return nil
}
