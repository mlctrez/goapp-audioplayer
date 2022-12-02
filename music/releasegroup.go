package music

//func (c *Catalog) ReleaseGroup(uu uuid.UUID) (group *model.ReleaseGroup, err error) {
//	err = c.db.View(func(tx *bbolt.Tx) error {
//		trackCursor := tx.Bucket([]byte(ReleaseDiscTrack)).Cursor()
//
//		var tracks []*model.ReleaseTrack
//
//		seek := bolt.Key(uu.String())
//		for k, v := trackCursor.Seek(seek.B()); k != nil; k, v = trackCursor.Next() {
//			if !strings.HasPrefix(string(k), uu.String()) {
//				break
//			}
//
//			md := &model.Metadata{}
//			if jErr := json.Unmarshal(v, md); jErr != nil {
//				return jErr
//			}
//			tracks = append(tracks, &model.ReleaseTrack{ID: string(k), Metadata: md})
//		}
//		if len(tracks) > 0 {
//			group = &model.ReleaseGroup{ID: uu.String(), Tracks: tracks}
//		}
//
//		return nil
//	})
//
//	return
//}

//func (c *Catalog) ReleaseGroups() (groups []*model.ReleaseGroup, err error) {
//	err = c.db.View(func(tx *bbolt.Tx) error {
//		covers := tx.Bucket([]byte(ReleaseCoverArt)).Cursor()
//
//		releases := tx.Bucket([]byte(ReleaseDiscTrack)).Cursor()
//
//		for k, v := covers.First(); k != nil; k, v = covers.Next() {
//			_ = v // ignored, just looking for keys
//			releaseGroup := &model.ReleaseGroup{ID: string(k)}
//
//			releases.First()
//			_, value := releases.Seek([]byte(releaseGroup.ID))
//			md := &model.Metadata{}
//			if err = json.Unmarshal(value, md); err != nil {
//				return err
//			}
//			releaseGroup.SortKey = md.Artist + md.Album
//
//			groups = append(groups, releaseGroup)
//		}
//		return nil
//	})
//
//	sort.SliceStable(groups, func(i, j int) bool {
//		return groups[i].SortKey < groups[j].SortKey
//	})
//
//	return
//}
