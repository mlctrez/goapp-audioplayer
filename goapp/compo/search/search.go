package search

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-natsws"
	"time"
)

type Search struct {
	app.Compo
	displayMode string
	natswsConn  *natsws.Connection
	results     []*model.Metadata
}

func (s *Search) Render() app.UI {

	if s.displayMode != "search.Search" {
		return nodisplay.NoDisplay("compo/search/Search")
	}

	var results []app.UI

	if s.results == nil {
		results = append(results, app.Div().Text(""))
	} else {
		for _, result := range s.results {
			results = append(results, &Card{md: result})
		}
	}

	return app.Div().Class("main-content").Body(
		app.Table().Class("search-table").Body(results...),
	)
}

func (s *Search) OnMount(ctx app.Context) {
	ctx.ObserveState("displayMode").Value(&s.displayMode)
	s.natswsConn = &natsws.Connection{}
	natsws.Observe(ctx, s.natswsConn)
	ctx.Handle("search.input", s.searchInput)
}

func (s *Search) searchInput(ctx app.Context, action app.Action) {

	if ss, ok := action.Value.(string); ok {
		conn, err := s.natswsConn.Nats()
		if err != nil {
			fmt.Println("error getting nats connection", err)
			return
		}
		response, err := model.NewNatsClientApi(conn).Search(&model.SearchRequest{Search: ss}, time.Second*3)
		if err != nil {
			fmt.Println("error getting response", err)
			return
		}
		s.results = response.Results
	}

}
