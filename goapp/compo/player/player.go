package player

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/player/control"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/queue"
)

type Player struct {
	app.Compo
	expanded bool
	playing  bool
	queue    queue.Queue
}

func (p *Player) OnMount(ctx app.Context) {
	ctx.ObserveState("queue").Value(&p.queue)
	ctx.Handle(audio.EventEnded, func(context app.Context, action app.Action) { p.queue.Next(context) })
}

func (p *Player) Render() app.UI {

	q := p.queue
	if q.HasCurrent() {

		var controls []app.UI
		controls = append(controls, &control.Previous{})
		controls = append(controls, &control.PlayPause{})
		controls = append(controls, &control.Next{})
		controls = append(controls, &CurrentTime{})

		md := q.CurrentTrack()
		imageSrc := fmt.Sprintf("/cover/%s", md.MusicbrainzReleaseGroupId)
		date := md.Date
		if len(date) >= 4 {
			date = " • " + date[0:4]
		} else {
			date = ""
		}

		playerUI := app.Div().Class("queue").Body(
			app.Div().Class("queue-vertical").Body(
				&AudioProgress{Seconds: md.Seconds},
				app.Div().Class("queue-controls").Body(
					app.Div().Class("audio-controls").Body(controls...),
					app.Div().Class("queue-playing").Body(
						app.Img().Src(imageSrc).Class("queue-playing-img"),
						app.Div().Class("queue-playing-text").Body(
							app.B().Text(md.Title), app.Br(),
							app.Text(md.Artist+" • "+md.Album+date),
						),
					),
					&Expand{},
				),
			),
		)

		return app.Div().Body(playerUI)

	}

	return app.Div()

}
