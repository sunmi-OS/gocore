package view

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type Home struct {
	app.Compo
}

func (h *Home) Render() app.UI {

	return app.Div().Class("layui-layout layui-layout-admin").Body(
		&Navigate{},
		app.Div().Class("layui-body").Body(
			app.Button().Class("layui-btn layui-btn-normal").Body(
				app.Text("Go2"),
			).OnClick(func(ctx app.Context, e app.Event) {
				//ctx.Navigate("https://baidu.com")
			}),
		))
}
