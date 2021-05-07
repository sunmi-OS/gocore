package constant

const (
	MainTemplate = `
package main

import (

	gocoreLog "github.com/sunmi-OS/gocore/log"

)

func main() {

	//初始化log
	gocoreLog.InitLogger("order")
	
}
`
)

var (
	// PkgMap imports head options. import包含选项
	PkgMap = map[string]string{
		"stirng":     `"string"`,
		"time.Time":  `"time"`,
		"gorm.Model": `"github.com/jinzhu/gorm"`,
		"fmt":        `"fmt"`,
	}
)
