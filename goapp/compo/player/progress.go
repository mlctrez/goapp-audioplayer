package player

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
)

type AudioProgress struct {
	app.Compo
	Seconds  uint64
	progress float64
}

func (ap *AudioProgress) Render() app.UI {
	return app.Div().Class("audio-progress").OnClick(func(ctx app.Context, e app.Event) {

		atPoint := e.Get("offsetX").Float()
		width := ap.JSValue().Get("clientWidth").Float()
		currentTime := float64(ap.Seconds) * atPoint / width

		audio.Action(ctx).CurrentTime(currentTime)

	}).Body(app.Div().Class("audio-slider").
		Style("width", fmt.Sprintf("%dpx", int(ap.progress))),
	)
}

func (ap *AudioProgress) OnMount(ctx app.Context) {
	ctx.Handle(audio.EventTimeUpdate, func(context app.Context, action app.Action) {
		if tu, ok := action.Value.(*audio.TimeUpdate); ok {
			width := ap.JSValue().Get("clientWidth").Float()
			ap.progress = width * tu.CurrentTime / tu.Duration
			ap.Update()
		}
	})
}
