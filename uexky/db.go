package uexky

import (
	"log"

	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky-go/config"
)

// Mongo ...
type Mongo struct {
	session *mgo.Session
}

// ConnectMongodb ...
func ConnectMongodb() *Mongo {
	session, err := mgo.Dial(config.Config.Mongo.URL)
	if err != nil {
		log.Fatal(errors.Wrap(err, "connect to mongodb"))
	}
	return &Mongo{session: session}
}

// Copy ...
func (m *Mongo) Copy() *Mongo {
	return &Mongo{session: m.session.Copy()}
}

// Close ...
func (m *Mongo) Close() {
	m.session.Close()
}

// DB ...
func (m *Mongo) DB() *mgo.Database {
	return m.session.DB(config.Config.Mongo.DB)
}

// C return collection
func (m *Mongo) C(name string) *mgo.Collection {
	return m.session.DB(config.Config.Mongo.DB).C(name)
}
