package model

import (
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
)

func TestSliceQuery_QueryObject(t *testing.T) {
	type fields struct {
		Limit  int
		Before string
		After  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bson.M
	}{
		{"empty", fields{10, "", ""}, nil},
		{"before", fields{10, "B", ""}, bson.M{"$lt": "B"}},
		{"after", fields{10, "", "A"}, bson.M{"$gt": "A"}},
		{"before & after", fields{10, "B", "A"}, bson.M{"$lt": "B", "$gt": "A"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sq := &SliceQuery{
				Limit:  tt.fields.Limit,
				Before: tt.fields.Before,
				After:  tt.fields.After,
			}
			if got := sq.QueryObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceQuery.QueryObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
