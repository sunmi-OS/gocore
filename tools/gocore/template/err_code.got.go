// Code generated by hero.
// DO NOT EDIT!
package template

import "bytes"

func FromErrCode(buffer *bytes.Buffer) {
	buffer.WriteString(`
package errcode

import (
	"github.com/sunmi-OS/gocore/v2/api/ecode"
)

var (
	ErrorNotFound = ecode.NewV2(50001, "record not found")
)`)

}
