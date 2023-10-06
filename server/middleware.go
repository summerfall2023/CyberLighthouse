package server

import (
	//"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	requestCount         int
	lastRequestTime      time.Time
	maxRequestsPerSecond = 5 // 允许的最大请求数
	lock                 sync.Mutex
)

// todo:change 监听对象
func CheckFrequency() {
	r := gin.Default()
	r.Use(RateLimitMiddleware)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "请求成功，总请求数：%d", requestCount)
	})
	r.Run(":8080")
}

func RateLimitMiddleware(c *gin.Context) {
	lock.Lock()
	defer lock.Unlock()

	now := time.Now()
	if now.Sub(lastRequestTime).Seconds() < 1.0/float64(maxRequestsPerSecond) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求频率过高"})
		c.Abort()
		return
	}

	requestCount++
	lastRequestTime = now
	c.Next()
}
