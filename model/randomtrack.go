// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"bytes"
	"encoding/json"
)

type RandomTrackRequest struct{}

var _ WebSocketMessage = (*RandomTrackRequest)(nil)

func (m *RandomTrackRequest) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("RandomTrackRequest")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *RandomTrackRequest) WebSocketMessageName() string {
	return "RandomTrackRequest"
}

type RandomTrackResponse struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

var _ WebSocketMessage = (*RandomTrackResponse)(nil)

func (m *RandomTrackResponse) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("RandomTrackResponse")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *RandomTrackResponse) WebSocketMessageName() string {
	return "RandomTrackResponse"
}
