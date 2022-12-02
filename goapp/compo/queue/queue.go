package queue

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"github.com/mlctrez/goapp-audioplayer/model"
)

type Queue struct {
	Index  int
	Tracks []*model.Metadata
}

func (q *Queue) persist(ctx app.Context) {
	ctx.SetState("queue", q, app.Persist)
}

func (q *Queue) HasCurrent() bool {
	if q == nil {
		return false
	}
	return len(q.Tracks) > 0
}

func (q *Queue) CurrentTrack() *model.Metadata {
	return q.Tracks[q.Index]
}

func (q *Queue) HasPrevious() bool {
	return q.Index > 0
}

func (q *Queue) HasNext() bool {
	return q.Index < len(q.Tracks)-1
}

func (q *Queue) SetCurrent(ctx app.Context) {
	ctx.Page().SetTitle(q.CurrentTrack().Title)
	audio.Action(ctx).Src(q.currentUrl())
}

func (q *Queue) StartCurrent(ctx app.Context) {
	ctx.Page().SetTitle(q.CurrentTrack().Title)
	audio.Action(ctx).Start(q.currentUrl())
}

func (q *Queue) currentUrl() string {
	return q.CurrentTrack().FlacUrl()
}

func (q *Queue) Clear(ctx app.Context) {
	q.Index = 0
	q.Tracks = []*model.Metadata{}
	q.persist(ctx)
}

func (q *Queue) Previous(ctx app.Context) {
	if q.HasPrevious() {
		q.Index--
		q.persist(ctx)
		q.StartCurrent(ctx)
	}
}

func (q *Queue) Next(ctx app.Context) {
	if q.HasNext() {
		q.Index++
		q.persist(ctx)
		q.StartCurrent(ctx)
	} else {
		ctx.Page().SetTitle("mlctrez Music")
		// enable play button when last song in queue ends
		ctx.NewAction(audio.EventPause)
	}
}

func (q *Queue) Seek(ctx app.Context, index int) {
	if index > -1 && index < len(q.Tracks)-1 {
		if q.Index == index {
			// don't seek to currently playing track
			return
		}
		q.Index = index
		q.persist(ctx)
		q.StartCurrent(ctx)
	}
}
