package redis

import (
	"context"
	"errors"
	"github.com/nova2018/easygin/logger"
	goLogger "github.com/nova2018/gologger/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net"
	"time"
)

var _ redis.Hook = &loggerHook{}

func newLoggerHook() redis.Hook {
	return &loggerHook{
		lg:      logger.Dynamic().WithOptions(zap.AddCallerSkip(1 + 4)),
		lgTrace: logger.Dynamic().WithOptions(zap.AddCallerSkip(1 + 5)),
	}
}

type loggerHook struct {
	lg      goLogger.Logger
	lgTrace goLogger.Logger
}

func (l loggerHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		l.lg.Debugf(ctx, "redis dial begin!")
		begin := time.Now()
		conn, err := next(ctx, network, addr)
		l.lg.Debugf(ctx, "redis dial end! timeUsed=%v", time.Now().Sub(begin))
		return conn, err
	}
}

func (l loggerHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		l.lg.Debugf(ctx, "redis begin!")
		begin := time.Now()
		err := next(ctx, cmd)
		l.lg.Debugf(ctx, "redis end!")
		_ = l.redisTrace(ctx, begin, false, cmd)
		return err
	}
}

func (l loggerHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		l.lg.Debugf(ctx, "redis pipeline begin!")
		begin := time.Now()
		err := next(ctx, cmds)
		l.lg.Debugf(ctx, "redis pipeline end!")
		_ = l.redisTrace(ctx, begin, false, cmds...)
		return err
	}
}

func (l loggerHook) redisTrace(ctx context.Context, begin time.Time, isPipeline bool, cmds ...redis.Cmder) error {
	elapsed := time.Since(begin)

	format := "redis cmd=[%s] elapsed=[%.3fms]"
	errFormat := "redis cmd=[%s] elapsed=[%.3fms] err=[%v]"
	if isPipeline {
		format = "redis pipeline cmd=[%s] elapsed=[%.3fms]"
		errFormat = "redis pipeline cmd=[%s] elapsed=[%.3fms] err=[%v]"
	}

	for _, c := range cmds {
		e := c.Err()
		strCmd := c.String()
		if e != nil && !errors.Is(e, redis.Nil) {
			l.lgTrace.Warnf(ctx, errFormat, strCmd, float64(elapsed.Nanoseconds())/1e6, e)
		} else {
			l.lgTrace.Debugf(ctx, format, strCmd, float64(elapsed.Nanoseconds())/1e6)
		}
	}

	return nil
}
