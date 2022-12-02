package nav

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/websocket"
	"github.com/mlctrez/goapp-audioplayer/model"
	"strings"
)

type Navigation struct {
	app.Compo
}

func (n *Navigation) Render() app.UI {
	return app.Div().Class("navigation").Body(
		app.Div().Class("navigation-left").Body(
			app.Div().Body(app.Img().Src("/web/logo-192.png")),
			app.Div().Body(app.Text("Music")),
		).OnClick(n.navigationLeft),
		app.Div().Class("navigation-center").Body(
			//app.Div().Body(app.Text("Home")),
			app.Div().Body(app.Text("Explore")).Style("cursor", "pointer").OnClick(func(ctx app.Context, e app.Event) {
				websocket.Action(ctx).Write(&model.AlbumsRequest{})
			}),
			//app.Div().Body(app.Text("Search")),
			//app.Raw(icon.Search48()),
		),
		n.version(),
	)
}

func (n *Navigation) version() app.UI {

	div := app.Div().Class("navigation-right")
	if goapp.IsDevelopment() {
		return div.Body(app.Div().Text(goapp.RuntimeVersion()[0:5]))
	} else {
		version := strings.Split(goapp.RuntimeVersion(), "@")[0]
		href := fmt.Sprintf("https://github.com/mlctrez/goapp-audioplayer/tree/%s", version)
		return div.Body(app.A().Href(href)).Text(version)
	}
}

func (n *Navigation) navigationLeft(ctx app.Context, e app.Event) {
	ctx.Reload()
}
