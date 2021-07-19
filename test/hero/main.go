package main

import (
	"bytes"
	"fmt"

	"github.com/sunmi-OS/gocore/v2/test/hero/template"
)

// hero -source=./test/hero/template -extensions=.got,.html,.docker,.md
func main() {

	buffer := new(bytes.Buffer)

	template.FromMain("test", []string{"Api"}, buffer)

	fmt.Println(buffer.String())

}
