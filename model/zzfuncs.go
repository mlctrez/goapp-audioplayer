package model

/*
TODO: extend spec to allow generation of these
*/

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	if d.Hours() > 1 {
		return fmt.Sprintf("%2d:%02d:%02d", int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60)
	}
	return fmt.Sprintf("%2d:%02d", int(d.Minutes()), int(d.Seconds())%60)
}

func (md *Metadata) FlacUrl() string {
	return "/flac/" + md.ReleaseDiscTrackID()
}

func (md *Metadata) ReleaseDiscTrackID() string {
	return fmt.Sprintf("%s_%s_%s", md.MusicbrainzReleaseGroupId, md.DiscNumber, md.TrackNumber)
}

func (md *Metadata) FormattedDuration() string {
	return FormatDuration(time.Second * time.Duration(md.Seconds))
}

func (m *AlbumResponse) TracksMetadata() []*Metadata {
	var allTracks []*Metadata
	for _, track := range m.Tracks {
		allTracks = append(allTracks, track.Metadata)
	}
	return allTracks
}
