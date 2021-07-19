// Code generated by hero.
// source: /Users/SM0286/code/core/gocore/tools/gocore/template/api.got
// DO NOT EDIT!
package template

import (
	"bytes"

	"github.com/shiyanhui/hero"
)

func FromApi(name, handler string, functions []string, req []string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package api

import (
	"`)
	hero.EscapeHTML(name, buffer)
	buffer.WriteString(`/app/domain"
	"`)
	hero.EscapeHTML(name, buffer)
	buffer.WriteString(`/app/errcode"
	"`)
	hero.EscapeHTML(name, buffer)
	buffer.WriteString(`/pkg/parse"
	"`)
	hero.EscapeHTML(name, buffer)
	buffer.WriteString(`/app/def"

	"github.com/labstack/echo/v4"
)

var `)
	hero.EscapeHTML(handler, buffer)
	buffer.WriteString(`Handler = ` + "`" + ` + handler + ` + "`" + `{}
type `)
	hero.EscapeHTML(handler, buffer)
	buffer.WriteString(` struct{}

`)
	for k1, v1 := range functions {
		buffer.WriteString(`
    // `)
		hero.EscapeHTML(v1, buffer)
		buffer.WriteString(`
    func (*`)
		hero.EscapeHTML(handler, buffer)
		buffer.WriteString(`) `)
		hero.EscapeHTML(v1, buffer)
		buffer.WriteString(`(c echo.Context) error {
        params := new(def.`)
		hero.EscapeHTML(req[k1], buffer)
		buffer.WriteString(`)
        //参数验证绑定
        _, response, err := parse.ParseJson(c, params)
        if err != nil {
            return response.RetError(err, errcode.Code0002)
        }
        resp, code, err := domain.`)
		hero.EscapeHTML(handler, buffer)
		buffer.WriteString(`Handler.`)
		hero.EscapeHTML(v1, buffer)
		buffer.WriteString(`(params)
        if err != nil {
            return response.RetError(err, code)
        }
        return response.RetSuccess(resp)
    }

`)
	}

}