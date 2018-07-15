package model

import (
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
	"gitlab.com/abyss.club/uexky/mw"
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
			"normal desk",
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

func TestSliceQuery_Find(t *testing.T) {
	type sample struct {
		OID  bson.ObjectId `bson:"_id"`
		Text string        `bson:"text"`
	}
	colle := "slice_sample"
	var samples []*sample
	for i := 0; i < 9; i++ {
		c := mw.GetMongo(testCtx).C(colle)
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
			results := []*sample{}
			err := sq.Find(testCtx, colle, bson.M{"text": "aha"}, &results)
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
