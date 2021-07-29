package cmd

import (
	"log"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/file"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/template"

	"github.com/urfave/cli/v2"
)

func creatToml(c *cli.Context) error {
	var writer = file.NewWriter()
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
