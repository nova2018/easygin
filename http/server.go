package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jpillora/overseer"
	"github.com/nova2018/easygin/config"
	"github.com/nova2018/goutils"
	"os"
	"sync"
)

type Server struct {
	g   *gin.Engine
	h   goutils.RecoveryHandle
	o   sync.Once
	m   []*middlewarePool
	sig os.Signal
}

type Router func(router *gin.Engine)

func NewServer() *Server {
	s := &Server{}
	r := gin.New()
	s.g = r
	s.m = make([]*middlewarePool, 0, 5)
	initDefaultMiddleware(s)
	return s
}

func (s *Server) SetPanicHandler(h goutils.RecoveryHandle) {
	s.h = h
}

func (s *Server) GetPanicHandler() goutils.RecoveryHandle {
	return s.h
}

func (s *Server) RegisterMiddleware(key string, h interface{}, isFront ...bool) {
	if len(isFront) == 0 {
		isFront = []bool{false}
	}
	p := s.GetMiddlewarePool(key)
	p.Put(h, isFront[0])
}

func (s *Server) GetMiddlewarePool(key string) *middlewarePool {
	for _, p := range s.m {
		if p.key == key {
			return p
		}
	}
	p := newMiddlewarePool(key)
	s.m = append(s.m, p)
	return p
}

func (s *Server) bindMiddleware() {
	for _, p := range s.m {
		p.Register2Router(s.g)
	}
}

func (s *Server) SetRestartSignal(signal os.Signal) {
	s.sig = signal
}

func (s *Server) initRouter() {
	r := s.g
	s.bindMiddleware()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

func (s *Server) RegisterRouter(router ...Router) {
	s.o.Do(s.initRouter)

	for _, r := range router {
		r(s.g)
	}
}

func (s *Server) Run() {
	addr := config.GetString("app_addr")
	debug := config.GetBool("app_debug")
	if debug {
		_ = s.g.Run(addr)
	} else {
		overseer.Run(overseer.Config{
			Program: func(state overseer.State) {
				if state.Listener != nil {
					err := s.g.RunListener(state.Listener)
					if err != nil {
						fmt.Printf("RunListener Failure! err=[%v]\n", err)
					}
				} else {
					panic("overseer start failure!")
				}
			},
			Address:       addr,
			RestartSignal: s.sig,
		})
	}

}

func (s *Server) RunSimple() {
	addr := config.GetString("app_addr")
	err := s.g.Run(addr)
	if err != nil {
		fmt.Printf("RunListener Failure! err=[%v]\n", err)
	}
}
