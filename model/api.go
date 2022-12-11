// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"encoding/json"
	"fmt"
)

type Api interface {
	Ping(clientId string, request *PingRequest) (response *PingResponse, err error)
	Search(clientId string, request *SearchRequest) (response *SearchResponse, err error)
	Albums(clientId string, request *AlbumsRequest) (response *AlbumsResponse, err error)
	Album(clientId string, request *AlbumRequest) (response *AlbumResponse, err error)
	PlayLists(clientId string, request *PlayListsRequest) (response *PlayListsResponse, err error)
	RandomTrack(clientId string, request *RandomTrackRequest) (response *RandomTrackResponse, err error)
}

type WebSocketMessage interface {
	WebSocketMessage() ([]byte, error)
	WebSocketMessageName() string
}

func InvokeApi(clientId string, data []byte, api Api) (result []byte, err error) {
	var messageType string
	var messageJson []byte
	for i := 0; i < len(data); i++ {
		if data[i] == 0 {
			messageType = string(data[0:i])
			messageJson = data[i+1:]
			break
		}
	}

	switch messageType {
	case "PingRequest":
		request := &PingRequest{}
		if err = json.Unmarshal(messageJson, request); err != nil {
			return nil, err
		}
		var response *PingResponse
		if response, err = api.Ping(clientId, request); err != nil {
			return nil, err
		}
		return response.WebSocketMessage()
	case "SearchRequest":
		request := &SearchRequest{}
		if err = json.Unmarshal(messageJson, request); err != nil {
			return nil, err
		}
		var response *SearchResponse
		if response, err = api.Search(clientId, request); err != nil {
			return nil, err
		}
		return response.WebSocketMessage()
	case "AlbumsRequest":
		request := &AlbumsRequest{}
		if err = json.Unmarshal(messageJson, request); err != nil {
			return nil, err
		}
		var response *AlbumsResponse
		if response, err = api.Albums(clientId, request); err != nil {
			return nil, err
		}
		return response.WebSocketMessage()
	case "AlbumRequest":
		request := &AlbumRequest{}
		if err = json.Unmarshal(messageJson, request); err != nil {
			return nil, err
		}
		var response *AlbumResponse
		if response, err = api.Album(clientId, request); err != nil {
			return nil, err
		}
		return response.WebSocketMessage()
	case "PlayListsRequest":
		request := &PlayListsRequest{}
		if err = json.Unmarshal(messageJson, request); err != nil {
			return nil, err
		}
		var response *PlayListsResponse
		if response, err = api.PlayLists(clientId, request); err != nil {
			return nil, err
		}
		return response.WebSocketMessage()
	case "RandomTrackRequest":
		request := &RandomTrackRequest{}
		if err = json.Unmarshal(messageJson, request); err != nil {
			return nil, err
		}
		var response *RandomTrackResponse
		if response, err = api.RandomTrack(clientId, request); err != nil {
			return nil, err
		}
		return response.WebSocketMessage()
	}

	return nil, fmt.Errorf("message type %q not mapped", messageType)
}

func DecodeMessage(data []byte) (result WebSocketMessage, err error) {
	var messageType string
	var messageJson []byte

	for i := 0; i < len(data); i++ {
		if data[i] == 0 {
			messageType = string(data[0:i])
			messageJson = data[i+1:]
			break
		}
	}

	switch messageType {
	case "PingRequest":
		result = &PingRequest{}
	case "PingResponse":
		result = &PingResponse{}
	case "SearchRequest":
		result = &SearchRequest{}
	case "SearchResponse":
		result = &SearchResponse{}
	case "AlbumsRequest":
		result = &AlbumsRequest{}
	case "AlbumsResponse":
		result = &AlbumsResponse{}
	case "AlbumRequest":
		result = &AlbumRequest{}
	case "AlbumResponse":
		result = &AlbumResponse{}
	case "PlayListsRequest":
		result = &PlayListsRequest{}
	case "PlayListsResponse":
		result = &PlayListsResponse{}
	case "RandomTrackRequest":
		result = &RandomTrackRequest{}
	case "RandomTrackResponse":
		result = &RandomTrackResponse{}
	default:
		return nil, fmt.Errorf("message type %q not mapped", messageType)
	}
	if err = json.Unmarshal(messageJson, result); err != nil {
		return nil, err
	}
	return
}
