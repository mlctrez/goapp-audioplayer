package album

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/nodisplay"
	"github.com/mlctrez/goapp-audioplayer/goapp/compo/websocket"
	"github.com/mlctrez/goapp-audioplayer/internal/icon"
	"github.com/mlctrez/goapp-audioplayer/model"
	"strings"
	"time"
)

type Album struct {
	app.Compo
	album       *model.AlbumResponse
	displayMode string
}

func (t *Album) OnMount(ctx app.Context) {
	ctx.ObserveState("displayMode").Value(&t.displayMode)
	websocket.Action(ctx).HandleAction(&model.AlbumResponse{}, t)
}

func (t *Album) OnWebsocketMessage(ctx app.Context, message model.WebSocketMessage) {
	t.album = message.(*model.AlbumResponse)
	ctx.SetState("displayMode", "album.Album")
}

func (t *Album) Render() app.UI {

	if t.album == nil || t.displayMode != "album.Album" {
		return nodisplay.NoDisplay("compo/album/List")
	}

	var albumTitle string
	var albumArtist string
	var albumDate string
	var songs int
	var totalDuration time.Duration

	for _, track := range t.album.Tracks {
		totalDuration += time.Second * time.Duration(track.Metadata.Seconds)
		songs++
		albumTitle = track.Metadata.Album
		if albumArtist == "" {
			albumArtist = track.Metadata.Artist
		} else {
			if albumArtist != track.Metadata.Artist {
				albumArtist = "Various Artists"
			}
		}
		if albumDate == "" && track.Metadata.Date != "" && len(track.Metadata.Date) > 3 {
			albumDate = " • " + track.Metadata.Date[0:4]
		}
	}

	var rows []app.UI

	rows = append(rows,
		app.Tr().Body(
			app.Td().ColSpan(2).Body(
				app.Div().Style("display", "flex").Style("flex-direction", "column").Body(
					app.H3().Text(albumTitle),
					app.Div().Text("Album • "+albumArtist+albumDate),
					app.Div().Text(fmt.Sprintf("%d songs • %d minutes", songs, int(totalDuration.Minutes()))),
				),
			),
			app.Td().Body(
				app.Div().Style("display", "flex").Style("flex-direction", "column").Body(
					app.Div().Style("position", "relative").Style("left", "-10px").
						Body(app.Raw(icon.Close48())).OnClick(func(ctx app.Context, e app.Event) {

						t.album = nil
						ctx.SetState("displayMode", "album.List")

					}),
					app.Div().Style("position", "relative").Style("left", "-10px").
						Body(app.Raw(icon.PlaylistAdd48())).OnClick(func(ctx app.Context, e app.Event) {

						ctx.NewActionWithValue("queue.add", t.album.TracksMetadata())

					}),
				),
			),
		),
		app.Tr().Body(app.Td().ColSpan(3).Style("height", "20px")),
	)

	for _, loopTrack := range t.album.Tracks {
		track := loopTrack
		rows = append(rows, &TrackRow{Metadata: track.Metadata})
	}

	table := app.Table().Style("width", "600px").Body(rows...)

	image := app.Img().Src(model.CoverArtUrl(t.album.ReleaseGroupID, 0))

	return app.Div().Class("main-content").Body(
		app.Div().Class("main-content-large-image").Body(image),
		app.Div().Class("main-content-album-tracks").Body(table),
	)
}

type TrackRow struct {
	app.Compo
	Metadata *model.Metadata
}

func (tr *TrackRow) Render() app.UI {
	md := tr.Metadata

	return app.Tr().Class("album-track-row").ID(md.ReleaseDiscTrackID()).Body(
		app.Td().Class("album-track-number").Text(strings.TrimLeft(md.TrackNumber, "0")),
		app.Td().Class("album-track-title").Text(md.Title),
		app.Td().Class("album-track-duration").Text(md.FormattedDuration()),
	).OnClick(tr.click)
}

func (tr *TrackRow) click(ctx app.Context, _ app.Event) {
	ctx.NewActionWithValue("queue.add", tr.Metadata)
}
