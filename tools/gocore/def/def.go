package def

import (
	"fmt"
	"regexp"
)

// EImportsHead imports head options. import包含选项
var EImportsHead = map[string]string{
	"stirng":     `"string"`,
	"time.Time":  `"time"`,
	"gorm.Model": `"gorm.io/gorm"`,
	"fmt":        `"fmt"`,
}

// TypeMysqlDicMp Accurate matching type.精确匹配类型
var TypeMysqlDicMp = map[string]string{
	"smallint":            "int16",
	"smallint unsigned":   "uint16",
	"int":                 "int",
	"int unsigned":        "uint",
	"bigint":              "int64",
	"bigint unsigned":     "uint64",
	"varchar":             "string",
	"char":                "string",
	"date":                "time.Time",
	"datetime":            "time.Time",
	"bit(1)":              "int8",
	"tinyint":             "int8",
	"tinyint unsigned":    "uint8",
	"tinyint(1)":          "int8",
	"tinyint(1) unsigned": "int8",
	"json":                "string",
	"text":                "string",
	"timestamp":           "time.Time",
	"double":              "float64",
	"mediumtext":          "string",
	"longtext":            "string",
	"float":               "float32",
	"tinytext":            "string",
	"enum":                "string",
}

// TypeMysqlMatchMp Fuzzy Matching Types.模糊匹配类型
var TypeMysqlMatchMp = map[string]string{
	`^(int)`:                  "int64",
	`^(tinyint)[(]\d+[)]`:     "int8",
	`^(smallint)[(]\d+[)]`:    "int16",
	`^(int)[(]\d+[)]`:         "int",
	`^(bigint)[(]\d+[)]`:      "int64",
	`^(char)[(]\d+[)]`:        "string",
	`^(enum)[(](.)+[)]`:       "string",
	`^(varchar)[(]\d+[)]`:     "string",
	`^(varbinary)[(]\d+[)]`:   "[]byte",
	`^(decimal)[(]\d+,\d+[)]`: "float64",
	`^(mediumint)[(]\d+[)]`:   "string",
	`^(double)[(]\d+,\d+[)]`:  "float64",
	`^(float)[(]\d+,\d+[)]`:   "float64",
	`^(datetime)[(]\d+[)]`:    "time.Time",
	`^(timestamp)`:            "time.Time",
	`^(text)`:                 "string",
	`^(date )`:                "time.Time",
}

type Table struct {
	Name    string `gorm:"column:Name"`
	Comment string `gorm:"column:Comment"`
}

type Field struct {
	Field      string `gorm:"column:Field"`
	Type       string `gorm:"column:Type"`
	Null       string `gorm:"column:Null"`
	Key        string `gorm:"column:Key"`
	Default    string `gorm:"column:Default"`
	Extra      string `gorm:"column:Extra"`
	Privileges string `gorm:"column:Privileges"`
	Comment    string `gorm:"column:Comment"`
}

// getTypeName Type acquisition filtering.类型获取过滤
func GetTypeName(name string) string {
	// Precise matching first.先精确匹配
	if v, ok := TypeMysqlDicMp[name]; ok {
		return v
	}

	// Fuzzy Regular Matching.模糊正则匹配
	for k, v := range TypeMysqlMatchMp {
		if ok, _ := regexp.MatchString(k, name); ok {
			return v
		}
	}

	panic(fmt.Sprintf("type (%v) not match in any way.maybe need to add on (https://github.com/xxjwxc/gormt/blob/master/data/view/cnf/def.go)", name))
}
