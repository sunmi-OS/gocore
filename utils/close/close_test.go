package close

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {

	AddShutdown(Close{
		Name:     "http",
		Priority: 0,
		Func: func() {
			fmt.Println("close http")
		},
	})
	go SignalClose()
	time.Sleep(1 * time.Second)

}

func TestSort(t *testing.T) {

	closeHandler = append(closeHandler, Close{
		Name:     "grom",
		Priority: 100,
		Func: func() {
			fmt.Println("close http")
		},
	}, Close{
		Name:     "http",
		Priority: 0,
		Func: func() {
			fmt.Println("close http")
		},
	})

	sort.Sort(closeHandler)

	assert.Equal(t, closeHandler[0].Priority, 0)

}
