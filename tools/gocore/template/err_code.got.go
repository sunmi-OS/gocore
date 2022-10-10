package template

import "bytes"

func FromErrCode(buffer *bytes.Buffer) {
	buffer.WriteString(`
package errcode

import (
	"github.com/sunmi-OS/gocore/v2/api/ecode"
	"gorm.io/gorm"
)

var (
	ErrorNotFound = ecode.New(50001, gorm.ErrRecordNotFound)
)`)

}
