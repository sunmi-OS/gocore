package template

import "bytes"

func FromConfBase(baseConf string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package conf

var BaseConfig = ` + "`" + `
[network]
ApiServiceHost = "`)
	buffer.WriteString(goCoreConfig.HttpApis.Host)
	buffer.WriteString(`"
ApiServicePort = "`)
	buffer.WriteString(goCoreConfig.HttpApis.Port)
	buffer.WriteString(`"

`)
	buffer.WriteString(baseConf)
	buffer.WriteString(`
` + "`" + ``)

}
