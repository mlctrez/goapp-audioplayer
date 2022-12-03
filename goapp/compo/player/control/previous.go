package control

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/queue"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
)

type Previous struct {
	app.Compo
	queue queue.Queue
}

func (p *Previous) OnMount(ctx app.Context) {
	ctx.ObserveState("queue").Value(&p.queue)
}

func (p *Previous) Render() app.UI {
	enabled := p.queue.HasPrevious()
	svg := app.Raw(icon.SkipPrevious48())

	return Div(enabled).Body(svg).OnClick(p.click)
}

func (p *Previous) click(ctx app.Context, _ app.Event) {
	p.queue.Previous(ctx)
}
