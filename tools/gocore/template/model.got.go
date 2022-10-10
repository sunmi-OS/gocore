package template

import (
	"bytes"
	"strings"
)

func FromModel(dbName, tabels string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package `)
	buffer.WriteString(dbName)
	buffer.WriteString(`

import (
	"fmt"

	"`)
	buffer.WriteString(goCoreConfig.Service.ProjectName)
	buffer.WriteString(`/conf"
	"gorm.io/gorm"
	"github.com/sunmi-OS/gocore/v2/db/orm"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/utils"
)

func Orm() *gorm.DB {
	db := orm.GetORM(conf.DB`)
	buffer.WriteString(strings.Title(dbName))
	buffer.WriteString(`)
	if 	viper.C.GetBool("base.debug") {
		db = db.Debug()
	}
	return db
}`)

}
