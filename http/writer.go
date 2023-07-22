package http

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/nova2018/goutils"
	"sync"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
	p    goutils.Pool[*responseWriter]
}

func (w *responseWriter) String() string {
	return w.body.String()
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (w *responseWriter) Reset() {
	w.body.Reset()
}

func (w *responseWriter) Free() {
	w.p.Put(w)
}

func (w *responseWriter) SetPool(p goutils.Pool[*responseWriter]) {
	w.p = p
}

var (
	writePool = goutils.NewPool[*responseWriter](func() *responseWriter {
		return &responseWriter{
			body: &bytes.Buffer{},
		}
	})
)

func AcquireWriter(writer gin.ResponseWriter) *responseWriter {
	w := writePool.Get()
	w.ResponseWriter = writer
	return w
}

type recoveryWriter struct {
	lock *sync.Mutex
	bf   *bytes.Buffer
}

func (r *recoveryWriter) Write(b []byte) (int, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.bf.Write(b)
}

func (r *recoveryWriter) ReadAndReset() string {
	r.lock.Lock()
	defer r.lock.Unlock()
	s := r.bf.String()
	r.bf.Reset()
	return s
}

func NewRecoveryWriter() *recoveryWriter {
	return &recoveryWriter{
		lock: &sync.Mutex{},
		bf:   &bytes.Buffer{},
	}
}
