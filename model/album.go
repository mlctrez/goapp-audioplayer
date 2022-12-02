// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"bytes"
	"encoding/json"
)

type AlbumRequest struct {
	ReleaseGroupID string `json:"releasegroupid,omitempty"`
}

var _ WebSocketMessage = (*AlbumRequest)(nil)

func (m *AlbumRequest) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("AlbumRequest")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *AlbumRequest) WebSocketMessageName() string {
	return "AlbumRequest"
}

type AlbumResponse struct {
	ReleaseGroupID string          `json:"releasegroupid,omitempty"`
	Tracks         []*ReleaseTrack `json:"tracks,omitempty"`
}

var _ WebSocketMessage = (*AlbumResponse)(nil)

func (m *AlbumResponse) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("AlbumResponse")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *AlbumResponse) WebSocketMessageName() string {
	return "AlbumResponse"
}
