package music

import (
	"encoding/json"
	"github.com/mlctrez/bolt"
	"github.com/mlctrez/goapp-audioplayer/model"
)

func (c *Catalog) ReleaseDiscTrackMetadata(key string) (md *model.Metadata, err error) {

	value := &bolt.Value{K: bolt.Key(key)}
	if err = c.db.Get(ReleaseDiscTrack, value); err != nil {
		return
	}
	md = &model.Metadata{}
	err = json.Unmarshal(value.V, md)

	return
}
