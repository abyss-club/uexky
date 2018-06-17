package model

import "github.com/globalsign/mgo/bson"

// SliceInfo ...
type SliceInfo struct {
	FirstCursor string
	LastCursor  string
}

// SliceQuery ...
type SliceQuery struct {
	Limit  int
	Before string
	After  string
}

// QueryObject ...
func (sq *SliceQuery) QueryObject() bson.M {
	if sq.After == "" && sq.Before == "" {
		return bson.M{}
	}
	if sq.After != "" && sq.Before != "" {
		return bson.M{"$gt": sq.After, "$lt": sq.Before}
	}
	if sq.After != "" {
		return bson.M{"$gt": sq.After}
	}
	return bson.M{"$lt": sq.Before}
}
