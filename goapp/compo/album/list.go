package album

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/websocket"
	"github.com/mlctrez/goapp-audioplayer/model"
)

type List struct {
	app.Compo
	goapp.Logging
	displayMode     string
	albumCards      []*model.AlbumCard
	albumsScrollTop int
}

func (l *List) Render() app.UI {

	if l.displayMode != "album.List" {
		return nodisplay.NoDisplay("compo/album/List")
	}

	var cards []app.UI
	for _, card := range l.albumCards {
		id := card.ReleaseGroupID
		cardUI := &Card{ReleaseGroupID: id, Album: card.Album, Artist: card.Artist}

		cards = append(cards, cardUI)
	}
	return app.Div().Class("main-content").Body(cards...)

}

func (l *List) OnMount(ctx app.Context) {
	l.Log("")
	ctx.ObserveState("displayMode").Value(&l.displayMode)
	websocket.Action(ctx).HandleAction(&model.AlbumsResponse{}, l)
	ctx.Defer(l.requestAlbums)
}

func (l *List) requestAlbums(ctx app.Context) {
	l.Log("")
	websocket.Action(ctx).Write(&model.AlbumsRequest{})
}

func (l *List) OnWebsocketMessage(ctx app.Context, message model.WebSocketMessage) {
	l.Log("")
	l.albumCards = message.(*model.AlbumsResponse).Results
	ctx.SetState("displayMode", "album.List")
}
