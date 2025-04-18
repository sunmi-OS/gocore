package errcode

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

var protoPath string

func init() {
	if protoPath = os.Getenv("GOCORE_PROTO_PATH"); protoPath == "" {
		protoPath = "./third_party"
	}
}

var ErrCode = &cli.Command{
	Name:        "errcode",
	Usage:       "Generate the errcode",
	Description: `Generate the proto errcode. Example: gocore proto errcode errcode.proto`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "proto_path",
			Aliases:     []string{"p"},
			Value:       protoPath,
			Usage:       "proto path",
			Destination: &protoPath,
		},
	},
	Action: run,
}

func run(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Println("Please enter the proto file or directory")
		return nil
	}
	var (
		err   error
		proto = strings.TrimSpace(c.Args().Get(0))
	)
	if err = look("protoc-gen-sm-go-errors"); err != nil {
		// update the gocore plugins
		cmd := exec.Command("gocore", "upgrade")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			fmt.Println(err)
			return nil
		}
	}
	if strings.HasSuffix(proto, ".proto") {
		err = generate(proto, c.Args().Slice())
	} else {
		err = walk(proto, c.Args().Slice())
	}
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func look(name ...string) error {
	for _, n := range name {
		if _, err := exec.LookPath(n); err != nil {
			return err
		}
	}
	return nil
}

func walk(dir string, args []string) error {
	if dir == "" {
		dir = "."
	}
	return filepath.Walk(dir, func(path string, _ os.FileInfo, _ error) error {
		fmt.Println("walk path = " + path)
		if ext := filepath.Ext(path); ext != ".proto" || strings.HasPrefix(path, "third_party") {
			return nil
		}
		return generate(path, args)
	})
}

// generate is used to execute the generate command for the specified proto file
// gocore proto errcode errcode.proto
// gocore proto errcode errcode.proto --sm-go-errors_out=fe_ecode=./api/docs/fe_code
func generate(proto string, args []string) error {
	// todo把third_party拷贝到/tmp/third_party
	theDir := path.Dir(proto)
	input := []string{
		"--proto_path=.",
	}
	if pathExists(protoPath) {
		input = append(input, "--proto_path="+protoPath)
	}

	inputExt := []string{
		//"--sm-go-errors_out=paths=source_relative:.",
		fmt.Sprintf("--sm-go-errors_out=fe_ecode=%v/docs/fe_ecode:.", theDir),
	}
	input = append(input, inputExt...)
	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if err := fd.Run(); err != nil {
		return err
	}
	fmt.Printf("proto: %s\n", proto)
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
