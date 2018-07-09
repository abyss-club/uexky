package model

import (
	"reflect"
	"testing"

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
