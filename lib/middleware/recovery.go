package middleware

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

type RecoverInfo struct {
	Time       string      `json:"time"`
	RequestURI string      `json:"request_uri"`
	Err        interface{} `json:"error"`
	Stack      string      `json:"stack"`
}

// Recovery gin middleware recovery
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				const size = 64 << 10
				stack := make([]byte, size)
				stack = stack[:runtime.Stack(stack, false)]
				bs, _ := sonic.Marshal(RecoverInfo{
					Time:       time.Now().Format("2006-01-02 15:04:05"),
					RequestURI: c.Request.Host + c.Request.RequestURI,
					Err:        err,
					Stack:      string(stack),
				})
				log.Printf("[GinPanic] %s\n", string(bs))
				c.JSON(http.StatusOK, struct {
					Code int         `json:"code"`
					Data interface{} `json:"data"`
					Msg  string      `json:"msg"`
				}{
					Code: -1,
					Data: nil,
					Msg:  "server panic",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
