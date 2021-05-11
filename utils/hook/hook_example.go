package hook

import (
	"fmt"
	"time"
)

func ShowExample() {
	//该方法是线程安全的
	AddShutdownHook(func() int {
		fmt.Println("first task exit...")
		return 0 // the program exit code
	}, func() int {
		fmt.Println("second task exit...")
		return 1
	})

	fmt.Println("some log")

	AddShutdownHook(func() int {
		fmt.Println("third task exit...")
		AddShutdownHook(nil)
		fmt.Println("success graceful shutdown the program")
		return 2
	})

	fmt.Println("service blocking..., please enter `ctr C` or `ctr Z`")
	time.Sleep(time.Hour * 999)
}
