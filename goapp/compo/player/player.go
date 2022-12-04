package player

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/player/control"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/queue"
	"github.com/mlctrez/goapp-audioplayer/model"
)

type Player struct {
	app.Compo
	goapp.Logging
	expanded bool
	playing  bool
	queue    queue.Queue
}

func (p *Player) OnMount(ctx app.Context) {
	p.Log("")
	ctx.ObserveState("queue").Value(&p.queue)
	ctx.Handle(audio.EventEnded, func(context app.Context, action app.Action) { p.queue.Next(context) })
}

func (p *Player) Render() app.UI {
	var controls []app.UI
	controls = append(controls, &control.Previous{})
	controls = append(controls, &control.PlayPause{})
	controls = append(controls, &control.Next{})
	controls = append(controls, &CurrentTime{})
	controls = append(controls, &Volume{})
	controls = append(controls, &Mute{})

	var md *model.Metadata

	var date string
	var imageSrc string

	if p.queue.HasCurrent() {
		md = p.queue.CurrentTrack()
		date = md.Date
		if len(date) >= 4 {
			date = " • " + date[0:4]
		} else {
			date = ""
		}
		imageSrc = fmt.Sprintf("/cover/%s", md.MusicbrainzReleaseGroupId)
	} else {
		md = &model.Metadata{
			Title:  "Title",
			Artist: "Artist",
			Album:  "Album",
		}
		imageSrc = "/web/logo-192.png"
	}

	playerUI := app.Div().Class("queue").Body(
		app.Div().Class("queue-vertical").Body(
			&AudioProgress{},
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

	if p.queue.HasCurrent() {
		return app.Div().Body(playerUI)
	} else {
		return nodisplay.NoDisplay("compo/player/Player").Body(playerUI)
	}

}
