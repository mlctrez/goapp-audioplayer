// Package model is auto generated from model/spec.json - DO NOT EDIT
package model

import "time"

type AlbumCard struct {
	ReleaseGroupID string `json:"releasegroupid,omitempty"`
	Album          string `json:"album,omitempty"`
	Artist         string `json:"artist,omitempty"`
}

type ReleaseGroupsResponse struct {
	Groups []*ReleaseGroup `json:"release_groups,omitempty"`
}

type ReleaseGroup struct {
	ID       string          `json:"ID,omitempty"`
	CoverArt string          `json:"CoverArt,omitempty"`
	SortKey  string          `json:"SortKey,omitempty"`
	Tracks   []*ReleaseTrack `json:"Tracks,omitempty"`
}

type ReleaseTrack struct {
	ID       string    `json:"ID,omitempty"`
	Metadata *Metadata `json:"Metadata,omitempty"`
}

type Metadata struct {
	Path                      string    `json:"path,omitempty"`
	Size                      int64     `json:"size,omitempty"`
	ModTime                   time.Time `json:"mod_time,omitempty"`
	Album                     string    `json:"album,omitempty"`
	AlbumArtist               string    `json:"albumartist,omitempty"`
	Artist                    string    `json:"artist,omitempty"`
	Composer                  string    `json:"composer,omitempty"`
	Date                      string    `json:"date,omitempty"`
	DiscNumber                string    `json:"discnumber,omitempty"`
	DiscTotal                 string    `json:"disctotal,omitempty"`
	MusicbrainzAlbumArtistId  string    `json:"musicbrainz_albumartistid,omitempty"`
	MusicbrainzAlbumId        string    `json:"musicbrainz_albumid,omitempty"`
	MusicbrainzArtistId       string    `json:"musicbrainz_artistid,omitempty"`
	MusicbrainzReleaseGroupId string    `json:"musicbrainz_releasegroupid,omitempty"`
	MusicbrainzReleaseTrackId string    `json:"musicbrainz_releasetrackid,omitempty"`
	MusicbrainzTrackId        string    `json:"musicbrainz_trackid,omitempty"`
	MusicbrainzWorkId         string    `json:"musicbrainz_workid,omitempty"`
	Performer                 string    `json:"performer,omitempty"`
	Title                     string    `json:"title,omitempty"`
	TrackNumber               string    `json:"tracknumber,omitempty"`
	TrackTotal                string    `json:"tracktotal,omitempty"`
	Seconds                   uint64    `json:"seconds,omitempty"`
}

type PlayList struct {
	Name         string      `json:"name,omitempty"`
	CurrentIndex int64       `json:"current_index,omitempty"`
	Tracks       []*Metadata `json:"tracks,omitempty"`
}
