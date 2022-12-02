package album

import (
	"fmt"
	app "github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/queue"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/websocket"
	"github.com/mlctrez/goapp-audioplayer/model"
)

type List struct {
	app.Compo
	queue           queue.Queue
	albumCards      []*model.AlbumCard
	albumsScrollTop int
}

func (l *List) Render() app.UI {

	if l.queue.Shown {
		return nodisplay.NoDisplay("compo/album/List")
	}

	var cards []app.UI
	for _, card := range l.albumCards {
		cardUI := &Card{
			ReleaseGroupID: card.ReleaseGroupID,
			Album:          card.Album,
			Artist:         card.Artist,
		}

		cards = append(cards, cardUI)
	}
	return app.Div().Class("main-content").Body(cards...)

}

func (l *List) OnMount(ctx app.Context) {
	ctx.ObserveState("queue").Value(&l.queue)
	websocket.Action(ctx).HandleAction(&model.AlbumsResponse{}, func(message model.WebSocketMessage) {
		l.albumCards = message.(*model.AlbumsResponse).Results
		l.Update()
	})

	ctx.Defer(func(context app.Context) {
		websocket.Action(context).Write(&model.AlbumsRequest{})
	})

	//// TODO: figure out how to do this without delay, panic on line 73 in web socket if not
	//ctx.After(time.Second*2, func(context app.Context) {
	//	websocket.Action(context).Write(&model.AlbumsRequest{})
	//})
}

func (l *List) scrollToAlbumsScrollTop() {
	fmt.Println("scroll to", l.albumsScrollTop)
	app.Window().Call("scrollTo", app.ValueOf(0), app.ValueOf(l.albumsScrollTop))
}

func getScrollTop() int {
	return app.Window().Get("document").Get("documentElement").Get("scrollTop").Int()
}
