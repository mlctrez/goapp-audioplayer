package music

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	caa "github.com/mineo/gocaa"
	"github.com/mlctrez/bolt"
	"github.com/mlctrez/goapp-audioplayer/model"
	"go.etcd.io/bbolt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const Repo = "https://github.com/mlctrez/goapp-audioplayer"
const UserAgent = "goapp-audioplayer (" + Repo + ")"
const tWidth = 400

type Catalog struct {
	db    *bolt.Bolt
	caa   *caa.CAAClient
	index bleve.Index
}

func (c *Catalog) DB() *bolt.Bolt {
	return c.db
}

// ReleaseDiscTrack bucket stores model.Metadata keyed by releaseGroupId,
// disc number (0 pad 2), and track number (0 pad 3).
//
//	ex. 81d27c75-9680-3b3f-acc3-a6e52b864c84_01_001
const ReleaseDiscTrack = bolt.Key("release_disc_track")

// ReleaseCoverArt bucket stores the cover art image keyed by release group id or and empty byte slice if none found.
const ReleaseCoverArt = bolt.Key("release_coverArt")

func OpenCatalog(path string, readonly ...bool) (c *Catalog, err error) {
	c = &Catalog{}
	if len(readonly) > 0 {
		if c.db, err = bolt.NewWithOptions(path, &bbolt.Options{ReadOnly: true}); err != nil {
			return nil, err
		}
	} else {
		if c.db, err = bolt.New(path); err != nil {
			return nil, err
		}
	}
	c.caa = caa.NewCAAClient(UserAgent)

	buckets := bolt.Keys{ReleaseDiscTrack, ReleaseCoverArt}
	if err = c.db.CreateBuckets(buckets); err != nil {
		return
	}

	c.index, err = bleve.Open("bleve_index")
	if err != nil {
		c.CloseCatalog()
		return
	}
	return
}

func (c *Catalog) CloseCatalog() {
	if c.db != nil {
		err := c.db.Close()
		if err != nil {
			log.Println(err)
		}
	}

	if c.index != nil {
		err := c.index.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

var _ model.Api = (*Catalog)(nil)

func (c *Catalog) Search(_ string, request *model.SearchRequest) (response *model.SearchResponse, err error) {
	response = &model.SearchResponse{}

	search := bleve.NewSearchRequest(bleve.NewQueryStringQuery(request.Search))
	search.Size = 50
	searchResult, _ := c.index.Search(search)

	err = c.DB().View(func(tx *bbolt.Tx) error {

		for _, hit := range searchResult.Hits {
			md := &model.Metadata{}
			result := tx.Bucket(ReleaseDiscTrack.B()).Get([]byte(hit.ID))
			if err = json.Unmarshal(result, md); err != nil {
				return err
			}
			response.Results = append(response.Results, md)
		}

		return nil
	})

	return
}

func MetadataKey(m *model.Metadata) bolt.Key {
	return bolt.Key(fmt.Sprintf("%s_%s_%s", m.MusicbrainzReleaseGroupId, m.DiscNumber, m.TrackNumber))
}

func (c *Catalog) AddMetaData(md *model.Metadata) (err error) {
	// add / overwrite metadata

	value := &bolt.Value{K: MetadataKey(md)}

	if value.V, err = json.Marshal(md); err != nil {
		return
	}

	if err = c.db.Put(ReleaseDiscTrack, value); err != nil {
		return
	}

	if err = c.addCoverArt(md); err != nil {
		return
	}

	return nil
}

var coverArtMissingLog = make(map[string]bool)

func (c *Catalog) addCoverArt(md *model.Metadata) (err error) {

	releaseGroupId := md.MusicbrainzReleaseGroupId

	value := &bolt.Value{K: bolt.Key(releaseGroupId)}
	if err = c.db.Get(ReleaseCoverArt, value); err == bolt.ErrValueNotFound {
		log.Println("addCoverArt", releaseGroupId, md.Artist, md.Album)

		var parse uuid.UUID
		parse, err = uuid.Parse(releaseGroupId)
		if err != nil {
			return err
		}
		if value.V, err = c.getCoverArt(parse); err != nil {
			// save a marker with the error message
			value.V = []byte("ERROR: " + err.Error())
			// clear error
			err = nil
		}
		err = c.db.Put(ReleaseCoverArt, value)
	}

	if len(value.V) < 1000 {
		// try for cover art files in release directory

		var img image.Image

		imgPath := filepath.Join(filepath.Dir(md.Path), "cover.jpg")
		var imgFile *os.File
		if imgFile, err = os.Open(imgPath); err == nil {
			if img, err = jpeg.Decode(imgFile); err != nil {
				return err
			}
		}
		imgPath = filepath.Join(filepath.Dir(md.Path), "cover.png")
		if imgFile, err = os.Open(imgPath); err == nil {
			if img, err = png.Decode(imgFile); err != nil {
				return err
			}
		}
		if img != nil {
			resize := imaging.Resize(img, tWidth, 0, imaging.Lanczos)
			newData := &bytes.Buffer{}
			if err = png.Encode(newData, resize); err != nil {
				log.Println("getCoverArt Encode", err)
				return
			}
			value.V = newData.Bytes()
			err = c.db.Put(ReleaseCoverArt, value)
		} else {
			err = nil
		}
	}

	if len(value.V) < 1000 && !coverArtMissingLog[releaseGroupId] {
		coverArtMissingLog[releaseGroupId] = true
		fmt.Println(releaseGroupId, md.Artist, md.Album, string(value.V))
	}

	return
}

func (c *Catalog) getCoverArt(uu uuid.UUID) (pngBytes []byte, err error) {

	time.Sleep(5 * time.Second)

	var img caa.CoverArtImage
	if img, err = c.caa.GetReleaseGroupFront(caa.StringToUUID(uu.String()), caa.ImageSizeOriginal); err != nil {
		log.Println("getCoverArt GetReleaseGroupFront", err)
		return
	}
	var decoded image.Image
	reader := bytes.NewReader(img.Data)
	switch img.Mimetype {
	case "image/jpeg":
		decoded, err = jpeg.Decode(reader)
	case "image/png":
		decoded, err = png.Decode(reader)
	default:
		err = fmt.Errorf("unsupported Mimetype %q", img.Mimetype)
		log.Println("getCoverArt unsupported Mimetype", err)
	}

	resize := imaging.Resize(decoded, tWidth, 0, imaging.Lanczos)
	newData := &bytes.Buffer{}
	if err = png.Encode(newData, resize); err != nil {
		log.Println("getCoverArt Encode", err)
		return
	}
	pngBytes = newData.Bytes()
	return
}

func (c *Catalog) Cleanup() (err error) {

	// spin through release tracks to remove any that don't have a file
	err = c.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(ReleaseDiscTrack.B())
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			m := &model.Metadata{}
			if err = json.Unmarshal(v, m); err != nil {
				return err
			}
			if _, err = os.Stat(m.Path); os.IsNotExist(err) {
				if err = bucket.Delete(k); err != nil {
					return err
				}
			}
		}
		return nil
	})

	// spin through release tracks to remove any that don't have a file
	err = c.db.Update(func(tx *bbolt.Tx) error {

		releaseCursor := tx.Bucket(ReleaseDiscTrack.B()).Cursor()
		bucket := tx.Bucket(ReleaseCoverArt.B())
		cursor := bucket.Cursor()
		for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
			releaseCursor.First()
			toMatch, _ := releaseCursor.Seek(k)
			if !strings.HasPrefix(string(toMatch), string(k)) {
				if err = bucket.Delete(k); err != nil {
					return err
				}
			}
		}
		return nil
	})

	return
}

func (c *Catalog) CoverArt(id uuid.UUID) (pngData []byte, err error) {

	value := &bolt.Value{K: bolt.Key(id.String())}
	if err = c.db.Get(ReleaseCoverArt, value); err != nil {
		return
	}
	return value.V, nil

}

func (c *Catalog) Albums(_ string, _ *model.AlbumsRequest) (response *model.AlbumsResponse, err error) {

	// for now just return all albums
	response = &model.AlbumsResponse{}

	artistForGroupId := make(map[string]string)
	foundArtistForGroupId := make(map[string]bool)

	err = c.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket(ReleaseDiscTrack.B()).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {

			groupIdFromKey := strings.Split(string(k), "_")[0]
			if foundArtistForGroupId[groupIdFromKey] {
				continue
			}

			md := &model.Metadata{}
			if jsonErr := json.Unmarshal(v, md); jsonErr != nil {
				return jsonErr
			}
			// compare previous artist
			previousArtist := artistForGroupId[groupIdFromKey]
			if previousArtist == "" {
				artistForGroupId[groupIdFromKey] = md.Artist
				continue
			}
			albumCard := &model.AlbumCard{
				ReleaseGroupID: md.MusicbrainzReleaseGroupId,
				Album:          md.Album,
				Artist:         md.Artist,
			}
			foundArtistForGroupId[groupIdFromKey] = true
			if previousArtist != md.Artist {
				albumCard.Artist = "Various Artists"
			}

			response.Results = append(response.Results, albumCard)

		}
		return nil
	})

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(response.Results), func(i, j int) {
		response.Results[i], response.Results[j] = response.Results[j], response.Results[i]
	})

	return
}

func (c *Catalog) Album(_ string, request *model.AlbumRequest) (response *model.AlbumResponse, err error) {
	// for now just return all albums
	response = &model.AlbumResponse{ReleaseGroupID: request.ReleaseGroupID}

	err = c.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket(ReleaseDiscTrack.B()).Cursor()
		seekKey := []byte(request.ReleaseGroupID)

		for k, v := cursor.Seek(seekKey); k != nil; k, v = cursor.Next() {
			if !strings.HasPrefix(string(k), request.ReleaseGroupID) {
				break
			}

			md := &model.Metadata{}
			if jsonErr := json.Unmarshal(v, md); jsonErr != nil {
				return jsonErr
			}

			response.Tracks = append(response.Tracks, &model.ReleaseTrack{
				ID:       string(k),
				Metadata: md,
			})

		}
		return nil
	})

	return

}

func (c *Catalog) PlayLists(_ string, _ *model.PlayListsRequest) (response *model.PlayListsResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *Catalog) SetCoverArt(uu uuid.UUID, img image.Image) error {

	resize := imaging.Resize(img, tWidth, 0, imaging.Lanczos)
	newData := &bytes.Buffer{}
	if err := png.Encode(newData, resize); err != nil {
		return err
	}
	pngBytes := newData.Bytes()

	value := &bolt.Value{K: bolt.Key(uu.String()), V: pngBytes}

	return c.db.Put(ReleaseCoverArt, value)
}

func (c *Catalog) Ping(clientId string, request *model.PingRequest) (response *model.PingResponse, err error) {
	return &model.PingResponse{Query: request.Query + " " + clientId}, nil
}

func (c *Catalog) RandomTrack(_ string, _ *model.RandomTrackRequest) (response *model.RandomTrackResponse, err error) {

	var allTrackIds []string

	err = c.db.View(func(tx *bbolt.Tx) error {
		cursor := tx.Bucket(ReleaseDiscTrack.B()).Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			_ = v
			allTrackIds = append(allTrackIds, string(k))
		}
		return nil
	})

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator

	trackId := allTrackIds[r.Intn(len(allTrackIds))]
	value := &bolt.Value{K: bolt.Key(trackId)}
	if err = c.db.Get(ReleaseDiscTrack, value); err != nil {
		return
	}

	response = &model.RandomTrackResponse{Metadata: &model.Metadata{}}
	err = json.Unmarshal(value.V, response.Metadata)

	return

}
