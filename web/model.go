package web

import (
	"io"
	"os"
)

// DefaultWriter is the default io.Writer used by Gin for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
// 		import "github.com/mattn/go-colorable"
// 		gin.DefaultWriter = colorable.NewColorableStdout()
var DefaultWriter io.Writer = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Gin to debug errors
var DefaultErrorWriter io.Writer = os.Stderr

type Config struct {
	// http export port. :8080
	Port string
}

type RecoverInfo struct {
	Time  string      `json:"time"`
	Url   string      `json:"url"`
	Err   string      `json:"error"`
	Query interface{} `json:"query"`
	Stack string      `json:"stack"`
}

type CommonRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
