package repo

import (
	"context"

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

func (n *NotiRepo) ToEntityNoti(user *entity.User, noti *Notification) *entity.Notification {
	e := &entity.Notification{
		Type:      entity.NotiType(noti.Type),
		EventTime: noti.CreatedAt,
		HasRead:   noti.SortKey <= int64(user.LastReadNoti),
		Key:       noti.Key,
		SortKey:   uid.UID(noti.SortKey),
		Receivers: noti.Receivers,
	}
	if err := e.DecodeContent(noti.Content); err != nil {
		panic(uerr.Errorf(uerr.InternalError, "read notification error: %w", err))
	}
	return e
}

func (n *NotiRepo) GetUserUnreadCount(ctx context.Context, user *entity.User) (int, error) {
	var count int
	_, err := n.db(ctx).Query(orm.Scan(&count),
		`SELECT count(n.*) FROM notification as n
		LEFT JOIN public."user" as u ON u.id = ?
		WHERE sort_key > u.last_read_noti AND n.updated_at > u.created_at
		AND n.receivers && ?`,
		user.ID, pg.Array(user.NotiReceivers()),
	)
	return count, err
}

func (n *NotiRepo) GetNotiByKey(ctx context.Context, user *entity.User, key string) (*entity.Notification, error) {
	var notification Notification
	err := n.db(ctx).Model(&notification).Where("key = ?", key).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return n.ToEntityNoti(user, &notification), err
}

func (n *NotiRepo) GetNotiSlice(
	ctx context.Context, user *entity.User, query entity.SliceQuery,
) (*entity.NotiSlice, error) {
	var notifications []Notification
	receivers := user.NotiReceivers()
	q := n.db(ctx).Model(&notifications).Where("receivers && ?", pg.Array(receivers))
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
		entities = append(entities, n.ToEntityNoti(user, &notifications[i]))
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
	content, err := noti.EncodeContent()
	if err != nil {
		return err
	}
	notification.Content = content
	_, err = n.db(ctx).Model(notification).Insert()
	return err
}

func (n *NotiRepo) UpdateNotiContent(ctx context.Context, noti *entity.Notification) error {
	content, err := noti.EncodeContent()
	if err != nil {
		return err
	}
	var notification Notification
	_, err = n.db(ctx).Model(&notification).
		Set("content = ?", content).
		Set("sort_key = ?", noti.SortKey).
		Where("key = ?", noti.Key).Update()
	return err
}

func (n *NotiRepo) UpdateReadID(ctx context.Context, userID int64, id uid.UID) error {
	user := &User{}
	_, err := n.db(ctx).Model(user).
		Set("last_read_noti = ?", id).Where("id = ?", userID).Update()
	return err
}
