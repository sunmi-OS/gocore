package api

import (
	"net/http"

	"github.com/sunmi-OS/gocore/v2/utils"

	"github.com/sunmi-OS/gocore/v2/api/ecode"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type Context struct {
	*gin.Context
	C context.Context
	R Response
	T *utils.TraceHeader
}

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
