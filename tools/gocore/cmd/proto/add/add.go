package add

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Add represents the add command.
// gocore proto add api/hello.proto
var Add = &cli.Command{
	Name:        "add",
	Usage:       "Add a proto API template",
	Description: `Add a proto API template. Example: gocore proto add helloworld/v1/hello.proto`,
	Action:      run,
}

func run(c *cli.Context) error {
	if c.NArg() == 0 {
		fmt.Println("Please enter the proto file or directory")
		return nil
	}
	input := c.Args().Get(0)
	n := strings.LastIndex(input, "/")
	if n == -1 {
		fmt.Println("The proto path needs to be hierarchical.")
		return nil
	}
	path := input[:n]
	fileName := input[n+1:]
	pkgName := strings.ReplaceAll(path, "/", ".")

	p := &Proto{
		Name:        fileName,
		Path:        path,
		Package:     pkgName,
		GoPackage:   goPackage(path),
		JavaPackage: javaPackage(pkgName),
		Service:     serviceName(fileName),
	}
	if err := p.Generate(); err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Printf("Created %s\n", input)
	return nil
}

func modName() string {
	modBytes, err := os.ReadFile("go.mod")
	if err != nil {
		if modBytes, err = os.ReadFile("../go.mod"); err != nil {
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}

func goPackage(path string) string {
	s := strings.Split(path, "/")
	return modName() + "/" + path + ";" + s[len(s)-1]
}

func javaPackage(name string) string {
	return name
}

func serviceName(name string) string {
	return toUpperCamelCase(strings.Split(name, ".")[0])
}

func toUpperCamelCase(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = cases.Title(language.Und, cases.NoLower).String(s)
	return strings.ReplaceAll(s, " ", "")
}
