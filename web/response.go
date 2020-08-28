package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/sunmi-OS/gocore/ecode"
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
	switch c.(type) {
	case *gin.Context:
		c.(*gin.Context).JSON(http.StatusOK, rsp)
	case echo.Context:
		_ = c.(echo.Context).JSON(http.StatusOK, rsp)
	default:
	}
}

// Redirect c: gin or echo Context
func Redirect(c interface{}, location string) {
	switch c.(type) {
	case *gin.Context:
		c.(*gin.Context).Redirect(http.StatusFound, location)
	case echo.Context:
		_ = c.(echo.Context).Redirect(http.StatusFound, location)
	default:
	}
}

// File c: gin or echo Context
func File(c interface{}, filePath, fileName string) {
	switch c.(type) {
	case *gin.Context:
		c.(*gin.Context).Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		c.(*gin.Context).File(filePath)
	case echo.Context:
		c.(echo.Context).Response().Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		_ = c.(echo.Context).File(filePath)
	default:
	}
}
