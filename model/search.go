// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import (
	"bytes"
	"encoding/json"
)

type SearchRequest struct {
	Search string `json:"search,omitempty"`
}

var _ WebSocketMessage = (*SearchRequest)(nil)

func (m *SearchRequest) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("SearchRequest")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *SearchRequest) WebSocketMessageName() string {
	return "SearchRequest"
}

type SearchResponse struct {
	Groups []*ReleaseGroup `json:"release_groups,omitempty"`
}

var _ WebSocketMessage = (*SearchResponse)(nil)

func (m *SearchResponse) WebSocketMessage() ([]byte, error) {
	result := bytes.NewBufferString("SearchResponse")
	result.WriteByte(0)

	err := json.NewEncoder(result).Encode(m)
	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func (m *SearchResponse) WebSocketMessageName() string {
	return "SearchResponse"
}
