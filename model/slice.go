package model

import (
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
)

// SliceInfo ...
type SliceInfo struct {
	FirstCursor string
	LastCursor  string
}

// SliceQuery ...
type SliceQuery struct {
	Limit  int
	Desc   bool
	Cursor string
}

// GenQueryByObjectID ...
func (sq *SliceQuery) GenQueryByObjectID() (bson.M, error) {
	if sq.Cursor == "" {
		return bson.M{}, nil
	}
	if !bson.IsObjectIdHex(sq.Cursor) {
		return nil, errors.New("Invalid cursor")
	}
	id := bson.ObjectIdHex(sq.Cursor)
	if sq.Desc {
		return bson.M{"_id": bson.M{"$lt": id}}, nil
	}
	return bson.M{"_id": bson.M{"$gt": id}}, nil
}
