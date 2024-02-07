package template

import "bytes"

func FromApiRequest(request, params string, buffer *bytes.Buffer) {
	buffer.WriteString(`

type `)
	buffer.WriteString(request)
	buffer.WriteString(` struct {
    `)
	buffer.WriteString(params)
	buffer.WriteString(`
}`)

}
