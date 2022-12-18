package search

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/model"
)

type Card struct {
	app.Compo
	md *model.Metadata
}

func (c *Card) OnMount(ctx app.Context) {
}

func (c *Card) Render() app.UI {
	src := model.CoverArtUrl(c.md.MusicbrainzReleaseGroupId, 0)
	return app.Tr().OnClick(c.rowClick).ID(c.md.ReleaseDiscTrackID()).Class("search-row-tr").Body(
		app.Td().Class("search-image-td").Body(
			app.Img().Alt(c.md.Album).Src(src).Width(48).Height(48),
		),
		app.Td().Class("search-album-td").Text(c.md.Album),
		app.Td().Class("search-artist-td").Text(c.md.Artist),
		app.Td().Class("search-title-td").Text(c.md.Title),
	)
}

func (c *Card) rowClick(ctx app.Context, e app.Event) {
	ctx.NewActionWithValue("queue.add", c.md)
}
