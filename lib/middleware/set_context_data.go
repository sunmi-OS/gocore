package middleware

import "github.com/labstack/echo/v4"

func SetContextData(key string, value interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(key, value)
			return next(c)
		}
	}
}
