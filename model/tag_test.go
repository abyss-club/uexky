package model

import (
	"log"
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/mw"
)

func Test_isMainTag(t *testing.T) {
	tests := []struct {
		name string
		tag  string
		want bool
	}{
		{"main", "MainA", true},
		{"not main", "Aha", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMainTag(tt.tag); got != tt.want {
				t.Errorf("isMainTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpsertTags(t *testing.T) {
	if err := UpsertTags(testCtx, "MainB", []string{"SubA", "SubB"}); err != nil {
		t.Fatalf("UpsertTags() error = %v", err)
	}
	c := mw.GetMongo(testCtx).C(colleTag)
	var tags []*Tag
	if err := c.Find(bson.M{"parent": "MainB"}).All(&tags); err != nil {
		t.Fatalf("UpsertTags(), find tag error = %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("UpsertTags() should insert two tags")
	}
	if !(tags[0].Name == "SubA" && tags[1].Name == "SubB") &&
		!(tags[0].Name == "SubB" && tags[1].Name == "SubA") {
		t.Fatalf("UpsertTags() should insert SubA and SubB, get %+v, %+v", *tags[0], *tags[1])
	}
}

func TestGetTagTree(t *testing.T) {
	log.Println("start TestGetTagTree()")
	ctx := ctxWithUser(mockUsers[2])
	subTagsList := [][]string{
		[]string{"Sub0", "Sub1"},
		[]string{"Sub2"},
		[]string{"Sub0", "Sub3"},
		[]string{"Sub2", "Sub4"},
		[]string{"Sub0", "Sub1", "Sub5"},
		[]string{"Sub6"},
		[]string{"Sub7"},
		[]string{"Sub8"},
		[]string{"Sub9"},
	}
	wantTags := []string{
		"Sub9", "Sub8", "Sub7", "Sub6", "Sub5",
		"Sub1", "Sub0", "Sub4", "Sub2", "Sub3",
	}
	log.Println("insert threads")
	for _, subTags := range subTagsList {
		if _, err := NewThread(ctx, &ThreadInput{
			Anonymous: true,
			Content:   "content",
			MainTag:   mgmt.Config.MainTags[2],
			SubTags:   &subTags,
		}); err != nil {
			t.Fatal(errors.Wrap(err, "create thread"))
		}
	}
	log.Println("GetTagTree")
	tree, err := GetTagTree(ctx)
	if err != nil {
		t.Fatal(errors.Wrap(err, "GetTagTree()"))
	}
	if len(tree.Nodes) != len(mgmt.Config.MainTags) {
		t.Fatalf("GetTagTree() should have %v node, got %v",
			len(mgmt.Config.MainTags), len(tree.Nodes))
	}
	if !reflect.DeepEqual(tree.Nodes[2].SubTags, wantTags) {
		t.Fatalf("GetTagTree().Nodes[2] = %q, want: %q",
			tree.Nodes[2].SubTags, wantTags)
	}
}
