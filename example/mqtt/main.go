package main

import (
	"fmt"
	"github.com/sunmi-OS/gocore/mqtt"
)

func main() {
	fmt.Println(mqtt.QOS_0)
	fmt.Println(mqtt.QOS_1)
	fmt.Println(mqtt.QOS_2)
}
