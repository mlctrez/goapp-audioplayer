package control

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
)

type PlayPause struct {
	app.Compo
	goapp.Logging
	playing bool
}

func (p *PlayPause) OnMount(ctx app.Context) {
	p.Log("")
	ctx.Handle(audio.EventPause, func(_ app.Context, _ app.Action) {
		p.playing = false
	})
	ctx.Handle(audio.EventPlay, func(_ app.Context, _ app.Action) {
		p.playing = true
	})
}

func (p *PlayPause) Render() app.UI {
	if p.playing {
		return Div(true).Body(app.Raw(icon.Pause48())).OnClick(p.pause)
	} else {
		return Div(true).Body(app.Raw(icon.PlayArrow48())).OnClick(p.play)
	}
}

func (p *PlayPause) play(ctx app.Context, _ app.Event) {
	p.Log("")
	audio.Action(ctx).Play()
}

func (p *PlayPause) pause(ctx app.Context, _ app.Event) {
	p.Log("")
	audio.Action(ctx).Pause()
}
