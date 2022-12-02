package compo

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/album"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nav"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/player"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/queue"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/updater"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/websocket"
)

type Root struct {
	app.Compo
}

func (r *Root) OnMount(ctx app.Context) {
	fmt.Println("RuntimeVersion", goapp.RuntimeVersion())
}

func (r *Root) Render() app.UI {
	return app.Div().Body(
		&updater.Updater{},
		&websocket.WebSocket{},
		&audio.Audio{},
		&nav.Navigation{},
		&queue.Component{},
		&album.Album{},
		&album.List{},
		&player.Player{},
	)
}
