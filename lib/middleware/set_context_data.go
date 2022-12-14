package middleware

import "github.com/gin-gonic/gin"

//SetContextData is used to store a new key/value pair exclusively for this context.
func SetContextData(key string, value interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, value)
		c.Next()
	}
}
