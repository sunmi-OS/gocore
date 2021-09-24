package cmd

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"

	//_ "github.com/sunmi-OS/gocore/v2/tools/gocore/ui/statik"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/rakyll/statik/fs"
	"github.com/sunmi-OS/gocore/v2/tools/gocore/ui/view"
	"github.com/urfave/cli/v2"
	"github.com/zserge/lorca"
)

// Ui 运行UI程序
var Ui = &cli.Command{
	Name: "ui",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "dir",
			Usage:       "dir path",
			DefaultText: ".",
		}},
	Usage:  "run ui [dir]",
	Action: ui,
}

// GOARCH=wasm GOOS=js go build -o web/app.wasm
// ui 运行UI程序
func ui(c *cli.Context) error {
	app.Route("/", &view.Home{})
	app.RunWhenOnBrowser()

	http.Handle("/", &app.Handler{
		Title: "gocore",
		Icon: app.Icon{
			Default: "https://file.cdn.sunmi.com/gocore-logo.png",
		},
		Image:       "https://file.cdn.sunmi.com/gocore-logo.png",
		Name:        "Hello",
		Description: "An Hello World! example",
		Styles: []string{
			"https://unpkg.com/layui@2.6.8/dist/css/layui.css",
		},
		Scripts: []string{
			"https://cdn.staticfile.org/jquery/2.1.1/jquery.min.js",
			"https://unpkg.com/layui@2.6.8/dist/layui.js",
		},
		LoadingLabel: "GoCore Loading... ...",
		Version:      "V1.0.0",
	})

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(statikFS)))

	go func() {
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatal(err)
		}
	}()

	err = gui("http://localhost:8000")
	if err != nil {
		log.Println(err)
		err = open("http://localhost:8000")
		if err != nil {
			return err
		}
	}
	return nil
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func gui(url string) error {
	ui, err := lorca.New(url, "", 1920, 1080, "--class=Lorca")
	if err != nil {
		return err
	}

	defer func(ui lorca.UI) {
		err := ui.Close()
		if err != nil {
			log.Println(err)
		}
	}(ui)
	<-ui.Done()
	return nil
}
