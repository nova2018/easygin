package context

import (
	"context"
	"github.com/nova2018/gologger/logger"
)

func NewContext() context.Context {
	return NewContextWithParent(context.TODO())
}

func NewContextWithParent(ctx context.Context) context.Context {
	logInfo := logger.GetInfoPool().Get()
	WaitContextDone(ctx, func() {
		logInfo.AutoFree()
	})
	return context.WithValue(ctx, logger.ContextField, logInfo)
}
