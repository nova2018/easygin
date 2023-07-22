package app

import (
	"github.com/nova2018/easygin/http"
	"github.com/nova2018/easygin/logger"
	"github.com/nova2018/goutils"
	"github.com/spf13/pflag"
)

type Server struct {
	BaseApplication

	s *http.Server
}

func NewServer() *Server {
	return NewServerWithFlags(nil)
}

func NewServerWithFlags(flag *pflag.FlagSet) *Server {
	s := &Server{}
	s.init(flag)
	s.s = http.NewServer()
	return s
}

func (s *Server) Router(r ...http.Router) {
	s.s.RegisterRouter(r...)
}

func (s *Server) Middleware(key string, h interface{}, isFront ...bool) {
	s.s.RegisterMiddleware(key, h, isFront...)
}

func (s *Server) init(flag *pflag.FlagSet) {
	SetApplication(s)

	s.BaseApplication.init()
	s.parse(flag)

	s.startUp()
}

func (s *Server) Run() {
	s.RunWithPanic(nil)
}

func (s *Server) GetServer() *http.Server {
	return s.s
}

func (s *Server) RunWithPanic(panicHandler goutils.RecoveryHandle) {
	defer logger.Sync()

	s.s.SetPanicHandler(panicHandler)
	s.s.Run()

	// 等待异步事件处理完成
	//event.GetDispatcher().WaitAsync()
}
