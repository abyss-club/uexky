package mw

import (
	"context"
	"log"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
)

// Mongo ...
type Mongo struct {
	session *mgo.Session
}

// ConnectMongodb ...
func ConnectMongodb() *Mongo {
	session, err := mgo.Dial(mgmt.Config.Mongo.URI)
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
	return m.session.DB(mgmt.Config.Mongo.DB)
}

// C return collection
func (m *Mongo) C(name string) *mgo.Collection {
	return m.session.DB(mgmt.Config.Mongo.DB).C(name)
}

// WithMongo set mongodb session to context
func WithMongo(handle httprouter.Handle) httprouter.Handle {
	mongo := ConnectMongodb()
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		m := mongo.Copy()
		defer m.Close()
		req = req.WithContext(context.WithValue(
			req.Context(), ContextKeyMongo, m))
		handle(w, req, p)
	}
}

// GetMongo from contenxt
func GetMongo(ctx context.Context) *Mongo {
	m, ok := ctx.Value(ContextKeyMongo).(*Mongo)
	if !ok {
		log.Fatal("Can't find mongodb in context")
	}
	return m
}
