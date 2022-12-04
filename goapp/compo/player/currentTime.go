package player

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/audio"
	"math"
)

type CurrentTime struct {
	app.Compo
	timeUpdate *audio.TimeUpdate
}

func (t *CurrentTime) Render() app.UI {
	timeUI := app.Div().Class("audio-time")
	if t.timeUpdate != nil {
		timeUI.Text(timeFormat(t.timeUpdate.CurrentTime) + " / " + timeFormat(t.timeUpdate.Duration))
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
	ctx.Handle(audio.EventTimeUpdate, t.update)
}

func (t *CurrentTime) update(_ app.Context, action app.Action) {
	if tu, ok := action.Value.(*audio.TimeUpdate); ok {
		t.timeUpdate = tu
	}
}
