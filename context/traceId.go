package context

import (
	"context"
	"github.com/nova2018/gologger/logger"
)

func GetTraceInfo(ctx context.Context) *logger.LogInfo {
	type getter interface {
		Get(key string) (value interface{}, exists bool)
	}
	key := logger.ContextField
	var info interface{}
	if s, ok := ctx.(getter); ok {
		if v, exists := s.Get(key); exists {
			info = v
		}
	}
	if info == nil {
		info = ctx.Value(key)
	}
	if v, ok := info.(*logger.LogInfo); ok {
		return v
	}
	return nil
}

func GetTraceId(ctx context.Context) string {
	traceInfo := GetTraceInfo(ctx)
	if traceInfo != nil {
		return traceInfo.TraceId
	}
	return ""
}
