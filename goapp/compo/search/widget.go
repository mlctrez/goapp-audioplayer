package search

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
)

type Widget struct {
	app.Compo
	displayMode string
}

func (w *Widget) Render() app.UI {
	div := app.Div().Class("navigation-search")
	if w.displayMode == "search.Search" {
		div.Body(&Input{}, app.Raw(icon.Search48()))
	} else {
		div.Body(app.Raw(icon.Search48())).OnClick(w.click)
	}
	return div
}

func (w *Widget) OnMount(ctx app.Context) {
	ctx.ObserveState("displayMode").Value(&w.displayMode)
}

func (w *Widget) click(ctx app.Context, e app.Event) {
	ctx.SetState("displayMode", "search.Search")
}

type Input struct {
	app.Compo
	input string
}

func (i *Input) Render() app.UI {
	return app.Input().ID("search-input").Value(i.input).OnInput(i.handleInput)
}

func (i *Input) handleInput(ctx app.Context, e app.Event) {
	i.ValueTo(&i.input)(ctx, e)
	ctx.NewActionWithValue("search.input", i.input)
}

func (i *Input) OnMount(ctx app.Context) {
	searchInput := app.Window().GetElementByID("search-input")
	if searchInput.Truthy() {
		searchInput.Call("focus")
	}
}
