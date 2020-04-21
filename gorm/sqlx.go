package gorm

import (
	"strings"

	"github.com/jinzhu/gorm"
)

// @desc 批量保存 如果唯一索引重复则更新，唯一索引不重复或者不存在唯一索引则插入
// @auth liuguoqiang 2020-04-21
// @param
// @return
func BulkSave(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	sql := "INSERT INTO `" + table + "` (" + strings.Join(fields, ",") + ") VALUES "
	updateArr := make([]string, 0)
	args := make([]interface{}, 0)
	valueArr := make([]string, 0)
	varArr := make([]string, 0)
	for _, value := range fields {
		updateArr = append(updateArr, value+"=VALUES("+value+")")
	}
	for _, obj := range params {
		varArr = varArr[:0]
		varStr := "("
		for _, value := range fields {
			varArr = append(varArr, "?")
			args = append(args, obj[value])
		}
		varStr += strings.Join(varArr, ",") + ")"
		valueArr = append(valueArr, varStr)
	}
	sql += strings.Join(valueArr, ",")
	sql += " ON DUPLICATE KEY UPDATE " + strings.Join(updateArr, ",")
	err := db.Exec(sql, args...).Error
	return err
}
