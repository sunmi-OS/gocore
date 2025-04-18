package server

import (
	"github.com/urfave/cli/v2"
)

var targetDir string

func init() {
	targetDir = "internal/service"
}

// Server represents the server command.
var Server = &cli.Command{
	Name:        "server",
	Usage:       "Generate the proto server implementations",
	Description: `Generate the proto server implementations. Example: gocore proto server api/xxx.proto --target-dir=internal/service`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "target-dir",
			Aliases:     []string{"t"},
			Value:       targetDir,
			Usage:       "generate target directory",
			Destination: &targetDir,
		},
	},
	Action: run,
}

func run(c *cli.Context) error {
	// 生成service的golang代码
	if c.NArg() == 0 {
		return cli.Exit("Please enter the proto file or directory", 1)
	}
	input := c.Args().Get(0)
	_ = input
	return nil
}
