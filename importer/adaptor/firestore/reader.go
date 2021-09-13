package firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/message"
	"github.com/appbaseio/abc/importer/message/ops"
	"github.com/appbaseio/abc/log"
	"google.golang.org/api/iterator"
)

var _ client.Reader = &Reader{}

type Reader struct{}

func (r *Reader) listCollections(fc firestore.Client, collectionFilterFn func(collectionID string) bool) (<-chan *firestore.CollectionRef, error) {
	out := make(chan *firestore.CollectionRef)

	go func() {
		defer close(out)
		collections := fc.Collections(context.Background())
		for {
			collectionRef, err := collections.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				continue
			}
			id, path := collectionRef.ID, collectionRef.Path
			if collectionFilterFn(id) {
				log.With("collectionID", id).With("collectionPath", path).Infoln("sending for iteration ...")
				out <- collectionRef
			} else {
				log.With("collectionID", id).Infoln("skipping collection iteration ...")
			}
		}
	}()

	return out, nil
}

func (r *Reader) Read(_ map[string]client.MessageSet, filterFn client.NsFilterFunc) client.MessageChanFunc {
	return func(s client.Session, done chan struct{}) (chan client.MessageSet, error) {
		out := make(chan client.MessageSet)
		session := s.(*Session)
		fc := *session.fc

		collectionRef, err := r.listCollections(fc, filterFn)
		if err != nil {
			return nil, err
		}

		go func() {
			defer close(out)
			for {
				collection, ok := <-collectionRef
				if !ok {
					break
				} else {
					docSnapshots := collection.Documents(context.Background())
					for {
						docSnapshot, err := docSnapshots.Next()
						if err == iterator.Done {
							break
						}
						if err != nil {
							continue
						}
						out <- client.MessageSet{
							Msg: message.From(ops.Insert, docSnapshot.Ref.Parent.ID, docSnapshot.Data()),
						}
					}
				}
			}
		}()

		return out, nil
	}
}
