package cmd

import (
	"log"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/conf"

	"github.com/urfave/cli/v2"
)

func creatYaml(c *cli.Context) error {
	_, err := CreatYoml(".", conf.GetGocoreConfig())
	if err != nil {
		return err
	}
	log.Println("gocore.yaml 已生成...")
	return nil
}
