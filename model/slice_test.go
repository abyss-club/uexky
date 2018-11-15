package model

import (
	"reflect"
	"testing"
	"time"

	"github.com/globalsign/mgo/bson"
)

func TestSliceQuery_GenQueryByObjectID(t *testing.T) {
	cursors := []bson.ObjectId{
		bson.NewObjectId(),
		bson.NewObjectId(),
	}
	type fields struct {
		Limit  int
		Desc   bool
		Cursor string
	}
	tests := []struct {
		name    string
		fields  fields
		want    bson.M
		wantErr bool
	}{
		{
			"normal",
			fields{10, false, cursors[0].Hex()},
			bson.M{"_id": bson.M{"$gt": cursors[0]}},
			false,
		},
		{
			"normal desc",
			fields{10, true, cursors[1].Hex()},
			bson.M{"_id": bson.M{"$lt": cursors[1]}},
			false,
		},
		{"invalid", fields{10, true, "invalid"}, nil, true},
		{"empty", fields{10, true, ""}, bson.M{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := &SliceQuery{
				Limit:  tt.fields.Limit,
				Desc:   tt.fields.Desc,
				Cursor: tt.fields.Cursor,
			}
			got, err := sq.GenQueryByObjectID()
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceQuery.GenQueryByObjectID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceQuery.GenQueryByObjectID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceQuery_GenQueryByTime(t *testing.T) {
	cursor0, _ := parseTimeCursor(genTimeCursor(time.Now()))
	cursor1, _ := parseTimeCursor(genTimeCursor(time.Now().Add(time.Hour)))
	type fields struct {
		Limit  int
		Desc   bool
		Cursor string
	}
	type args struct {
		field string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bson.M
		wantErr bool
	}{
		{
			"normal",
			fields{10, false, genTimeCursor(cursor0)},
			args{"time"},
			bson.M{"time": bson.M{"$gt": cursor0}},
			false,
		},
		{
			"normal desc",
			fields{10, true, genTimeCursor(cursor1)},
			args{"time"},
			bson.M{"time": bson.M{"$lt": cursor1}},
			false,
		},
	}
	getTime := func(b bson.M, desc bool) time.Time {
		dotTime := b["time"].(bson.M)
		key := "$gt"
		if desc {
			key = "$lt"
		}
		return dotTime[key].(time.Time)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := &SliceQuery{
				Limit:  tt.fields.Limit,
				Desc:   tt.fields.Desc,
				Cursor: tt.fields.Cursor,
			}
			got, err := sq.GenQueryByTime(tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceQuery.GenQueryByTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !getTime(got, sq.Desc).Equal(getTime(tt.want, sq.Desc)) {
				t.Errorf("SliceQuery.GenQueryByTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceQuery_Find(t *testing.T) {
	type sample struct {
		OID  bson.ObjectId `bson:"_id"`
		Text string        `bson:"text"`
	}
	colle := "slice_sample"
	var samples []*sample
	for i := 0; i < 9; i++ {
		c := mu[0].Mongo.C(colle)
		s := &sample{bson.NewObjectId(), "aha"}
		if err := c.Insert(s); err != nil {
			t.Fatal("insert slice sample", err)
		}
		samples = append(samples, s)
	}

	type fields struct {
		Limit  int
		Desc   bool
		Cursor string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*sample
		wantErr bool
	}{
		{"asc", fields{3, false, samples[4].OID.Hex()}, samples[5:8], false},
		{"desc", fields{3, true, samples[4].OID.Hex()}, []*sample{
			samples[3], samples[2], samples[1],
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := &SliceQuery{
				Limit:  tt.fields.Limit,
				Desc:   tt.fields.Desc,
				Cursor: tt.fields.Cursor,
			}
			queryObj, err := sq.GenQueryByObjectID()
			if err != nil {
				t.Errorf("sq.GenQueryByObjectID error = %v", err)
				return
			}
			queryObj["text"] = "aha"
			results := []*sample{}
			err = sq.Find(mu[0], colle, "_id", queryObj, &results)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceQuery.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(results) != tt.fields.Limit {
				t.Errorf("SliceQuery.Find() length = %v, want length %v", len(results), tt.fields.Limit)
			}
			for i := 0; i < tt.fields.Limit; i++ {
				if results[i].OID != tt.want[i].OID {
					t.Errorf("SliceQuery.Find() got = %+v, want %+v", *results[i], *tt.want[i])
				}
			}
		})
	}
}
