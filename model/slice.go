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

// GenQueryByTime in milliSeconds
func (sq *SliceQuery) GenQueryByTime(field string) (bson.M, error) {
	if sq.Cursor == "" {
		return bson.M{}, nil
	}
	cursorTime, err := parseTimeCursor(sq.Cursor)
	if err != nil {
		return bson.M{}, err
	}
	if sq.Desc {
		return bson.M{field: bson.M{"$lt": cursorTime}}, nil
	}
	return bson.M{field: bson.M{"$gt": cursorTime}}, nil
}

func genTimeCursor(t time.Time) string {
	return fmt.Sprint(t.UnixNano() / 1000 / 1000)
}

func parseTimeCursor(s string) (time.Time, error) {
	var t time.Time
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return t, err
	}
	t = time.Unix(0, ms*1000*1000)
	return t, nil
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
