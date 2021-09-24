package generate

import (
	"strings"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/def"
	"gorm.io/gorm"
)

func Genertate(gormDb *gorm.DB, mysqlDb *conf.MysqlDb, config *conf.GoCore) *conf.GoCore {
	tableNames := []string{}
	for k1 := range mysqlDb.Tables {
		tableNames = append(tableNames, mysqlDb.Tables[k1])
	}

	//生成所有表信息
	tables := getTables(gormDb, mysqlDb.Database, tableNames)
	models := make([]conf.Model, 0)
	for _, table := range tables {
		fileds := make([]string, 0)
		fieldList := getFields(gormDb, table.Name)
		for _, v1 := range fieldList {
			field := "column:" + v1.Field + ";type:" + v1.Type
			if v1.Null == "NO" {
				field += "  NOT NULL"
			}
			field += ";default:" + v1.Default + " " + v1.Extra + ";" + "comment:'" + v1.Comment + "';"
			fileds = append(fileds, field)
		}
		model := conf.Model{
			Name:    table.Name,
			Comment: table.Comment,
			Fields:  fileds,
		}
		models = append(models, model)
	}
	databaseName := strings.ReplaceAll(mysqlDb.Database, "-", "_")
	mysql := conf.Mysql{
		Name:   databaseName,
		Models: models,
	}
	config.Config.CMysql = append(config.Config.CMysql, mysql)
	return config
}

//获取表信息
func getTables(gormDb *gorm.DB, databaseName string, tableNames []string) []def.Table {
	var tables []def.Table
	if len(tableNames) == 0 {
		gormDb.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema=?;", databaseName).Find(&tables)
	} else {
		gormDb.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE TABLE_NAME IN (?) AND table_schema=?;", tableNames, databaseName).Find(&tables)
	}
	return tables
}

//获取所有字段信息
func getFields(gormDb *gorm.DB, tableName string) []def.Field {
	var fields []def.Field
	gormDb.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}
