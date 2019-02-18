package model

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky-go/uexky"
)

// SliceInfo ...
type SliceInfo struct {
	FirstCursor string
	LastCursor  string
	HasPrev     bool
	HasNext     bool
}

// SliceQuery ...
type SliceQuery struct {
	GT    string
	LT    string
	Limit int
}

// GenQueryByObjectID ...
func (sq *SliceQuery) GenQueryByObjectID() (bson.M, error) {
	qry := bson.M{}
	if sq.GT != "" {
		id, err := parseObjectID(sq.GT)
		if err != nil {
			return nil, err
		}
		qry["$gt"] = id
	}
	if sq.LT != "" {
		id, err := parseObjectID(sq.LT)
		if err != nil {
			return nil, err
		}
		qry["$lt"] = id
	}
	if len(qry) == 0 {
		return bson.M{}, nil
	}
	return bson.M{"_id": qry}, nil
}

// GenQueryByTime in milliSeconds
func (sq *SliceQuery) GenQueryByTime(field string) (bson.M, error) {
	qry := bson.M{}
	if sq.GT != "" {
		time, err := parseTimeCursor(sq.GT)
		if err != nil {
			return nil, err
		}
		qry["$gt"] = time
	}
	if sq.LT != "" {
		time, err := parseTimeCursor(sq.LT)
		if err != nil {
			return nil, err
		}
		qry["$lt"] = time
	}
	if len(qry) == 0 {
		return bson.M{}, nil
	}
	return bson.M{field: qry}, nil
}

// Find ...
func (sq *SliceQuery) Find(
	u *uexky.Uexky, collection, sort string,
	queryObj bson.M, result interface{},
) error {
	if sq.Limit <= 0 {
		return errors.New("limit must greater than 0")
	}
	if err := u.Flow.CostQuery(sq.Limit); err != nil {
		return err
	}
	log.Printf("slice do query '%+v'", queryObj)
	query := u.Mongo.C(collection).Find(queryObj).Limit(sq.Limit).Sort(sort)
	return query.All(result)
}

func genTimeCursor(t time.Time) string {
	return fmt.Sprint(t.UnixNano() / 1000 / 1000)
}

func parseObjectID(s string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(s) {
		return "", errors.New("invalid ObjectId")
	}
	return bson.ObjectIdHex(s), nil
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
