package queue

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
	"github.com/mlctrez/goapp-audioplayer/model"
	"time"
)

type Component struct {
	app.Compo
	queue Queue
}

func (c *Component) Render() app.UI {

	if !c.queue.Shown {
		return nodisplay.NoDisplay("compo/queue/Component")
	}

	var rows []app.UI

	formatDuration := func(d time.Duration) string {
		if d.Hours() > 1 {
			return fmt.Sprintf("%2d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
		}
		return fmt.Sprintf("%2d:%02d", int(d.Minutes()), int(d.Seconds())%60)
	}

	var totalDuration time.Duration

	for indexLoop, trackLoop := range c.queue.Tracks {
		index := indexLoop
		md := trackLoop
		image := app.Img().Width(30).Height(30).Src(fmt.Sprintf("/cover/%s", md.MusicbrainzReleaseGroupId))

		duration := time.Second * time.Duration(md.Seconds)
		totalDuration += duration

		tr := app.Tr().Body(
			app.Td().Body(image),
			app.Td().Text(md.Album),
			app.Td().Text(md.Title),
			app.Td().Text(md.Artist),
			app.Td().Style("text-align", "right").Text(formatDuration(duration)),
		)

		if index == c.queue.Index {
			tr.Style("background-color", "#222")
		}
		tr.Style("cursor", "pointer").OnClick(func(ctx app.Context, e app.Event) {
			c.queue.Seek(ctx, index)
		})

		rows = append(rows, tr)
	}

	queueClear := func(ctx app.Context, e app.Event) { ctx.NewAction("queue.clear") }
	tr := app.Tr().Body(
		app.Td().Body(app.Div().Class("queue-clear").Body(app.Raw(icon.Close48())).OnClick(queueClear)),
		app.Td().Body(app.Div().Class("queue-clear").Body(app.Text("Clear Queue")).OnClick(queueClear)),
		app.Td().ColSpan(2).Text(""),
		app.Td().Style("text-align", "right").Text(formatDuration(totalDuration)),
	)
	rows = append(rows, tr)

	tableStyle := map[string]string{"width": "60vw", "border-spacing": "0px"}
	table := app.Table().Styles(tableStyle).Class("queue-table").Body(rows...)

	return app.Div().Class("main-content").Body(table)

}

func (c *Component) OnMount(ctx app.Context) {
	ctx.ObserveState("queue").Value(&c.queue)

	//fmt.Printf("queue.Component.OnMount c.queue = %+v\n", c.queue)

	if c.queue.HasCurrent() {
		c.queue.SetCurrent(ctx)
		// this enables the play button
		ctx.NewAction(audio.EventPause)
	}

	ctx.Handle("queue.add", c.add)
	ctx.Handle("queue.clear", c.clear)
	ctx.Handle("queue.toggle", c.toggle)

}

func (c *Component) add(ctx app.Context, action app.Action) {

	wasEmpty := len(c.queue.Tracks) == 0

	switch v := action.Value.(type) {
	case *model.Metadata:
		c.queue.Tracks = append(c.queue.Tracks, v)
	case []*model.Metadata:
		c.queue.Tracks = append(c.queue.Tracks, v...)
	default:
		return
	}
	if wasEmpty {
		c.queue.StartCurrent(ctx)
	}
	c.queue.persist(ctx)
}

func (c *Component) clear(ctx app.Context, _ app.Action) {
	// stop the playing audio if any
	audio.Action(ctx).Src("")
	// clear the queue
	c.queue.Clear(ctx)
}

func (c *Component) toggle(ctx app.Context, _ app.Action) {
	ctx.GetState("queue", &c.queue)
	c.Update()
}
