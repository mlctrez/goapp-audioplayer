// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"context"
	"fmt"
	natsgo "github.com/nats-io/nats.go"
)

type NatsBackendApi interface {
	Ping(request *PingRequest) (response *PingResponse)

	Search(request *SearchRequest) (response *SearchResponse)

	Albums(request *AlbumsRequest) (response *AlbumsResponse)

	Album(request *AlbumRequest) (response *AlbumResponse)

	PlayLists(request *PlayListsRequest) (response *PlayListsResponse)

	RandomTrack(request *RandomTrackRequest) (response *RandomTrackResponse)
}

type NatsBackend interface {
	Start() error

	Stop()
}

func NewNatsBackend(ctx context.Context, conn *natsgo.Conn, api NatsBackendApi) (backend NatsBackend) {
	b := &natsBackend{}

	b.ctx = ctx
	b.conn = conn
	b.api = api

	return b
}

var _ NatsBackend = (*natsBackend)(nil)

type natsBackend struct {
	ctx       context.Context
	conn      *natsgo.Conn
	api       NatsBackendApi
	subCtx    context.Context
	subCancel context.CancelFunc
}

func (nb *natsBackend) Start() (err error) {
	nb.subCtx, nb.subCancel = context.WithCancel(nb.ctx)
	var subs []*natsgo.Subscription
	errUnsubAll := func() {
		for _, sub := range subs {
			_ = sub.Unsubscribe()
		}
	}

	pingChan := make(chan *natsgo.Msg, 10)
	var pingSub *natsgo.Subscription
	if pingSub, err = nb.conn.ChanSubscribe("ping", pingChan); err != nil {
		errUnsubAll()
		return
	}
	subs = append(subs, pingSub)

	searchChan := make(chan *natsgo.Msg, 10)
	var searchSub *natsgo.Subscription
	if searchSub, err = nb.conn.ChanSubscribe("search", searchChan); err != nil {
		errUnsubAll()
		return
	}
	subs = append(subs, searchSub)

	albumsChan := make(chan *natsgo.Msg, 10)
	var albumsSub *natsgo.Subscription
	if albumsSub, err = nb.conn.ChanSubscribe("albums", albumsChan); err != nil {
		errUnsubAll()
		return
	}
	subs = append(subs, albumsSub)

	albumChan := make(chan *natsgo.Msg, 10)
	var albumSub *natsgo.Subscription
	if albumSub, err = nb.conn.ChanSubscribe("album", albumChan); err != nil {
		errUnsubAll()
		return
	}
	subs = append(subs, albumSub)

	playListsChan := make(chan *natsgo.Msg, 10)
	var playListsSub *natsgo.Subscription
	if playListsSub, err = nb.conn.ChanSubscribe("playLists", playListsChan); err != nil {
		errUnsubAll()
		return
	}
	subs = append(subs, playListsSub)

	randomTrackChan := make(chan *natsgo.Msg, 10)
	var randomTrackSub *natsgo.Subscription
	if randomTrackSub, err = nb.conn.ChanSubscribe("randomTrack", randomTrackChan); err != nil {
		errUnsubAll()
		return
	}
	subs = append(subs, randomTrackSub)

	go func() {
		defer errUnsubAll()
		for {
			select {
			case <-nb.subCtx.Done():
				return
			case msg := <-pingChan:
				if msg == nil {
					continue
				}

				var response WebSocketMessage
				if response, err = DecodeMessage(msg.Data); err != nil {
					fmt.Println("bad message received", err)
				}

				if wsm, ok := response.(*PingRequest); ok {
					resp := nb.api.Ping(wsm)

					var message []byte
					if message, err = resp.WebSocketMessage(); err != nil {
						fmt.Println("bad message from api", err)
					}

					if err = nb.conn.Publish(msg.Reply, message); err != nil {
						fmt.Println("cannot publish reply", err)
					}
				}
			case msg := <-searchChan:
				if msg == nil {
					continue
				}

				var response WebSocketMessage
				if response, err = DecodeMessage(msg.Data); err != nil {
					fmt.Println("bad message received", err)
				}

				if wsm, ok := response.(*SearchRequest); ok {
					resp := nb.api.Search(wsm)

					var message []byte
					if message, err = resp.WebSocketMessage(); err != nil {
						fmt.Println("bad message from api", err)
					}

					if err = nb.conn.Publish(msg.Reply, message); err != nil {
						fmt.Println("cannot publish reply", err)
					}
				}
			case msg := <-albumsChan:
				if msg == nil {
					continue
				}

				var response WebSocketMessage
				if response, err = DecodeMessage(msg.Data); err != nil {
					fmt.Println("bad message received", err)
				}

				if wsm, ok := response.(*AlbumsRequest); ok {
					resp := nb.api.Albums(wsm)

					var message []byte
					if message, err = resp.WebSocketMessage(); err != nil {
						fmt.Println("bad message from api", err)
					}

					if err = nb.conn.Publish(msg.Reply, message); err != nil {
						fmt.Println("cannot publish reply", err)
					}
				}
			case msg := <-albumChan:
				if msg == nil {
					continue
				}

				var response WebSocketMessage
				if response, err = DecodeMessage(msg.Data); err != nil {
					fmt.Println("bad message received", err)
				}

				if wsm, ok := response.(*AlbumRequest); ok {
					resp := nb.api.Album(wsm)

					var message []byte
					if message, err = resp.WebSocketMessage(); err != nil {
						fmt.Println("bad message from api", err)
					}

					if err = nb.conn.Publish(msg.Reply, message); err != nil {
						fmt.Println("cannot publish reply", err)
					}
				}
			case msg := <-playListsChan:
				if msg == nil {
					continue
				}

				var response WebSocketMessage
				if response, err = DecodeMessage(msg.Data); err != nil {
					fmt.Println("bad message received", err)
				}

				if wsm, ok := response.(*PlayListsRequest); ok {
					resp := nb.api.PlayLists(wsm)

					var message []byte
					if message, err = resp.WebSocketMessage(); err != nil {
						fmt.Println("bad message from api", err)
					}

					if err = nb.conn.Publish(msg.Reply, message); err != nil {
						fmt.Println("cannot publish reply", err)
					}
				}
			case msg := <-randomTrackChan:
				if msg == nil {
					continue
				}

				var response WebSocketMessage
				if response, err = DecodeMessage(msg.Data); err != nil {
					fmt.Println("bad message received", err)
				}

				if wsm, ok := response.(*RandomTrackRequest); ok {
					resp := nb.api.RandomTrack(wsm)

					var message []byte
					if message, err = resp.WebSocketMessage(); err != nil {
						fmt.Println("bad message from api", err)
					}

					if err = nb.conn.Publish(msg.Reply, message); err != nil {
						fmt.Println("cannot publish reply", err)
					}
				}
			}
		}
	}()

	return
}

func (nb *natsBackend) Stop() {
	if nb.subCancel != nil {
		nb.subCancel()
	}
}
