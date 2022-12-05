package album

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/websocket"
	"github.com/mlctrez/goapp-audioplayer/model"
)

type Card struct {
	app.Compo
	ReleaseGroupID string
	Album          string
	Artist         string
}

func (c *Card) Render() app.UI {
	return app.Div().ID(c.ReleaseGroupID).Class("main-content-album-card").Body(
		app.Img().Alt(c.Album).Title(c.Album).
			Src(model.CoverArtUrl(c.ReleaseGroupID, 0)).OnClick(c.click),
		app.Div().Title(c.Album).Class("main-content-album-title").Text(c.Album),
		app.Div().Title(c.Artist).Class("main-content-album-artist").Text(c.Artist),
	)
}

func (c *Card) click(ctx app.Context, _ app.Event) {
	websocket.Action(ctx).Write(&model.AlbumRequest{ReleaseGroupID: c.ReleaseGroupID})
}
