package compo

import (
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
	goapp.Logging
}

func (r *Root) OnMount(_ app.Context) {
	r.Logf("version %s", goapp.RuntimeVersion())
}

func (r *Root) Render() app.UI {
	return app.Div().ID("compo-Root").Body(
		&updater.Updater{},
		&websocket.WebSocket{},
		&audio.Audio{},
		&nav.Navigation{},
		&player.Player{},
		&queue.Display{},
		&album.List{},
		&album.Album{},
	)
}
