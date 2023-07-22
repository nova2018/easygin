package app

import (
	gContext "context"
	"github.com/nova2018/easygin/context"
	"github.com/nova2018/easygin/logger"
	"github.com/nova2018/goutils"
	"github.com/spf13/pflag"
	"sync"
)

type Cli struct {
	BaseApplication

	h   []Handler
	ctx gContext.Context
}

func NewCli() *Cli {
	return NewCliWithFlag(nil)
}

func NewCliWithFlag(f *pflag.FlagSet) *Cli {
	c := &Cli{}
	c.init(f)
	return c
}

type Handler func(ctx gContext.Context)

func (s *Cli) AttachHandler(h ...Handler) {
	if s.h == nil {
		s.h = h
	} else {
		s.h = append(s.h, h...)
	}
}

func (s *Cli) init(f *pflag.FlagSet) {
	SetApplication(s)

	s.BaseApplication.init()
	s.parse(f)
	s.startUp()
}

func (s *Cli) Context(ctx gContext.Context) gContext.Context {
	if ctx != nil {
		s.ctx = ctx
	}
	return s.ctx
}

func (s *Cli) RunWithPanic(recoveryHandle goutils.RecoveryHandle) {
	defer logger.Sync()

	ctx := s.Context(nil)
	if ctx == nil {
		ctx = context.NewContext()
	} else {
		ctx = context.NewContextWithParent(ctx)
	}
	if recoveryHandle != nil {
		ctx = gContext.WithValue(ctx, goutils.ContextRecoveryField, recoveryHandle)
	}

	wg := sync.WaitGroup{}
	for _, h := range s.h {
		wg.Add(1)
		done := goutils.GoWithContextHandler(ctx, h, recoveryHandle)
		go func(done <-chan struct{}) {
			<-done
			wg.Done()
		}(done)
	}
	wg.Wait()

	// 等待异步事件处理完成
	//event.GetDispatcher().WaitAsync()
}

func (s *Cli) Run() {
	s.RunWithPanic(nil)
}
