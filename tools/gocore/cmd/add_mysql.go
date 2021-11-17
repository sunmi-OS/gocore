package cmd

import (
	"os"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/generate"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// AddMysql 添加已有mysql
var AddMysql = &cli.Command{
	Name:   "mysql",
	Usage:  "mysql",
	Action: addMysql,
	Subcommands: []*cli.Command{
		{
			Name: "create_yaml",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dir",
					Aliases:     []string{"d"},
					Usage:       "dir",
					DefaultText: ".",
				}},
			Usage:  "mysql create_yaml -d xxx",
			Action: createMysqlYaml,
		},
		{
			Name: "add",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "dir",
					Aliases:     []string{"d"},
					Usage:       "dir",
					DefaultText: ".",
				}},
			Usage:  "mysql add -d .",
			Action: addMysql,
		},
	},
}

// addMysql 添加已有mysql
func addMysql(c *cli.Context) error {
	root := c.String("dir")
	mysqlPath := "mysql.yaml"
	gocorePath := "gocore.yaml"
	if root != "" {
		mysqlPath = root + "/mysql.yaml"
		gocorePath = root + "gocore.yaml"
	}
	mysqlByte, err := os.ReadFile(mysqlPath)
	if err != nil {
		return err
	}
	mysqlDb := conf.MysqlDb{}
	err = yaml.Unmarshal(mysqlByte, &mysqlDb)
	if err != nil {
		return err
	}
	db, err := openORM(&mysqlDb)
	if err != nil {
		return err
	}
	config, err := InitYaml(root, conf.GetGocoreConfig())
	if err != nil {
		return err
	}
	config = generate.Genertate(db, &mysqlDb, config)
	_, err = CreateYaml(gocorePath, config)
	if err != nil {
		return err
	}
	printHint("Welcome to GoCore, Configuration file has been generated.")
	return nil
}

func openORM(mysqlDb *conf.MysqlDb) (*gorm.DB, error) {
	dsn := mysqlDb.User + ":" + mysqlDb.Password + "@tcp(" + mysqlDb.Host + ":" + mysqlDb.Port + ")/" + mysqlDb.Database + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return nil, err
	}
	return db, err
}

// createMysqlYaml 创建 MysqlYaml
func createMysqlYaml(c *cli.Context) error {
	root := c.String("dir")
	mysqlPath := "mysql.yaml"
	if root != "" {
		mysqlPath = root + "/mysql.yaml"
	}
	var writer = file.NewWriter()
	yamlByte, err := yaml.Marshal(new(conf.MysqlDb))
	if err != nil {
		return err
	}
	writer.Add(yamlByte)
	writer.WriteToFile(mysqlPath)
	printHint("Welcome to GoCore, Mysql Configuration file has been generated.")
	return nil
}
