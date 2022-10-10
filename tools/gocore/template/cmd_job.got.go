package template

import "bytes"

func FromCmdJob(name, jobCmd, jobFunctions string, buffer *bytes.Buffer) {
	buffer.WriteString(`
package cmd

import (
	"`)
	buffer.WriteString(name)
	buffer.WriteString(`/job"
	"github.com/urfave/cli/v2"
	"github.com/sunmi-OS/gocore/v2/utils/closes"
)

// Job cmd 任务相关
var Job = &cli.Command{
	Name:    "job",
	Aliases: []string{"j"},
	Usage:   "job",
	Subcommands: []*cli.Command{
		`)
	buffer.WriteString(jobCmd)
	buffer.WriteString(`
	},
}
`)
	buffer.WriteString(jobFunctions)

}
