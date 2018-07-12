package model

import (
	"context"
	"log"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/api"
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

// Find ...
func (sq *SliceQuery) Find(
	ctx context.Context, collection string, extra bson.M, result interface{},
) error {
	if sq.Limit <= 0 {
		return errors.New("limit must greater than 0")
	}

	queryObj, err := sq.GenQueryByObjectID()
	if err != nil {
		return err
	}
	for k, v := range extra {
		queryObj[k] = v
	}

	log.Printf("slice do query '%+v'", queryObj)
	query := api.GetMongo(ctx).C(collection).Find(queryObj).Limit(sq.Limit)
	if sq.Desc {
		query = query.Sort("-_id")
	} else {
		query = query.Sort("_id")
	}

	return query.All(result)
}
