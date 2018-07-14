package resolver

import (
	"context"

	"gitlab.com/abyss.club/uexky/mgmt"
	"gitlab.com/abyss.club/uexky/model"
)

// TagResolver ...
type TagResolver struct{}

// MainTags ...
func (tr *TagResolver) MainTags(ctx context.Context) ([]string, error) {
	return mgmt.Config.MainTags, nil
}

// Recommand ...
func (tr *TagResolver) Recommand(ctx context.Context) (*[]string, error) {
	return &mgmt.Config.MainTags, nil // TODO
}

// Tree ...
func (tr *TagResolver) Tree(ctx context.Context) (*[]*TagTreeNodeResolver, error) {
	tree, err := model.GetTagTree(ctx)
	if err != nil {
		return nil, err
	}
	if len(tree.Nodes) == 0 {
		return nil, nil
	}
	ttnrs := []*TagTreeNodeResolver{}
	for _, node := range tree.Nodes {
		ttnrs = append(ttnrs, &TagTreeNodeResolver{node})
	}
	return &ttnrs, nil
}

// TagTreeNodeResolver ...
type TagTreeNodeResolver struct {
	node *model.TagTreeNode
}

// MainTag ...
func (ttnr *TagTreeNodeResolver) MainTag(ctx context.Context) (string, error) {
	return ttnr.node.MainTag, nil
}

// SubTags ...
func (ttnr *TagTreeNodeResolver) SubTags(ctx context.Context) (*[]string, error) {
	if len(ttnr.node.SubTags) == 0 {
		return nil, nil
	}
	return &ttnr.node.SubTags, nil
}
