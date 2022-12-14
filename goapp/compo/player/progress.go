package player

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
)

type AudioProgress struct {
	app.Compo
	timeUpdate *audio.TimeUpdate
	progress   float64
}

func (ap *AudioProgress) Render() app.UI {
	return app.Div().Class("audio-progress").OnClick(ap.seek).
		Body(app.Div().Class("audio-slider").
			Style("width", fmt.Sprintf("%dpx", int(ap.progress))),
		)
}

func (ap *AudioProgress) seek(ctx app.Context, e app.Event) {

	if ap.timeUpdate != nil {

		atPoint := e.Get("offsetX").Float()
		width := ap.JSValue().Get("clientWidth").Float()
		currentTime := ap.timeUpdate.Duration * atPoint / width

		audio.Action(ctx).CurrentTime(currentTime)
	}
}

func (ap *AudioProgress) OnMount(ctx app.Context) {
	ctx.Handle(audio.EventTimeUpdate, ap.update)
}

func (ap *AudioProgress) update(_ app.Context, action app.Action) {
	if tu, ok := action.Value.(*audio.TimeUpdate); ok {
		width := ap.JSValue().Get("clientWidth").Float()
		ap.progress = width * tu.CurrentTime / tu.Duration
		ap.timeUpdate = tu
	}
}
