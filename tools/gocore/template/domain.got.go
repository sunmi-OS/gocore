package template

import "bytes"

func FromDomain(buffer *bytes.Buffer) {
	buffer.WriteString(`
package biz
`)

}
