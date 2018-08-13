package model

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mw"
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

// GenQueryByTime ...
func (sq *SliceQuery) GenQueryByTime(field string) (bson.M, error) {
	if sq.Cursor == "" {
		return bson.M{}, nil
	}
	cTimeUnix, err := strconv.ParseInt(sq.Cursor, 10, 64)
	if err != nil {
		return bson.M{}, err
	}
	cTime := time.Unix(0, cTimeUnix)
	if sq.Desc {
		return bson.M{field: bson.M{"$lt": cTime}}, nil
	}
	return bson.M{field: bson.M{"$gt": cTime}}, nil
}

// Find ...
func (sq *SliceQuery) Find(
	ctx context.Context, collection, field string, queryObj bson.M, result interface{},
) error {
	if sq.Limit <= 0 {
		return errors.New("limit must greater than 0")
	}
	log.Printf("slice do query '%+v'", queryObj)
	query := mw.GetMongo(ctx).C(collection).Find(queryObj).Limit(sq.Limit)
	if sq.Desc {
		query = query.Sort(fmt.Sprintf("-%s", field))
	} else {
		query = query.Sort(field)
	}
	return query.All(result)
}
