// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"bytes"
	"encoding/json"
)

type AlbumsRequest struct {
	Query string `json:"query,omitempty"`
}

var _ WebSocketMessage = (*AlbumsRequest)(nil)

func (m *AlbumsRequest) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("AlbumsRequest")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *AlbumsRequest) WebSocketMessageName() string {
	return "AlbumsRequest"
}

type AlbumsResponse struct {
	Results []*AlbumCard `json:"results,omitempty"`
}

var _ WebSocketMessage = (*AlbumsResponse)(nil)

func (m *AlbumsResponse) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("AlbumsResponse")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *AlbumsResponse) WebSocketMessageName() string {
	return "AlbumsResponse"
}
