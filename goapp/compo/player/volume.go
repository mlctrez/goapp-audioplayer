package player

import (
	app "github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
)

type Volume struct {
	app.Compo
	value float64
}

func (v *Volume) Render() app.UI {
	input := app.Input().ID("volume-slider").Type("range").Min(0).Max(1)
	input.Value(v.value).Step(0.02).Class("volume-slider")
	input.OnInput(v.valueChanged)
	return app.Div().Class("volume-slider-container").Body(input)
}

func (v *Volume) OnMount(ctx app.Context) {
	var val float64
	ctx.GetState("volume", &val)
	if val == 0 {
		val = 1
		ctx.SetState("volume", val)
	}
	ctx.ObserveState("volume").Value(&v.value)
}

func (v *Volume) valueChanged(ctx app.Context, e app.Event) {
	sliderValue := e.Get("target").Get("valueAsNumber").Float()

	ctx.SetState("volume", sliderValue, app.Persist)
	audio.Action(ctx).Volume(sliderValue)

	//app.Window().Get("console").Call("log", sliderValue)
}
