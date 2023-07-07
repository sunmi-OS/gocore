package orm

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Deprecated
// BulkInsert 批量插入 不对唯一索引做更新
func BulkInsert(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	if len(params) == 0 {
		return nil
	}
	sql := "INSERT INTO `" + table + "` (`" + strings.Join(fields, "`,`") + "`) VALUES "
	args := make([]interface{}, 0)
	valueArr := make([]string, 0)
	varArr := make([]string, 0)
	for _, obj := range params {
		varArr = varArr[:0]
		varStr := "("
		for _, value := range fields {
			if _, ok := obj[value]; !ok {
				return fmt.Errorf("%s:字段在map中不存在", value)
			}
			varArr = append(varArr, "?")
			args = append(args, obj[value])
		}
		varStr += strings.Join(varArr, ",") + ")"
		valueArr = append(valueArr, varStr)
	}
	sql += strings.Join(valueArr, ",")
	err := db.Exec(sql, args...).Error
	return err
}

// Deprecated
// BulkSave 批量插入 如果唯一索引重复则更新，唯一索引不重复或者不存在唯一索引则插入
func BulkSave(db *gorm.DB, table string, fields []string, params []map[string]interface{}) error {
	if len(params) == 0 {
		return nil
	}
	sql := "INSERT INTO `" + table + "` (`" + strings.Join(fields, "`,`") + "`) VALUES "
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
			if _, ok := obj[value]; !ok {
				return fmt.Errorf("%s字段在map中不存在", value)
			}
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
