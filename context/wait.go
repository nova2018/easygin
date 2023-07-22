package context

import (
	"context"
	"github.com/gin-gonic/gin"
	"sync/atomic"
)

func WaitContextDone(ctx context.Context, fn func()) {
	// context done check
	var flg int32 = 0
	closeFn := func() {
		if atomic.AddInt32(&flg, 1) == 1 {
			fn()
		}
	}

	{
		done := ctx.Done()
		if c, ok := ctx.(*gin.Context); ok {
			done = c.Request.Context().Done()
		}

		if done != nil {
			go func() {
				<-done
				closeFn()
			}()
		}
	}
}
