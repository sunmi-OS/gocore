package view

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type Home struct {
	app.Compo
}

func (h *Home) Render() app.UI {

	return app.Div().Body(
		app.Div().Body(
			app.Button().Body(
				app.Text("Go2"),
			).OnClick(func(ctx app.Context, e app.Event) {
				ctx.Navigate("https://baidu.com")
			}),
		),
	)
}
