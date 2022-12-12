package actions

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-natsws"
	"time"
)

func RequestAlbums(ctx app.Context, natswsConn *natsws.Connection, request *model.AlbumsRequest) {
	conn, err := natswsConn.Nats()
	if err != nil {
		fmt.Println("RequestAlbums", err)
		return
	}
	var response *model.AlbumsResponse
	response, err = model.NewNatsClientApi(conn).Albums(request, time.Second)
	if err != nil {
		fmt.Println("RequestAlbums", err)
	} else {
		ctx.NewActionWithValue("model.AlbumsResponse", response)
	}
}

func RequestAlbum(ctx app.Context, natswsConn *natsws.Connection, request *model.AlbumRequest) {
	conn, err := natswsConn.Nats()
	if err != nil {
		fmt.Println("RequestAlbum", err)
		return
	}
	var response *model.AlbumResponse
	response, err = model.NewNatsClientApi(conn).Album(request, time.Second)
	if err != nil {
		fmt.Println("RequestAlbum", err)
	} else {
		ctx.NewActionWithValue("model.AlbumResponse", response)
	}
}
