package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type NotiRepo struct{}

func (n *NotiRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (n *NotiRepo) ToEntityNoti(noti *Notification) *entity.Notification {
	e := &entity.Notification{
		Type:      entity.NotiType(noti.Type),
		EventTime: noti.CreatedAt,
		Key:       noti.Key,
		SortKey:   uid.UID(noti.SortKey),
		Receivers: noti.Receivers,
	}
	var err error
	switch e.Type {
	case entity.NotiTypeSystem:
		var content entity.SystemNoti
		err = json.Unmarshal(noti.Content, &content)
		e.Content = content
	case entity.NotiTypeReplied:
		var content entity.RepliedNoti
		err = json.Unmarshal(noti.Content, &content)
		e.Content = content
	case entity.NotiTypeQuoted:
		var content entity.QuotedNoti
		err = json.Unmarshal(noti.Content, &content)
		e.Content = content
	default:
		err = fmt.Errorf("can't marshal noti content, invalid type '%s'", e.Type)
	}
	if err != nil {
		panic(uerr.Errorf(uerr.InternalError, "read notification error: %w", err))
	}
	return e
}

func (n *NotiRepo) GetUserUnreadCount(ctx context.Context, user *entity.User) (int, error) {
	var count int
	_, err := n.db(ctx).Query(orm.Scan(&count),
		`SELECT count(n.*) FROM notification as n
		LEFT JOIN public."user" as u ON u.id = ?
		WHERE sort_key > u.last_read_noti AND updated_at > u.created_at
		AND n.receivers && ?`,
		user.ID, pg.Array(user.NotiReceivers()),
	)
	return count, err
}

func (n *NotiRepo) GetNotiSlice(
	ctx context.Context, search *entity.NotiSearch, query entity.SliceQuery,
) (*entity.NotiSlice, error) {
	var notifications []Notification
	receivers := (&entity.User{ID: search.UserID}).NotiReceivers()
	q := n.db(ctx).Model(&notifications).Where("receivers && ?", receivers)
	applySlice := func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error) {
		if cursor != "" {
			c, err := uid.ParseUID(cursor)
			if err != nil {
				return nil, err
			}
			if !isAfter {
				q = q.Where("sort_key > ?", c)
			} else {
				q = q.Where("sort_key < ?", c)
			}
		}
		if !isAfter {
			return q.Order("sort_key"), nil
		}
		return q.Order("sort_key DESC"), nil
	}

	var err error
	q, err = applySliceQuery(applySlice, q, &query)
	if err != nil {
		return nil, err
	}
	if err := q.Select(); err != nil {
		return nil, err
	}

	sliceInfo := &entity.SliceInfo{HasNext: len(notifications) > query.Limit}
	var entities []*entity.Notification
	dealSlice := func(i int, isFirst bool, isLast bool) {
		entities = append(entities, n.ToEntityNoti(&notifications[i]))
		if isFirst {
			sliceInfo.FirstCursor = uid.UID(notifications[i].SortKey).ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = uid.UID(notifications[i].SortKey).ToBase64String()
		}
	}
	dealSliceResult(dealSlice, &query, len(notifications), query.Before != nil)
	return &entity.NotiSlice{
		Notifications: entities,
		SliceInfo:     sliceInfo,
	}, nil
}

func (n *NotiRepo) InsertNoti(ctx context.Context, noti *entity.Notification) error {
	notification := &Notification{
		Key:       noti.Key,
		SortKey:   int64(noti.SortKey),
		Type:      noti.Type.String(),
		Receivers: noti.Receivers,
	}
	var err error
	switch noti.Type {
	case entity.NotiTypeSystem:
		notification.Content, err = json.Marshal(noti.Content.(entity.SystemNoti))
		if err != nil {
			return err
		}
	case entity.NotiTypeReplied:
		notification.Content, err = json.Marshal(noti.Content.(entity.RepliedNoti))
		if err != nil {
			return err
		}
	case entity.NotiTypeQuoted:
		notification.Content, err = json.Marshal(noti.Content.(entity.QuotedNoti))
		if err != nil {
			return err
		}
	default:
		return uerr.Errorf(uerr.ParamsError, "invalid noti type '%s'", noti.Type)
	}
	return n.db(ctx).Insert(&notification)
}

func (n *NotiRepo) UpdateReadID(ctx context.Context, userID int, id int) error {
	user := &User{}
	_, err := n.db(ctx).Model(user).
		Set("last_read_noti = ?", id).Where("id = ?", userID).Update()
	return err
}
