package main

import (
	"encoding/json"
	"fmt"
	"github.com/dhowden/tag"
	"github.com/mewkiz/flac"
	caa "github.com/mineo/gocaa"
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-audioplayer/music"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Index struct {
	c   *music.Catalog
	caa *caa.CAAClient
}

func (i *Index) processPath(path string, info fs.FileInfo) (err error) {

	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	var tagData tag.Metadata
	if tagData, err = tag.ReadFrom(f); err != nil {
		return
	}
	var md *model.Metadata
	if md, err = MetadataFromRaw(tagData.Raw()); err != nil {
		return
	}
	// back to beginning of file
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return
	}
	var stream *flac.Stream
	if stream, err = flac.New(f); err != nil {
		return
	}

	md.Path = path
	md.Size = info.Size()
	md.ModTime = info.ModTime()
	md.Seconds = stream.Info.NSamples / uint64(stream.Info.SampleRate)

	return i.c.AddMetaData(md)
}

func (i *Index) walkFunc(path string, info fs.FileInfo, err error) error {

	if err != nil || info.IsDir() {
		return err
	}
	if strings.HasSuffix(path, ".flac") && !strings.HasPrefix(info.Name(), "00-") {
		return i.processPath(path, info)
	}
	return nil
}

func (i *Index) walk() error {
	return filepath.Walk(filepath.Join(userHome(), "Music"), i.walkFunc)
}

func (i *Index) run() (err error) {
	if i.c, err = music.OpenCatalog("bolt.db"); err != nil {
		return err
	}
	defer i.c.CloseCatalog()

	// this puts data into the catalog
	if err = i.walk(); err != nil {
		return
	}

	return i.c.Cleanup()

}

func main() {
	err := (&Index{}).run()
	if err != nil {
		log.Fatal(err)
	}
}

func userHome() (path string) {
	var err error
	if path, err = os.UserHomeDir(); err != nil {
		panic(err)
	}
	return
}

func MetadataFromRaw(rawMap map[string]interface{}) (data *model.Metadata, err error) {

	var marshal []byte
	if marshal, err = json.Marshal(rawMap); err != nil {
		return
	}

	data = &model.Metadata{}
	if err = json.Unmarshal(marshal, data); err != nil {
		return
	}

	if !ValidMetadata(data) {
		err = fmt.Errorf("bad metadata %+v", rawMap)
	}

	data.TrackNumber = padThree(data.TrackNumber)
	data.TrackTotal = padThree(data.TrackTotal)

	data.DiscNumber = padTwo(data.DiscNumber)
	data.DiscTotal = padTwo(data.DiscTotal)

	return
}

func ValidMetadata(m *model.Metadata) bool {

	return m.Album != "" &&
		m.Title != "" &&
		m.Artist != "" &&
		m.MusicbrainzAlbumId != "" &&
		m.MusicbrainzTrackId != "" &&
		m.MusicbrainzArtistId != "" &&
		m.MusicbrainzReleaseGroupId != "" &&
		m.MusicbrainzReleaseTrackId != "" &&
		m.TrackNumber != "" &&
		m.TrackTotal != ""

}

func padTwo(in string) string {
	return leftPad(in, "0", 2)
}

func padThree(in string) string {
	return leftPad(in, "0", 3)
}

func leftPad(in, padWith string, length int) (result string) {
	result = in
	for len(result) < length {
		result = padWith + result
	}
	return
}
