package natsapi

import (
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-audioplayer/music"
	"log"
)

var _ model.NatsBackendApi = (*Api)(nil)

type Api struct {
	Catalog *music.Catalog
}

func (a Api) Ping(request *model.PingRequest) (response *model.PingResponse) {
	return &model.PingResponse{Query: request.Query}
}

func (a Api) Search(sr *model.SearchRequest) (response *model.SearchResponse) {
	var err error

	response, err = a.Catalog.Search("", sr)
	if err != nil {
		log.Println("error searching", err)
		response = &model.SearchResponse{}
	}

	return
}

func (a Api) Albums(request *model.AlbumsRequest) (response *model.AlbumsResponse) {
	var err error

	response, err = a.Catalog.Albums("", request)
	if err != nil {
		log.Println("error getting catalog albums", err)
		response = &model.AlbumsResponse{}
	}

	return
}

func (a Api) Album(request *model.AlbumRequest) (response *model.AlbumResponse) {
	album, err := a.Catalog.Album("", request)

	if err != nil {
		log.Println("error getting catalog album", err)
	}

	return album
}

func (a Api) PlayLists(_ *model.PlayListsRequest) (response *model.PlayListsResponse) {
	return &model.PlayListsResponse{}
}

func (a Api) RandomTrack(_ *model.RandomTrackRequest) (response *model.RandomTrackResponse) {
	return &model.RandomTrackResponse{}
}
