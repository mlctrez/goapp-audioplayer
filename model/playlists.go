// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"bytes"
	"encoding/json"
)

type PlayListsRequest struct {
	Names []string `json:"names,omitempty"`
}

var _ WebSocketMessage = (*PlayListsRequest)(nil)

func (m *PlayListsRequest) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("PlayListsRequest")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *PlayListsRequest) WebSocketMessageName() string {
	return "PlayListsRequest"
}

type PlayListsResponse struct {
	PlayLists []*PlayList `json:"play_lists,omitempty"`
}

var _ WebSocketMessage = (*PlayListsResponse)(nil)

func (m *PlayListsResponse) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("PlayListsResponse")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *PlayListsResponse) WebSocketMessageName() string {
	return "PlayListsResponse"
}
