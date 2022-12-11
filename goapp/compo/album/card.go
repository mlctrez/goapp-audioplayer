package album

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/actions"
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-natsws/natsws"
)

type Card struct {
	app.Compo
	ReleaseGroupID string
	Album          string
	Artist         string
	natswsConn     *natsws.Connection
	albumResponse  *model.AlbumResponse
}

func (c *Card) OnMount(ctx app.Context) {
	c.natswsConn = &natsws.Connection{}
	natsws.Observe(ctx, c.natswsConn)
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
	actions.RequestAlbum(ctx, c.natswsConn, &model.AlbumRequest{ReleaseGroupID: c.ReleaseGroupID})
}
