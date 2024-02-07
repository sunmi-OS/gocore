package template

import "bytes"

func FromApiRoutes(name, routes string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package route

import (
	"`)
	buffer.WriteString(name)
	buffer.WriteString(`/api"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
    `)
	buffer.WriteString(routes)
	buffer.WriteString(`
}
`)

}
