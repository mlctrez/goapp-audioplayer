package player

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
)

type Mute struct {
	app.Compo
	volume         float64
	previousVolume float64
}

func (m *Mute) Render() app.UI {
	control := app.Div().Class("audio-control-mute")
	if m.volume == 0 {
		control.Body(app.Raw(icon.VolumeMute48())).Style("position", "relative").Style("left", "-7px")
	} else if m.volume < .5 {
		control.Body(app.Raw(icon.VolumeDown48())).Style("position", "relative").Style("left", "-3px")
	} else {
		control.Body(app.Raw(icon.VolumeUp48()))
	}
	control.OnClick(m.clicked)
	return control
}

func (m *Mute) OnMount(ctx app.Context) {
	ctx.ObserveState("volume").Value(&m.volume)
}

func (m *Mute) clicked(ctx app.Context, _ app.Event) {
	if m.volume == 0 {
		// restore previous volume value
		ctx.SetState("volume", m.previousVolume)
		audio.Action(ctx).Volume(m.previousVolume)
	} else {
		// save current volume and mute
		m.previousVolume = m.volume
		ctx.SetState("volume", float64(0))
		audio.Action(ctx).Volume(float64(0))
	}

}
