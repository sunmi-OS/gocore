package template

import "bytes"

func FromJob(job, comment string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package job
// `)
	buffer.WriteString(job)
	buffer.WriteString(" " + comment)
	buffer.WriteString(`
func `)
	buffer.WriteString(job)
	buffer.WriteString(`() {
}
`)

}
