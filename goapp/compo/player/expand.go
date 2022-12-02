package player

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/queue"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
)

type Expand struct {
	app.Compo
	queue queue.Queue
}

func (ex *Expand) OnMount(ctx app.Context) {
	ctx.ObserveState("queue").Value(&ex.queue)
}

func (ex *Expand) Render() app.UI {
	var svg app.UI
	// inverted since we're using it from the bottom up
	if ex.queue.Shown {
		svg = app.Raw(icon.ExpandMore48())
	} else {
		svg = app.Raw(icon.ExpandLess48())
	}
	return app.Div().Class("queue-expand").Body(svg).OnClick(ex.queue.Toggle)
}
