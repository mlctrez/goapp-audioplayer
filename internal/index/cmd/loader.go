package main

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/mlctrez/bolt"
	"github.com/mlctrez/goapp-audioplayer/model"
	"github.com/mlctrez/goapp-audioplayer/music"
	"go.etcd.io/bbolt"
	"log"
)

func main() {

	db, err := bolt.NewWithOptions("bolt.db", &bbolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	var count int

	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("bleve_index", mapping)
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bbolt.Tx) error {

		cursor := tx.Bucket(music.ReleaseDiscTrack.B()).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			count++
			md := &model.Metadata{}
			err = json.Unmarshal(v, md)
			if err != nil {
				return err
			}

			err = index.Index(md.ReleaseDiscTrackID(), md)
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	err = index.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(count)

}
