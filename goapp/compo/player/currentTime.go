package player

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"math"
)

type CurrentTime struct {
	app.Compo
	current  float64
	duration float64
}

func (t *CurrentTime) Render() app.UI {
	timeUI := app.Div().Class("audio-time")
	if t.duration > 1 {
		timeUI.Text(timeFormat(t.current) + " / " + timeFormat(t.duration))
	} else {
		timeUI.Text("----:---- / ----:----")
	}
	return timeUI
}

func timeFormat(seconds float64) string {
	sec := int(math.RoundToEven(seconds))
	return fmt.Sprintf("%02d:%02d", sec/60, sec%60)
}

var _ app.Mounter = (*CurrentTime)(nil)

func (t *CurrentTime) OnMount(ctx app.Context) {
	ctx.Handle(audio.EventTimeUpdate, func(context app.Context, action app.Action) {
		if tu, ok := action.Value.(*audio.TimeUpdate); ok {
			t.current = tu.CurrentTime
			t.duration = tu.Duration
			t.Update()
		}
	})
}
