package player

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
)

type Expand struct {
	app.Compo
	displayMode  string
	previousMode string
}

func (ex *Expand) OnMount(ctx app.Context) {
	ctx.ObserveState("displayMode").Value(&ex.displayMode)
}

func (ex *Expand) Render() app.UI {
	var svg app.UI
	// inverted since we're using it from the bottom up
	if ex.displayMode == "queue.Display" {
		svg = app.Raw(icon.ExpandMore48())
	} else {
		svg = app.Raw(icon.ExpandLess48())
	}
	return app.Div().Class("queue-expand").Body(svg).OnClick(ex.click)
}

func (ex *Expand) click(ctx app.Context, _ app.Event) {
	if ex.displayMode == "queue.Display" {
		ctx.SetState("displayMode", ex.previousMode)
	} else {
		ex.previousMode = ex.displayMode
		ctx.SetState("displayMode", "queue.Display")
	}
}
