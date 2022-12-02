package nav

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/websocket"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
	"github.com/mlctrez/goapp-audioplayer/model"
)

type Navigation struct {
	app.Compo
}

func (n *Navigation) Render() app.UI {
	return app.Div().Class("navigation").Body(
		app.Div().Class("navigation-left").Body(
			app.Div().Body(app.Img().Src("/web/logo-192.png")),
			app.Div().Body(app.Text("Music")),
		),
		app.Div().Class("navigation-center").Body(
			app.Div().Body(app.Text("Home")),
			app.Div().Body(app.Text("Explore")).Style("cursor", "pointer").OnClick(func(ctx app.Context, e app.Event) {
				websocket.Action(ctx).Write(&model.AlbumsRequest{})
			}),
			app.Div().Body(app.Text("Search")),
			app.Raw(icon.Search48()),
		),
		app.Div().Class("navigation-right").Body(
			app.Div().Body(app.Text("PROFILE")),
		),
	)
}
