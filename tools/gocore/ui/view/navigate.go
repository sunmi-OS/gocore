package view

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type HeaderItem struct {
	href string
	name string
}

type Navigate struct {
	app.Compo

	hitem []HeaderItem
}

func (h *Navigate) Render() app.UI {
	h.initHeaderItem()
	return app.Div().Class("layui-header").Body(
		//app.Div().Class("layui-logo layui-hide-xs layui-bg-black").Text("GoCore"),
		app.Ul().Class("layui-nav  layui-layout-left").Body(
			app.Range(h.hitem).Slice(func(i int) app.UI {
				f := h.hitem[i]
				return app.Li().
					Class("layui-nav-item").Body(
					app.A().Href(f.href).Text(f.name),
				)
			}),
		),
	)
}

func (h *Navigate) initHeaderItem() {
	h.hitem = []HeaderItem{
		{
			href: "",
			name: "基础配置",
		}, {
			href: "",
			name: "Api接口",
		}, {
			href: "",
			name: "常驻任务",
		}, {
			href: "",
			name: "定时任务",
		}, {
			href: "",
			name: "数据库配置",
		},
	}
}
