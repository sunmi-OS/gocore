package template

import (
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cast"
)

var bar *progressbar.ProgressBar
var schedule int

func newProgress(max int, description string) {
	schedule = 0
	bar = progressbar.NewOptions(max,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetDescription("[cyan]["+cast.ToString(schedule)+"/"+cast.ToString(max)+"][reset] "+description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

}

func progressNext(description string) {
	time.Sleep(200 * time.Millisecond)
	schedule++
	bar.Describe("[cyan][" + cast.ToString(schedule) + "/" + cast.ToString(bar.GetMax()) + "][reset] " + description)
	err := bar.Add(1)
	if err != nil {
		return
	}
}
