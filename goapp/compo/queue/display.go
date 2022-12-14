package queue

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
	"github.com/mlctrez/goapp-audioplayer/model"
	"time"
)

type Display struct {
	app.Compo
	goapp.Logging
	queue       Queue
	displayMode string
}

func (d *Display) Render() app.UI {

	if d.displayMode != "queue.Display" {
		return nodisplay.NoDisplay("compo/queue/Display")
	}

	var rows []app.UI

	formatDuration := func(d time.Duration) string {
		if d.Hours() > 1 {
			return fmt.Sprintf("%2d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
		}
		return fmt.Sprintf("%2d:%02d", int(d.Minutes()), int(d.Seconds())%60)
	}

	var totalDuration time.Duration

	for indexLoop, trackLoop := range d.queue.Tracks {
		index := indexLoop
		md := trackLoop
		size := 40
		image := app.Img().Width(size).Height(size).Src(md.CoverArtUrl(size))

		duration := time.Second * time.Duration(md.Seconds)
		totalDuration += duration

		tr := app.Tr().Body(
			app.Td().Body(image),
			app.Td().Class("queue-table-album").Text(md.Album),
			app.Td().Text(md.Title),
			app.Td().Class("queue-table-artist").Text(md.Artist),
			app.Td().Style("text-align", "right").Text(formatDuration(duration)),
		)

		if index == d.queue.Index {
			tr.Style("background-color", "#222")
		}
		tr.Style("cursor", "pointer").OnClick(func(ctx app.Context, e app.Event) {
			d.queue.Seek(ctx, index)
		})

		rows = append(rows, tr)
	}

	queueClear := func(ctx app.Context, e app.Event) { ctx.NewAction("queue.clear") }
	tr := app.Tr().Body(
		app.Td().Body(app.Div().Class("queue-clear").Body(app.Raw(icon.Close48())).OnClick(queueClear)),
		app.Td().Class("queue-table-album").Text(""),
		app.Td().Body(app.Div().Class("queue-clear").Body(app.Text("Clear Queue")).OnClick(queueClear)),
		app.Td().Class("queue-table-artist").Text(""),
		app.Td().Style("text-align", "right").Text(formatDuration(totalDuration)),
	)
	rows = append(rows, tr)

	table := app.Table().Class("queue-table").Class("queue-table").Body(rows...)

	return app.Div().Class("main-content").Body(table)

}

func (d *Display) OnMount(ctx app.Context) {
	ctx.ObserveState("queue").Value(&d.queue)
	ctx.ObserveState("displayMode").Value(&d.displayMode)

	if d.queue.HasCurrent() {
		d.queue.SetCurrent(ctx)
	}

	ctx.Handle("queue.add", d.add)
	ctx.Handle("queue.clear", d.clear)

	ctx.Handle("mediaSession.previoustrack", d.previousTrack)
	ctx.Handle("mediaSession.nexttrack", d.nextTrack)

}

func (d *Display) add(ctx app.Context, action app.Action) {

	wasEmpty := len(d.queue.Tracks) == 0

	switch v := action.Value.(type) {
	case *model.Metadata:
		d.queue.Tracks = append(d.queue.Tracks, v)
	case []*model.Metadata:
		d.queue.Tracks = append(d.queue.Tracks, v...)
	default:
		return
	}
	if wasEmpty {
		d.queue.StartCurrent(ctx)
		ctx.NewAction(audio.EventPlay)
	}
	d.queue.persist(ctx)
}

func (d *Display) clear(ctx app.Context, _ app.Action) {
	// stop the playing audio if any
	audio.Action(ctx).Src(nil)
	ctx.SetState("displayMode", "album.List")
	// clear the queue
	d.queue.Clear(ctx)
}

func (d *Display) previousTrack(ctx app.Context, _ app.Action) {
	d.queue.Previous(ctx)
}

func (d *Display) nextTrack(ctx app.Context, _ app.Action) {
	d.queue.Next(ctx)
}
