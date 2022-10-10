package template

import "bytes"

func FromCronJob(cron, comment string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package cronjob
// `)
	buffer.WriteString(cron)
	buffer.WriteString(" " + comment)
	buffer.WriteString(`
func `)
	buffer.WriteString(cron)
	buffer.WriteString(`() {
}`)

}
