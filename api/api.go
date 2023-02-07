package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sunmi-OS/gocore/v2/api/ecode"
	"github.com/sunmi-OS/gocore/v2/utils"
)

type Context struct {
	*gin.Context
	C context.Context
	R Response
	T *utils.TraceHeader
}

var (
	ErrorBind      = errors.New("missing required parameters")
	TraceHeaderKey struct{}
)

//const TraceHeaderKey = "TraceHeaderKey"

// NewContext 初始化上下文包含context.Context
// 对链路信息进行判断并且在Response时返回TraceId信息
func NewContext(g *gin.Context) Context {
	c := Context{
		Context: g,
		C:       context.Background(),
		R:       NewResponse(),
	}

	if g.GetHeader(utils.XB3TraceId) != "" {
		g.Header(utils.XB3TraceId, g.GetHeader(utils.XB3TraceId))
		c.T = utils.SetHttp(g.Request.Header)
		//g.Set(TraceHeaderKey, c.T)
		c.C = context.WithValue(c.C, TraceHeaderKey, c.T)
	}
	return c
}

// Success 返回正常数据
func (c *Context) Success(data interface{}) {
	c.R.Data = data
	c.JSON(http.StatusOK, c.R)
}

// Error 返回异常信息，自动识别Code码
func (c *Context) Error(err error) {
	c.R.Code = ecode.Transform(err)
	c.R.Msg = err.Error()
	c.JSON(http.StatusOK, c.R)
}

// RetJSON 针对 ecode v2
func (c *Context) RetJSON(data interface{}, err error) {
	e := ecode.FromError(err)
	c.R.Code = e.Code()
	c.R.Data = data
	c.R.Msg = e.Message()
	c.JSON(http.StatusOK, c.R)
}

// ErrorCodeMsg 直接指定code和msg
func (c *Context) ErrorCodeMsg(code int, msg string) {
	c.R.Code = code
	c.R.Msg = msg
	c.JSON(http.StatusOK, c.R)
}

// Response 直接指定code和msg和data
func (c *Context) Response(code int, msg string, data interface{}) {
	c.R.Code = code
	c.R.Msg = msg
	c.R.Data = data
	c.JSON(http.StatusOK, c.R)
}

// BindValidator 参数绑定结构体，并且按照tag进行校验返回校验结果
func (c *Context) BindValidator(obj interface{}) error {
	err := c.ShouldBind(obj)
	if err != nil {
		if utils.IsRelease() {
			return ErrorBind
		}
		return err
	}
	return nil
}
