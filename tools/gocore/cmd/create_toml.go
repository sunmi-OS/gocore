package cmd

import (
	"log"

	"github.com/sunmi-OS/gocore/tools/gocore/template"

	"github.com/urfave/cli"
)

func creatToml(c *cli.Context) error {
	dir := c.String("dir")
	if dir != "" {
		dir = dir + "/gocore.toml"
	} else {
		dir = "gocore.toml"
	}
	writer.Add([]byte(template.CreateToml()))
	writer.WriteToFile(dir)
	log.Println("gocore.toml 已生成...")
	return nil
}
