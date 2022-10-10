package template

import "bytes"

func FromConfLocal(env, content string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package conf

var `)
	buffer.WriteString(env)
	buffer.WriteString(` = ` + "`" + ``)
	buffer.WriteString(content)
	buffer.WriteString(`
` + "`" + `

`)

}
