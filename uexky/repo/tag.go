package repo

import (
	"context"
	"fmt"

	"gitlab.com/abyss.club/uexky/lib/config"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type TagRepo struct {
}

func (r *TagRepo) GetMainTags(ctx context.Context) ([]string, error) {
	var tags []Tag
	if err := db(ctx).Model(&tags).Where("type= ?", "main").Select(); err != nil {
		return nil, dbErrWrap(err, "GetMainTags")
	}
	var mainTags []string
	for _, t := range tags {
		mainTags = append(mainTags, t.Name)
	}
	return mainTags, nil
}

func (r *TagRepo) SetMainTags(ctx context.Context, mainTags []string) error {
	var tags []Tag
	tagType := "main"
	for _, t := range mainTags {
		tags = append(tags, Tag{
			Name:    t,
			TagType: &tagType,
		})
	}
	if _, err := db(ctx).Model(&mainTags).Insert(); err != nil {
		return dbErrWrapf(err, "SetMainTags(tags=%v)", tags)
	}
	return nil
}

func (r *TagRepo) Search(ctx context.Context, search *entity.TagSearch) ([]*entity.Tag, error) {
	type tag struct {
		Tag string `pg:"tag"`
	}
	var tags []tag
	var where, limit string
	if search.Text != "" {
		where = fmt.Sprintf("WHERE tag LIKE '%%%s%%'", search.Text)
	}
	if search.Limit != 0 {
		limit = fmt.Sprintf("LIMIT %v", search.Limit)
	}
	sql := fmt.Sprintf(`SELECT tag FROM (
		SELECT unnest(tags) as tag, max(created_at) as updated_at
		FROM thread group by tag
	) as tags %s ORDER BY updated_at DESC %s`, where, limit)
	if _, err := db(ctx).Query(&tags, sql); err != nil {
		return nil, dbErrWrapf(err, "GetTags(search=%+v)", search)
	}
	var entities []*entity.Tag
	for _, t := range tags {
		entities = append(entities, &entity.Tag{
			Name:   t.Tag,
			IsMain: config.IsMainTag(t.Tag),
		})
	}
	return entities, nil
}
