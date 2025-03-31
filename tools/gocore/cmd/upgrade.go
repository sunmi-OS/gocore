package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"
)

var Upgrade = &cli.Command{
	Name:        "upgrade",
	Usage:       "Upgrade the gocore",
	Action:      upgradeGocore,
	Description: `Upgrade the gocore.`,
}

func upgradeGocore(c *cli.Context) error {
	err := GoInstall(
		"github.com/sunmi-OS/gocore/v2/tools/gocore@latest",
		"github.com/sunmi-OS/gocore/v2/tools/protoc-gen/protoc-gen-sm-go-errors@latest",
		"github.com/sunmi-OS/gocore/v2/tools/protoc-gen/protoc-gen-sm-go-gin@latest",
		"github.com/sunmi-OS/gocore/v2/tools/protoc-gen/protoc-gen-sm-go-openapi@latest",
	)
	if err != nil {
		printHint("Upgrade gocore failed")
		return err
	}

	return nil
}

// GoInstall go get path.
func GoInstall(path ...string) error {
	for _, p := range path {
		if !strings.Contains(p, "@") {
			p += "@latest"
		}
		fmt.Printf("go install %s\n", p)
		cmd := exec.Command("go", "install", p)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
