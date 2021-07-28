package web

import (
	"fmt"
	"net/http"

	"github.com/sunmi-OS/gocore/v2/api/ecode"

	"github.com/gin-gonic/gin"
)

const (
	TypeOctetStream = "application/octet-stream"
	TypeForm        = "application/x-www-form-urlencoded"
	TypeJson        = "application/json"
	TypeXml         = "application/xml"
	TypeJpg         = "image/jpeg"
	TypePng         = "image/png"
)

// JSON c: gin or echo Context
func JSON(c interface{}, data interface{}, err error) {
	e := ecode.AnalyseError(err)
	rsp := CommonRsp{
		Code:    e.Code(),
		Message: e.Message(),
		Data:    data,
	}

	c.(*gin.Context).JSON(http.StatusOK, rsp)

}

// Redirect c: gin or echo Context
func Redirect(c interface{}, location string) {

	c.(*gin.Context).Redirect(http.StatusFound, location)

}

// File c: gin or echo Context
func File(c interface{}, filePath, fileName string) {

	c.(*gin.Context).Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.(*gin.Context).File(filePath)

}
