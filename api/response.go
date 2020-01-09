//	PhalGo-Response
//	返回json参数,默认结构code,data,msg
//	喵了个咪 <wenzhenxi@vip.qq.com> 2016/5/11
//  依赖情况:
//          "github.com/labstack/echo" 必须基于echo路由

package api

import (
	"github.com/labstack/echo"
	"net/http"
)

type Response struct {
	Context   echo.Context
	parameter *RetParameter
}

type RetParameter struct {
	Code int         `json:"code";xml:"code"`
	Data interface{} `json:"data";xml:"data"`
	Msg  string      `json:"msg";xml:"msg"`
}

const DefaultCode = 1

var HttpStatus = http.StatusOK

// 初始化Response
func NewResponse(c echo.Context) *Response {

	R := new(Response)
	R.Context = c
	R.parameter = new(RetParameter)
	R.parameter.Data = nil
	return R
}

// 设置返回的Status值默认http.StatusOK
func (this *Response) SetStatus(i int) {
	HttpStatus = i
}

func (this *Response) SetMsg(s string) {
	this.parameter.Msg = s
}

func (this *Response) SetData(d interface{}) {
	this.parameter.Data = d
}

func (this *Response) GetParameter() *RetParameter {
	return this.parameter
}

// 返回自定自定义的消息格式
func (this *Response) RetCustomize(code int, d interface{}, msg string) error {

	this.parameter.Code = code
	this.parameter.Data = d
	this.parameter.Msg = msg

	return this.Context.JSON(HttpStatus, this.parameter)
}

// 返回成功的结果 默认code为1
func (this *Response) RetSuccess(d interface{}) error {

	this.parameter.Code = DefaultCode
	this.parameter.Data = d

	return this.Context.JSON(HttpStatus, this.parameter)
}

// 返回失败结果
func (this *Response) RetError(e error, c int) error {

	this.parameter.Code = c
	this.parameter.Msg = e.Error()

	return this.Context.JSON(HttpStatus, this.parameter)
}

// 输出返回结果
func (this *Response) Write(b []byte) {

	_, e := this.Context.Response().Write(b)
	if e != nil {
		print(e.Error())
	}
}
