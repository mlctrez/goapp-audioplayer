package audio

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp"
	"time"
)

var _ app.Mounter = (*Audio)(nil)

type Actions struct {
	app.Context
}

func Action(ctx app.Context) *Actions {
	return &Actions{ctx}
}

const audioSrc = "audio.action.src"
const audioPlay = "audio.action.play"
const audioStart = "audio.action.start"
const audioPause = "audio.action.pause"
const audioCurrentTime = "audio.action.currentTime"
const audioVolume = "audio.action.volume"

func (ac *Actions) Src(url string)            { ac.NewActionWithValue(audioSrc, url) }
func (ac *Actions) Play()                     { ac.NewAction(audioPlay) }
func (ac *Actions) Start(url string)          { ac.NewActionWithValue(audioStart, url) }
func (ac *Actions) Pause()                    { ac.NewAction(audioPause) }
func (ac *Actions) CurrentTime(value float64) { ac.NewActionWithValue(audioCurrentTime, value) }
func (ac *Actions) Volume(value float64)      { ac.NewActionWithValue(audioVolume, value) }

func (ac *Actions) handle(audio *Audio) {
	ac.Handle(audioSrc, audio.src)
	ac.Handle(audioPlay, audio.play)
	ac.Handle(audioStart, audio.start)
	ac.Handle(audioPause, audio.pause)
	ac.Handle(audioCurrentTime, audio.currentTime)
	ac.Handle(audioVolume, audio.volume)
}

type Audio struct {
	app.Compo
	goapp.Logging
}

func (a *Audio) Render() app.UI {
	return app.Audio().Src("").Preload("auto")
}

const EventCanPlay = "audio.event.canplay"
const EventEnded = "audio.event.ended"
const EventPause = "audio.event.pause"
const EventPlay = "audio.event.play"
const EventSeeked = "audio.event.seeked"
const EventTimeUpdate = "audio.event.timeupdate"

func (a *Audio) eventListener(ctx app.Context) app.Func {

	lastTimeUpdate := time.Now()

	return app.FuncOf(func(this app.Value, args []app.Value) any {
		if len(args) < 1 || args[0].Get("type").IsUndefined() {
			return nil
		}
		eventType := args[0].Get("type").String()

		switch eventType {
		case "canplay":
			ctx.NewAction(EventCanPlay)
		case "ended":
			ctx.NewAction(EventEnded)
		case "pause":
			ctx.NewAction(EventPause)
		case "play":
			lastTimeUpdate = time.Now().Add(-1 * time.Second)
			ctx.NewAction(EventPlay)
		case "seeked":
			lastTimeUpdate = time.Now().Add(-1 * time.Second)
			ctx.NewAction(EventSeeked)
		case "timeupdate":
			now := time.Now()
			if now.Sub(lastTimeUpdate) > time.Millisecond*200 {
				lastTimeUpdate = now
				a.timeUpdate(ctx)
			}
		}

		return nil
	})
}

func (a *Audio) timeUpdate(ctx app.Context) {
	duration := a.JSValue().Get("duration")
	if duration.IsNaN() {
		return
	}
	currentTime := a.JSValue().Get("currentTime")
	actionValue := &TimeUpdate{CurrentTime: currentTime.Float(), Duration: duration.Float()}
	ctx.NewActionWithValue(EventTimeUpdate, actionValue)
}

type TimeUpdate struct {
	CurrentTime float64
	Duration    float64
}

func (a *Audio) OnMount(ctx app.Context) {
	Action(ctx).handle(a)
	handleEvents := func(eventNames ...string) {
		listener := a.eventListener(ctx)
		for _, name := range eventNames {
			a.JSValue().Call("addEventListener", name, listener)
		}
	}
	handleEvents("canplay", "ended", "pause", "play", "seeked", "timeupdate")

	// load last volume state
	var volume float64
	ctx.GetState("volume", &volume)
	a.JSValue().Set("volume", app.ValueOf(volume))
}

func (a *Audio) src(ctx app.Context, action app.Action) {
	if url, ok := action.Value.(string); ok {
		a.Logf("src=%q", url)
		if url == "" {
			a.pause(ctx, action)
		}
		a.JSValue().Set("src", app.ValueOf(url))
	}
}

func (a *Audio) play(_ app.Context, _ app.Action) {
	a.JSValue().Call("play")
}

func (a *Audio) pause(_ app.Context, _ app.Action) {
	a.JSValue().Call("pause")
}

func (a *Audio) currentTime(ctx app.Context, action app.Action) {
	if value, ok := action.Value.(float64); ok {
		a.JSValue().Set("currentTime", app.ValueOf(value))
		a.timeUpdate(ctx)
	}
}

func (a *Audio) start(ctx app.Context, action app.Action) {
	a.src(ctx, action)
	a.play(ctx, action)
}

func (a *Audio) volume(_ app.Context, action app.Action) {
	if volume, ok := action.Value.(float64); ok {
		if volume >= 0 || volume <= 1 {
			a.JSValue().Set("volume", app.ValueOf(volume))
		}
	}
}
