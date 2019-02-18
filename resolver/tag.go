package resolver

import (
	"context"

	"gitlab.com/abyss.club/uexky-go/config"
	"gitlab.com/abyss.club/uexky-go/model"
	"gitlab.com/abyss.club/uexky-go/uexky"
)

// queries:

// Tags ...
func (r *Resolver) Tags(ctx context.Context) (*TagResolver, error) {
	return &TagResolver{}, nil
}

// types:

// TagResolver ...
type TagResolver struct{}

// MainTags ...
func (tr *TagResolver) MainTags(ctx context.Context) ([]string, error) {
	return config.Config.MainTags, nil
}

// Recommended ...
func (tr *TagResolver) Recommended(ctx context.Context) ([]string, error) {
	return config.Config.MainTags, nil // TODO
}

// Tree ...
func (tr *TagResolver) Tree(
	ctx context.Context,
	args struct{ Query *string },
) (*[]*TagTreeNodeResolver, error) {
	u := uexky.Pop(ctx)
	query := ""
	if args.Query != nil {
		query = *args.Query
	}
	tree, err := model.GetTagTree(u, query)
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
