//nolint
//lint:file-ignore U1000 ignore unused code, it's generated
package model

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

const condition = "?.? = ?"

// base filters
type applier func(query *orm.Query) (*orm.Query, error)

type search struct {
	appliers []applier
}

func (s *search) apply(query *orm.Query) {
	for _, applier := range s.appliers {
		query.Apply(applier)
	}
}

func (s *search) where(query *orm.Query, table, field string, value interface{}) {
	query.Where(condition, pg.F(table), pg.F(field), value)
}

func (s *search) WithApply(a applier) {
	if s.appliers == nil {
		s.appliers = []applier{}
	}
	s.appliers = append(s.appliers, a)
}

func (s *search) With(condition string, params ...interface{}) {
	s.WithApply(func(query *orm.Query) (*orm.Query, error) {
		return query.Where(condition, params...), nil
	})
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query *orm.Query) *orm.Query
	Q() applier

	With(condition string, params ...interface{})
	WithApply(a applier)
}

type AnonymouIdSearch struct {
	search

	ID          *int
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	ThreadID    *int64
	UserID      *int
	AnonymousID *int64
}

func (s *AnonymouIdSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.AnonymouId.Alias, Columns.AnonymouId.ID, s.ID)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.AnonymouId.Alias, Columns.AnonymouId.CreatedAt, s.CreatedAt)
	}
	if s.UpdatedAt != nil {
		s.where(query, Tables.AnonymouId.Alias, Columns.AnonymouId.UpdatedAt, s.UpdatedAt)
	}
	if s.ThreadID != nil {
		s.where(query, Tables.AnonymouId.Alias, Columns.AnonymouId.ThreadID, s.ThreadID)
	}
	if s.UserID != nil {
		s.where(query, Tables.AnonymouId.Alias, Columns.AnonymouId.UserID, s.UserID)
	}
	if s.AnonymousID != nil {
		s.where(query, Tables.AnonymouId.Alias, Columns.AnonymouId.AnonymousID, s.AnonymousID)
	}

	s.apply(query)

	return query
}

func (s *AnonymouIdSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type ConfigSearch struct {
	search

	ID *int
}

func (s *ConfigSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.Config.Alias, Columns.Config.ID, s.ID)
	}

	s.apply(query)

	return query
}

func (s *ConfigSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type CounterSearch struct {
	search

	Name  *string
	Count *int
}

func (s *CounterSearch) Apply(query *orm.Query) *orm.Query {
	if s.Name != nil {
		s.where(query, Tables.Counter.Alias, Columns.Counter.Name, s.Name)
	}
	if s.Count != nil {
		s.where(query, Tables.Counter.Alias, Columns.Counter.Count, s.Count)
	}

	s.apply(query)

	return query
}

func (s *CounterSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type NotificationSearch struct {
	search

	ID          *int
	Key         *string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	Type        *string
	SendTo      *int
	SendToGroup *string
}

func (s *NotificationSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.Notification.Alias, Columns.Notification.ID, s.ID)
	}
	if s.Key != nil {
		s.where(query, Tables.Notification.Alias, Columns.Notification.Key, s.Key)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.Notification.Alias, Columns.Notification.CreatedAt, s.CreatedAt)
	}
	if s.UpdatedAt != nil {
		s.where(query, Tables.Notification.Alias, Columns.Notification.UpdatedAt, s.UpdatedAt)
	}
	if s.Type != nil {
		s.where(query, Tables.Notification.Alias, Columns.Notification.Type, s.Type)
	}
	if s.SendTo != nil {
		s.where(query, Tables.Notification.Alias, Columns.Notification.SendTo, s.SendTo)
	}
	if s.SendToGroup != nil {
		s.where(query, Tables.Notification.Alias, Columns.Notification.SendToGroup, s.SendToGroup)
	}

	s.apply(query)

	return query
}

func (s *NotificationSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type PgmigrationSearch struct {
	search

	ID    *int
	Name  *string
	RunOn *time.Time
}

func (s *PgmigrationSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.Pgmigration.Alias, Columns.Pgmigration.ID, s.ID)
	}
	if s.Name != nil {
		s.where(query, Tables.Pgmigration.Alias, Columns.Pgmigration.Name, s.Name)
	}
	if s.RunOn != nil {
		s.where(query, Tables.Pgmigration.Alias, Columns.Pgmigration.RunOn, s.RunOn)
	}

	s.apply(query)

	return query
}

func (s *PgmigrationSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type PostSearch struct {
	search

	ID          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	ThreadID    *int64
	Anonymous   *bool
	UserID      *int
	UserName    *string
	AnonymousID *int64
	Blocked     *bool
	Content     *string
}

func (s *PostSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.ID, s.ID)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.CreatedAt, s.CreatedAt)
	}
	if s.UpdatedAt != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.UpdatedAt, s.UpdatedAt)
	}
	if s.ThreadID != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.ThreadID, s.ThreadID)
	}
	if s.Anonymous != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.Anonymous, s.Anonymous)
	}
	if s.UserID != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.UserID, s.UserID)
	}
	if s.UserName != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.UserName, s.UserName)
	}
	if s.AnonymousID != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.AnonymousID, s.AnonymousID)
	}
	if s.Blocked != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.Blocked, s.Blocked)
	}
	if s.Content != nil {
		s.where(query, Tables.Post.Alias, Columns.Post.Content, s.Content)
	}

	s.apply(query)

	return query
}

func (s *PostSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type PostsQuoteSearch struct {
	search

	ID       *int
	QuoterID *int64
	QuotedID *int64
}

func (s *PostsQuoteSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.PostsQuote.Alias, Columns.PostsQuote.ID, s.ID)
	}
	if s.QuoterID != nil {
		s.where(query, Tables.PostsQuote.Alias, Columns.PostsQuote.QuoterID, s.QuoterID)
	}
	if s.QuotedID != nil {
		s.where(query, Tables.PostsQuote.Alias, Columns.PostsQuote.QuotedID, s.QuotedID)
	}

	s.apply(query)

	return query
}

func (s *PostsQuoteSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type TagSearch struct {
	search

	Name      *string
	IsMain    *bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (s *TagSearch) Apply(query *orm.Query) *orm.Query {
	if s.Name != nil {
		s.where(query, Tables.Tag.Alias, Columns.Tag.Name, s.Name)
	}
	if s.IsMain != nil {
		s.where(query, Tables.Tag.Alias, Columns.Tag.IsMain, s.IsMain)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.Tag.Alias, Columns.Tag.CreatedAt, s.CreatedAt)
	}
	if s.UpdatedAt != nil {
		s.where(query, Tables.Tag.Alias, Columns.Tag.UpdatedAt, s.UpdatedAt)
	}

	s.apply(query)

	return query
}

func (s *TagSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type TagsMainTagSearch struct {
	search

	ID        *int
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Name      *string
	BelongsTo *string
}

func (s *TagsMainTagSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.TagsMainTag.Alias, Columns.TagsMainTag.ID, s.ID)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.TagsMainTag.Alias, Columns.TagsMainTag.CreatedAt, s.CreatedAt)
	}
	if s.UpdatedAt != nil {
		s.where(query, Tables.TagsMainTag.Alias, Columns.TagsMainTag.UpdatedAt, s.UpdatedAt)
	}
	if s.Name != nil {
		s.where(query, Tables.TagsMainTag.Alias, Columns.TagsMainTag.Name, s.Name)
	}
	if s.BelongsTo != nil {
		s.where(query, Tables.TagsMainTag.Alias, Columns.TagsMainTag.BelongsTo, s.BelongsTo)
	}

	s.apply(query)

	return query
}

func (s *TagsMainTagSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type ThreadSearch struct {
	search

	ID          *int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	Anonymous   *bool
	UserID      *int
	UserName    *string
	AnonymousID *int64
	Title       *string
	Content     *string
	Locked      *bool
	Blocked     *bool
	LastPostID  *int64
}

func (s *ThreadSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.ID, s.ID)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.CreatedAt, s.CreatedAt)
	}
	if s.UpdatedAt != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.UpdatedAt, s.UpdatedAt)
	}
	if s.Anonymous != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.Anonymous, s.Anonymous)
	}
	if s.UserID != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.UserID, s.UserID)
	}
	if s.UserName != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.UserName, s.UserName)
	}
	if s.AnonymousID != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.AnonymousID, s.AnonymousID)
	}
	if s.Title != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.Title, s.Title)
	}
	if s.Content != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.Content, s.Content)
	}
	if s.Locked != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.Locked, s.Locked)
	}
	if s.Blocked != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.Blocked, s.Blocked)
	}
	if s.LastPostID != nil {
		s.where(query, Tables.Thread.Alias, Columns.Thread.LastPostID, s.LastPostID)
	}

	s.apply(query)

	return query
}

func (s *ThreadSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type ThreadsTagSearch struct {
	search

	ID        *int
	CreatedAt *time.Time
	UpdatedAt *time.Time
	ThreadID  *int64
	TagName   *string
}

func (s *ThreadsTagSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.ThreadsTag.Alias, Columns.ThreadsTag.ID, s.ID)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.ThreadsTag.Alias, Columns.ThreadsTag.CreatedAt, s.CreatedAt)
	}
	if s.UpdatedAt != nil {
		s.where(query, Tables.ThreadsTag.Alias, Columns.ThreadsTag.UpdatedAt, s.UpdatedAt)
	}
	if s.ThreadID != nil {
		s.where(query, Tables.ThreadsTag.Alias, Columns.ThreadsTag.ThreadID, s.ThreadID)
	}
	if s.TagName != nil {
		s.where(query, Tables.ThreadsTag.Alias, Columns.ThreadsTag.TagName, s.TagName)
	}

	s.apply(query)

	return query
}

func (s *ThreadsTagSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type UserSearch struct {
	search

	ID                  *int
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
	Email               *string
	Name                *string
	Role                *string
	LastReadSystemNoti  *int
	LastReadRepliedNoti *int
	LastReadQuotedNoti  *int
}

func (s *UserSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.User.Alias, Columns.User.ID, s.ID)
	}
	if s.CreatedAt != nil {
		s.where(query, Tables.User.Alias, Columns.User.CreatedAt, s.CreatedAt)
	}
	if s.UpdatedAt != nil {
		s.where(query, Tables.User.Alias, Columns.User.UpdatedAt, s.UpdatedAt)
	}
	if s.Email != nil {
		s.where(query, Tables.User.Alias, Columns.User.Email, s.Email)
	}
	if s.Name != nil {
		s.where(query, Tables.User.Alias, Columns.User.Name, s.Name)
	}
	if s.Role != nil {
		s.where(query, Tables.User.Alias, Columns.User.Role, s.Role)
	}
	if s.LastReadSystemNoti != nil {
		s.where(query, Tables.User.Alias, Columns.User.LastReadSystemNoti, s.LastReadSystemNoti)
	}
	if s.LastReadRepliedNoti != nil {
		s.where(query, Tables.User.Alias, Columns.User.LastReadRepliedNoti, s.LastReadRepliedNoti)
	}
	if s.LastReadQuotedNoti != nil {
		s.where(query, Tables.User.Alias, Columns.User.LastReadQuotedNoti, s.LastReadQuotedNoti)
	}

	s.apply(query)

	return query
}

func (s *UserSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}

type UsersTagSearch struct {
	search

	ID      *int
	UserID  *int
	TagName *string
}

func (s *UsersTagSearch) Apply(query *orm.Query) *orm.Query {
	if s.ID != nil {
		s.where(query, Tables.UsersTag.Alias, Columns.UsersTag.ID, s.ID)
	}
	if s.UserID != nil {
		s.where(query, Tables.UsersTag.Alias, Columns.UsersTag.UserID, s.UserID)
	}
	if s.TagName != nil {
		s.where(query, Tables.UsersTag.Alias, Columns.UsersTag.TagName, s.TagName)
	}

	s.apply(query)

	return query
}

func (s *UsersTagSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		return s.Apply(query), nil
	}
}
