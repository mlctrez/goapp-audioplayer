package control

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/queue"
)

type BaseControl struct {
	app.Compo
	queue queue.Queue
}

func Div(enabled bool) app.HTMLDiv {
	if enabled {
		return app.Div().Class("audio-control-enabled")
	}
	return app.Div().Class("audio-control")
}
