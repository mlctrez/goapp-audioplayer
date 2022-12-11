// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"fmt"
	natsgo "github.com/nats-io/nats.go"
	"time"
)

type NatsClientApi interface {
	Ping(request *PingRequest, timeout time.Duration) (response *PingResponse, err error)

	Search(request *SearchRequest, timeout time.Duration) (response *SearchResponse, err error)

	Albums(request *AlbumsRequest, timeout time.Duration) (response *AlbumsResponse, err error)

	Album(request *AlbumRequest, timeout time.Duration) (response *AlbumResponse, err error)

	PlayLists(request *PlayListsRequest, timeout time.Duration) (response *PlayListsResponse, err error)

	RandomTrack(request *RandomTrackRequest, timeout time.Duration) (response *RandomTrackResponse, err error)
}

func NewNatsClientApi(conn *natsgo.Conn) (api NatsClientApi) {
	return &natsClientApi{conn: conn}
}

type natsClientApi struct {
	conn *natsgo.Conn
}

func (na *natsClientApi) invokeNats(subject string, message WebSocketMessage, timeout time.Duration) (result WebSocketMessage, err error) {
	var bytes []byte
	if bytes, err = message.WebSocketMessage(); err != nil {
		return
	}

	var reply *natsgo.Msg
	if reply, err = na.conn.Request(subject, bytes, timeout); err != nil {
		return
	}

	return DecodeMessage(reply.Data)
}

func (na *natsClientApi) Ping(request *PingRequest, timeout time.Duration) (response *PingResponse, err error) {
	var wsm WebSocketMessage
	if wsm, err = na.invokeNats("ping", request, timeout); err != nil {
		return
	}

	var ok bool
	if response, ok = wsm.(*PingResponse); ok {
		return
	}
	return nil, fmt.Errorf("incorrect message received")
}

func (na *natsClientApi) Search(request *SearchRequest, timeout time.Duration) (response *SearchResponse, err error) {
	var wsm WebSocketMessage
	if wsm, err = na.invokeNats("search", request, timeout); err != nil {
		return
	}

	var ok bool
	if response, ok = wsm.(*SearchResponse); ok {
		return
	}
	return nil, fmt.Errorf("incorrect message received")
}

func (na *natsClientApi) Albums(request *AlbumsRequest, timeout time.Duration) (response *AlbumsResponse, err error) {
	var wsm WebSocketMessage
	if wsm, err = na.invokeNats("albums", request, timeout); err != nil {
		return
	}

	var ok bool
	if response, ok = wsm.(*AlbumsResponse); ok {
		return
	}
	return nil, fmt.Errorf("incorrect message received")
}

func (na *natsClientApi) Album(request *AlbumRequest, timeout time.Duration) (response *AlbumResponse, err error) {
	var wsm WebSocketMessage
	if wsm, err = na.invokeNats("album", request, timeout); err != nil {
		return
	}

	var ok bool
	if response, ok = wsm.(*AlbumResponse); ok {
		return
	}
	return nil, fmt.Errorf("incorrect message received")
}

func (na *natsClientApi) PlayLists(request *PlayListsRequest, timeout time.Duration) (response *PlayListsResponse, err error) {
	var wsm WebSocketMessage
	if wsm, err = na.invokeNats("playLists", request, timeout); err != nil {
		return
	}

	var ok bool
	if response, ok = wsm.(*PlayListsResponse); ok {
		return
	}
	return nil, fmt.Errorf("incorrect message received")
}

func (na *natsClientApi) RandomTrack(request *RandomTrackRequest, timeout time.Duration) (response *RandomTrackResponse, err error) {
	var wsm WebSocketMessage
	if wsm, err = na.invokeNats("randomTrack", request, timeout); err != nil {
		return
	}

	var ok bool
	if response, ok = wsm.(*RandomTrackResponse); ok {
		return
	}
	return nil, fmt.Errorf("incorrect message received")
}
