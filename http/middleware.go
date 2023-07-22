package http

import (
	"bytes"
	"container/list"
	"github.com/gin-gonic/gin"
	"github.com/nova2018/easygin/config"
	"github.com/nova2018/easygin/file"
	"github.com/nova2018/easygin/logger"
	goLogger "github.com/nova2018/gologger/logger"
	"github.com/nova2018/goutils"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	MiddlewareAccessLog       = "access_log"
	MiddlewareRecoveryHandler = "recovery_handler"
	MiddlewareTraceInfo       = "trace_info"
	MiddlewareRequestLog      = "request_log"
)

var (
	HeaderExclude = []string{
		"Accept",
		"Accept-Encoding",
		"Connection",
		"User-Agent",
		"Content-Type",
		"Content-Length",
	}

	ListTextMIME = []string{
		"text/html",
		"text/plain",
		"text/xml",
		"application/json",
	}
)

type middlewarePool struct {
	key  string
	list *list.List
}

func newMiddlewarePool(key string) *middlewarePool {
	return &middlewarePool{
		key:  key,
		list: list.New(),
	}
}

func (p *middlewarePool) List() *list.List {
	return p.list
}

func (p *middlewarePool) Put(h interface{}, isFront bool) {
	p.checkHandler(h)
	if isFront {
		p.list.PushFront(h)
	} else {
		p.list.PushBack(h)
	}
}

func (p *middlewarePool) checkHandler(h interface{}) {
	if _, ok := h.(gin.HandlerFunc); ok {
		return
	}
	if _, ok := h.(func(*gin.Engine)); ok {
		return
	}
	panic("middleware handler invalid")
}

func (p *middlewarePool) Register2Router(router *gin.Engine) {
	current := p.list.Front()
	for current != nil {
		if h, ok := current.Value.(gin.HandlerFunc); ok {
			router.Use(h)
		} else if h, ok := current.Value.(func(*gin.Engine)); ok {
			h(router)
		}
		current = current.Next()
	}
}

func initDefaultMiddleware(server *Server) {
	lg := AccessLogger()
	if lg != nil {
		server.RegisterMiddleware(MiddlewareAccessLog, lg)
	}
	server.RegisterMiddleware(MiddlewareRecoveryHandler, RecoveryWithPanicHandler(func() goutils.RecoveryHandle {
		return server.GetPanicHandler()
	}))
	server.RegisterMiddleware(MiddlewareTraceInfo, TraceInfo())
	server.RegisterMiddleware(MiddlewareRequestLog, RequestLogger())
}

func AccessLogger() gin.HandlerFunc {
	writer := getAccessLogWriter()
	if writer == nil {
		return nil
	}
	return gin.LoggerWithWriter(writer)
}

func getAccessLogWriter() io.Writer {
	var cfg *file.WriteConfig
	access := config.Sub("access")
	if access == nil {
		return nil
	}
	err := access.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	if cfg == nil {
		return nil
	}

	return file.NewWriter(cfg)
}

func RequestLogger() gin.HandlerFunc {
	return _middleware.requestLogger
}

func TraceInfo() gin.HandlerFunc {
	return _middleware.putTraceInfo
}

func Recovery() gin.HandlerFunc {
	return RecoveryWithPanicHandler(nil)
}

func RecoveryWithPanicHandler(fn func() goutils.RecoveryHandle) gin.HandlerFunc {
	writer := NewRecoveryWriter()
	multi := io.MultiWriter(gin.DefaultErrorWriter, writer)
	return gin.CustomRecoveryWithWriter(multi, _middleware.RecoveryHandler(writer, fn))
}

type middleware struct {
}

var _middleware = &middleware{}

func (middleware) putTraceInfo(ctx *gin.Context) {
	logInfo := goLogger.GetInfoPool().Get()
	defer logInfo.Free()
	ctx.Set(goLogger.ContextField, logInfo)

	ctx.Next()
}

func (middleware) requestLogger(ctx *gin.Context) {
	s1 := time.Now()
	writer := AcquireWriter(ctx.Writer)
	ctx.Writer = writer
	defer func() {
		if ctx != nil {
			ctx.Writer = writer.ResponseWriter
		}
		writer.Free()
	}()

	body, err := ctx.GetRawData()
	if err == nil && ctx.Request.Body != nil {
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	listCookies := ctx.Request.Cookies()
	mapCookie := make(map[string]string, len(listCookies))
	for _, c := range listCookies {
		mapCookie[c.Name] = c.Value
	}
	header := ctx.Request.Header
	mapHeader := make(map[string]string, len(header))
	mapHeaderLink := make(map[string]string, len(header))
	for k := range header {
		mapHeader[k] = header.Get(k)
		mapHeaderLink[strings.ToLower(k)] = k
	}
	for _, k := range HeaderExclude {
		if key, ok := mapHeaderLink[strings.ToLower(k)]; ok {
			delete(mapHeader, key)
		}
	}

	// 请求日志
	logger.Infow(ctx, "REQUEST:BaseInfo",
		"method", ctx.Request.Method,
		"query", ctx.Request.URL.Query(),
		"header", mapHeader,
		"body", string(body),
		"cookie", mapCookie,
		"serverIp", ctx.Request.Host,
		"domain", ctx.Request.Host,
		"uri", ctx.Request.URL.RequestURI(),
		"path", ctx.Request.URL.Path,
		"userAgent", ctx.Request.UserAgent(),
		"ip", ctx.ClientIP(),
		"remote_addr", ctx.Request.RemoteAddr,
	)

	// Process request
	ctx.Next()

	timeUsed := time.Now().Sub(s1)

	resp := "<binary>"
	if isTextResponse(ctx) {
		resp = writer.body.String()
	}

	// 响应日志
	logger.Infow(ctx, "ACCESS_LOG",
		"uri", ctx.Request.URL.RequestURI(),
		"path", ctx.Request.URL.Path,
		"status", ctx.Writer.Status(),
		"size", ctx.Writer.Size(),
		"redirect", ctx.Writer.Header().Get("Location"),
		"content-type", ctx.Writer.Header().Get("Content-Type"),
		"response", resp,
		"cost", timeUsed.Truncate(time.Millisecond),
	)
}

/**
 * isTextResponse
 * @Description: 检查返回值类型是否是文本格式
 * @author lijunpeng<lijunpeng@weimiao.cn>
 * @date: 2022-07-29 15:18:26
 * @param ctx
 * @return bool
 */
func isTextResponse(ctx *gin.Context) bool {
	contentType := ctx.Writer.Header().Get("Content-Type")
	realType := strings.SplitN(contentType, ";", 2)
	if len(realType) == 0 {
		return true
	}
	for _, x := range ListTextMIME {
		if x == realType[0] {
			return true
		}
	}
	return false
}

func (middleware) RecoveryHandler(writer *recoveryWriter, fn func() goutils.RecoveryHandle) gin.RecoveryFunc {
	return func(ctx *gin.Context, recovered interface{}) {
		result := writer.ReadAndReset()
		if err, ok := recovered.(error); ok {
			//event.Dispatch(alarm.NewSimpleEvent(
			//	"Panic! Recovery",
			//	fmt.Sprintf("error: %s", err.Error()),
			//	result,
			//	"",
			//	utils.GetTraceId(ctx),
			//	"default",
			//))

			if fn != nil {
				if h := fn(); h != nil {
					h(ctx, err, []byte(result))
				}
			}
		}

		if !ctx.IsAborted() {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
