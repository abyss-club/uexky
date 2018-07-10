package api

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

// Dial ...
func (m *Mongo) Dial() error {
	session, err := mgo.Dial(mgmt.Config.Mongo.URI)
	if err != nil {
		return errors.Wrap(err, "connect to mongodb")
	}
	m.session = session
	return nil
}

// Copy ...
func (m *Mongo) Copy() *Mongo {
	return &Mongo{session: m.session.Copy()}
}

// Close ...
func (m *Mongo) Close() {
	m.session.Close()
}

// C return collection
func (m *Mongo) C(name string) *mgo.Collection {
	return m.session.DB(mgmt.Config.Mongo.DB).C(name)
}

// WithMongo set mongodb session to context
func WithMongo(handle httprouter.Handle, mongo *Mongo) httprouter.Handle {
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
