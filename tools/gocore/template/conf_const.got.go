package template

import (
	"bytes"
	"strings"
)

func FromConfConst(projectName string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package conf

const (
	ProjectName    = "`)
	buffer.WriteString(projectName)
	buffer.WriteString(`"
	ProjectVersion = "v1.0.0"
	`)
	for _, v1 := range goCoreConfig.Config.CMysql {
		buffer.WriteString(`
		DB`)
		buffer.WriteString(strings.Title(v1.Name))
		buffer.WriteString(` = "db`)
		buffer.WriteString(strings.Title(v1.Name))
		buffer.WriteString(`"
	`)
	}
	for _, v1 := range goCoreConfig.Config.CRedis {
		for k2 := range v1.Index {
			buffer.WriteString(strings.Title(v1.Name) + strings.Title(k2))
			buffer.WriteString(`Redis = "`)
			buffer.WriteString(v1.Name)
			buffer.WriteString(`.`)
			buffer.WriteString(k2)
			buffer.WriteString(`"
	  `)
		}
	}
	buffer.WriteString(`
)`)

}
