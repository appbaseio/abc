package redis

import (
	"fmt"

	"github.com/appbaseio/abc/importer/client"
	"github.com/appbaseio/abc/importer/message"
	"github.com/appbaseio/abc/importer/message/ops"
	"github.com/appbaseio/abc/log"
)

var _ client.Reader = &Reader{}

// Reader fulfills the client.Reader interface for use with copying a Redis
// database
type Reader struct{}

func newReader() client.Reader {
	return &Reader{}
}

func (r *Reader) Read(_ map[string]client.MessageSet, filterFn client.NsFilterFunc) client.MessageChanFunc {
	return func(s client.Session, done chan struct{}) (chan client.MessageSet, error) {
		out := make(chan client.MessageSet)
		session := s.(*Session)
		go func() {
			defer close(out)
			log.With("db", session.conn.Options().DB).Infoln("starting Read func")
			sets, err := r.listCollections(session, filterFn)
			if err != nil {
				log.With("db", session.conn.Options().DB).Errorf("unable to list collections, %s", err)
				return
			}
			r.iterateCollections(session, sets, out, done)
		}()
		return out, nil
	}
}

// Scans all the sets present in a Redis database and stores them in a channel.
func (r *Reader) listCollections(redisSession *Session, filterFn func(name string) bool) (<-chan string, error) {
	out := make(chan string)
	go func() {
		var cursor uint64
		for {
			defer close(out)
			var keys []string
			var err error

			keys, cursor, err := redisSession.conn.Scan(cursor, "", 10).Result()
			if err != nil {
				log.With("db", redisSession.conn.Options().DB).Errorf("unable to scan collections, %s", err)
				return
			}

			for _, k := range keys {
				log.With("db", redisSession.conn.Options().DB).Infoln("sending for iteration...")
				out <- k
			}
			if cursor == 0 {
				log.With("db", redisSession.conn.Options().DB).Infoln("scan iteration complete...")
				break
			}
		}
	}()

	return out, nil
}

// Iterates over all the sets present in the channel returned by listCollections.
func (r *Reader) iterateCollections(redisSession *Session, in <-chan string, out chan<- client.MessageSet, done chan struct{}) error {
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return nil
			}
			log.With("db", redisSession.conn.Options().DB).Infoln("iterating...")
			keyType := redisSession.conn.Type(msg).Val()
			result, err := r.getData(redisSession, keyType, msg)
			if err != nil {
				log.With("db", redisSession.conn.Options().DB).Errorf("error while fetching data, %s", err)
			}
			out <- client.MessageSet{
				Msg: message.From(ops.Insert, msg, result),
			}
		case <-done:
			log.With("db", redisSession.conn.Options().DB).Infoln("iterating no more...")
			return nil
		}
	}
}

// Fetches the appropriate data structure associated with the key
func (r *Reader) getData(redisSession *Session, keyType string, key string) (map[string]interface{}, error) {
	switch keyType {
	case "string":
		result, err := redisSession.conn.Get(key).Result()
		fmt.Printf("Redis DB %d: Bitmaps or HLLs stored in the db might lead to non-human readable values being indexed\n", redisSession.conn.Options().DB)
		return stringMapToMessageData(map[string]string{key: result}), err
	case "hash":
		result, err := redisSession.conn.HGetAll(key).Result()
		return stringMapToMessageData(result), err
	case "list":
		result, err := redisSession.conn.LRange(key, 0, -1).Result()
		return stringSliceMapToMessageData(map[string][]string{key: result}), err
	case "set":
		result, err := redisSession.conn.SMembers(key).Result()
		return stringSliceMapToMessageData(map[string][]string{key: result}), err
	case "zset":
		result := r.getZSets(redisSession, key)
		return stringMapToMessageData(result), nil
	}
	return nil, nil
}

// Scan and iterated over the sorted sets
func (r *Reader) getZSets(redisSession *Session, key string) map[string]string {
	done := make(chan bool)
	out := make(chan string, 2)
	go func() {
		var cursor uint64

		defer func() {
			close(out)
			done <- true
		}()

		for {
			var zset []string
			var err error
			zset, cursor, err := redisSession.conn.ZScan(key, cursor, "", 10).Result()
			if err != nil {
				return
			}
			for _, z := range zset {
				out <- z
			}
			if cursor == 0 {
				break
			}
		}
	}()

	res := make(map[string]string)

loop:
	for {
		select {
		case k, ok := <-out:
			if !ok {
				done <- true
			}
			value, ok := <-out
			if !ok {
				done <- true
			}
			res[k] = value
		case <-done:
			break loop
		}
	}
	return res
}

// Converts the map[string]string map returned by the go-redis API to a map[string]interface{}
func stringMapToMessageData(setMap map[string]string) map[string]interface{} {
	m := make(map[string]interface{})
	for key, value := range setMap {
		m[key] = value
	}
	return m
}

// Converts the map[string][]string map returned by the go-redis API to a map[string]interface{}
func stringSliceMapToMessageData(iterMap map[string][]string) map[string]interface{} {
	m := make(map[string]interface{})
	for key, value := range iterMap {
		m[key] = value
	}
	return m
}
