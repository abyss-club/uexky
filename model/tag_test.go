package model

import (
	"log"
	"reflect"
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/config"
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
	if err := UpsertTags(mu[0], "MainB", []string{"SubA", "SubB"}); err != nil {
		t.Fatalf("UpsertTags() error = %v", err)
	}
	c := mu[0].Mongo.C(colleTag)
	var tags []*Tag
	if err := c.Find(bson.M{"main_tags": "MainB"}).All(&tags); err != nil {
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
	t.Log("insert threads")
	for _, subTags := range subTagsList {
		if _, err := NewThread(mu[2], &ThreadInput{
			Anonymous: true,
			Content:   "content",
			MainTag:   config.Config.MainTags[2],
			SubTags:   &subTags,
		}); err != nil {
			t.Fatal(errors.Wrap(err, "create thread"))
		}
	}
	wantTags := []string{
		"Sub9", "Sub8", "Sub7", "Sub6", "Sub5",
		"Sub1", "Sub0", "Sub4", "Sub2", "Sub3",
	}
	tests := []struct {
		name     string
		query    string
		wantTags []string
	}{
		{"no query", "", wantTags},
		{"query 'ub'", "ub", wantTags},
		{"query '2'", "2", []string{"Sub2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := GetTagTree(mu[2], tt.query)
			if err != nil {
				t.Fatal(errors.Wrap(err, "GetTagTree()"))
			}
			if len(tree.Nodes) != len(config.Config.MainTags) {
				t.Fatalf("GetTagTree() should have %v node, got %v",
					len(config.Config.MainTags), len(tree.Nodes))
			}
			if !reflect.DeepEqual(tree.Nodes[2].SubTags, tt.wantTags) {
				t.Fatalf("GetTagTree().Nodes[2] = %q, want: %q",
					tree.Nodes[2].SubTags, wantTags)
			}
		})
	}
}
