package control

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/queue"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
)

type Next struct {
	app.Compo
	queue queue.Queue
}

func (n *Next) OnMount(ctx app.Context) {
	ctx.ObserveState("queue").Value(&n.queue)
}

func (n *Next) Render() app.UI {
	enabled := n.queue.HasNext()
	svg := app.Raw(icon.SkipNext48())

	return Div(enabled).Body(svg).OnClick(n.click)
}

func (n *Next) click(ctx app.Context, _ app.Event) {
	n.queue.Next(ctx)
}
