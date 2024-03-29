package template

import "bytes"

func FromModelTable(dbName, tableStruct, tableName, fields string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package `)
	buffer.WriteString(dbName)
	buffer.WriteString(`

var `)
	buffer.WriteString(tableStruct)
	buffer.WriteString(`Handler = &`)
	buffer.WriteString(tableStruct)
	buffer.WriteString(`{}

type `)
	buffer.WriteString(tableStruct)
	buffer.WriteString(` struct {
	`)
	buffer.WriteString(fields)
	buffer.WriteString(`
}

func (* `)
	buffer.WriteString(tableStruct)
	buffer.WriteString(`) TableName() string {
	return "`)
	buffer.WriteString(tableName)
	buffer.WriteString(`"
}

`)

}
