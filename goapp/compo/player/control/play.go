package control

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
)

type PlayPause struct {
	app.Compo
	enabled     bool
	audioPaused bool
}

func (p *PlayPause) OnMount(ctx app.Context) {
	ctx.Handle(audio.EventCanPlay, func(_ app.Context, _ app.Action) {
		p.enabled = true
	})
	ctx.Handle(audio.EventEnded, func(_ app.Context, _ app.Action) {
		p.enabled = false
	})
	ctx.Handle(audio.EventPause, func(_ app.Context, _ app.Action) {
		p.audioPaused = true
	})
	ctx.Handle(audio.EventPlay, func(_ app.Context, _ app.Action) {
		p.audioPaused = false
	})
}

func (p *PlayPause) Render() app.UI {
	if p.audioPaused {
		return Div(p.enabled).Body(app.Raw(icon.PlayArrow48())).OnClick(p.play)
	} else {
		return Div(p.enabled).Body(app.Raw(icon.Pause48())).OnClick(p.pause)
	}
}

func (p *PlayPause) play(ctx app.Context, _ app.Event) {
	audio.Action(ctx).Play()
}

func (p *PlayPause) pause(ctx app.Context, _ app.Event) {
	audio.Action(ctx).Pause()
}
