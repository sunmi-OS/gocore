package proto

import (
	"github.com/sunmi-OS/gocore/v2/tools/gocore/cmd/proto/errcode"
	"github.com/urfave/cli/v2"

	"github.com/sunmi-OS/gocore/v2/tools/gocore/cmd/proto/add"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/cmd/proto/client"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/cmd/proto/server"
)

// Proto represents the proto command.
var Proto = &cli.Command{
	Name:        "proto",
	Usage:       "Generate the proto files",
	Description: `Generate the proto files.`,
	Subcommands: []*cli.Command{
		add.Add,
		client.Client,
		server.Server,
		errcode.ErrCode,
	},
}
