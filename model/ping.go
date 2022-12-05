// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"bytes"
	"encoding/json"
)

type PingRequest struct {
	Query string `json:"query,omitempty"`
}

var _ WebSocketMessage = (*PingRequest)(nil)

func (m *PingRequest) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("PingRequest")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *PingRequest) WebSocketMessageName() string {
	return "PingRequest"
}

type PingResponse struct {
	Query string `json:"query,omitempty"`
}

var _ WebSocketMessage = (*PingResponse)(nil)

func (m *PingResponse) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("PingResponse")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *PingResponse) WebSocketMessageName() string {
	return "PingResponse"
}
