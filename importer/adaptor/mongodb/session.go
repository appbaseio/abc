package mongodb

import (
	"github.com/appbaseio/abc/importer/client"
	"github.com/globalsign/mgo"
)

// Session serves as a wrapper for the underlying mgo.Session
type Session struct {
	mgoSession *mgo.Session
}

var _ client.Session = &Session{}

// Close implements necessary calls to cleanup the underlying mgo.Session
func (s *Session) Close() {
	s.mgoSession.Close()
}
