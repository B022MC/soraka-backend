package middleware

import (
	"github.com/gin-gonic/gin"
	"go-utils/utils/ecode"
	"go-utils/utils/lang/browserx"
	"go-utils/utils/response"

	"sync"
	"time"
)

type debounceItem struct {
	lastTime time.Time
	mutex    sync.Mutex
}

var debounceStore = make(map[string]*debounceItem)
var storeMutex sync.Mutex

func Debounce(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		storeMutex.Lock()
		key := browserx.GenerateBrowserKey(c)
		item, exists := debounceStore[key]
		if !exists {
			item = &debounceItem{}
			debounceStore[key] = item
		}
		storeMutex.Unlock()

		item.mutex.Lock()
		defer item.mutex.Unlock()

		now := time.Now()
		if now.Sub(item.lastTime) < duration {
			c.Abort()
			response.Fail(c, ecode.RateLimitAllowErrFailed, ecode.Text(ecode.RateLimitAllowErrFailed))
			return
		}

		item.lastTime = now
		c.Next()
	}
}
