package template

import "bytes"

func FromDomainHandler(handlers []string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package biz
`)
	for _, v1 := range handlers {
		buffer.WriteString(`
    var `)
		buffer.WriteString(v1)
		buffer.WriteString(`Handler = &`)
		buffer.WriteString(v1)
		buffer.WriteString(`{}
    type `)
		buffer.WriteString(v1)
		buffer.WriteString(` struct{}
`)
	}

}
