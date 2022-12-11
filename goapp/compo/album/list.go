package album

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/actions"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-natsws/natsws"
)

type List struct {
	app.Compo
	goapp.Logging
	displayMode     string
	albumCards      []*model.AlbumCard
	listPosition    int
	albumsScrollTop int
	natswsConn      *natsws.Connection
}

func (l *List) Render() app.UI {

	if l.displayMode != "album.List" {
		return nodisplay.NoDisplay("compo/album/List")
	}

	var cards []app.UI
	for i, card := range l.albumCards {
		if i < l.listPosition {
			continue
		}

		id := card.ReleaseGroupID
		cardUI := &Card{ReleaseGroupID: id, Album: card.Album, Artist: card.Artist}

		cards = append(cards, cardUI)
		if len(cards) >= 12 {
			break
		}
	}
	return app.Div().Class("main-content").Body(cards...)

}

func (l *List) OnMount(ctx app.Context) {
	l.Log("")
	ctx.ObserveState("displayMode").Value(&l.displayMode)

	ctx.Handle("model.AlbumsResponse", l.onAlbumsResponse)
	ctx.Handle("navigation.previous", l.previous)
	ctx.Handle("navigation.next", l.next)

	l.natswsConn = &natsws.Connection{}
	natsws.Observe(ctx, l.natswsConn).OnChange(func() {
		fmt.Println("List.OnMount natsws.Connection")
		actions.RequestAlbums(ctx, l.natswsConn, &model.AlbumsRequest{})
	})
}

func (l *List) onAlbumsResponse(ctx app.Context, action app.Action) {
	l.Log("")
	l.albumCards = action.Value.(*model.AlbumsResponse).Results
	l.listPosition = 0
	ctx.SetState("navigation.previous", "")
	ctx.SetState("navigation.next", "on")
	ctx.SetState("displayMode", "album.List")
}

func (l *List) next(ctx app.Context, _ app.Action) {
	l.Log("")
	l.listPosition += 12
	if l.listPosition+12 > len(l.albumCards)-1 {
		ctx.SetState("navigation.next", "")
	} else {
		ctx.SetState("navigation.next", "on")
	}
	ctx.SetState("navigation.previous", "on")
	ctx.SetState("displayMode", "album.List")
}

func (l *List) previous(ctx app.Context, _ app.Action) {
	l.Log("")
	l.listPosition -= 12
	if l.listPosition <= 0 {
		l.listPosition = 0
		ctx.SetState("navigation.previous", "")
	}
	ctx.SetState("displayMode", "album.List")
}
