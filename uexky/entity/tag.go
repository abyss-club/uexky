package entity

import "context"

type TagSearch struct {
	Text  string
	Limit int
}

type TagRepo interface {
	SetMainTags(ctx context.Context, mainTags []string) error
	GetMainTags(ctx context.Context) ([]string, error)

	Search(ctx context.Context, search *TagSearch) ([]*Tag, error)
}
