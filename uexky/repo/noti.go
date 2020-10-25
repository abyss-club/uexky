package repo

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/abyss.club/uexky/lib/uerr"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type NotiRepo struct{}

func (r *NotiRepo) GetUnreadCount(ctx context.Context, user *entity.User) (int, error) {
	var count int
	_, err := db(ctx).Query(orm.Scan(&count),
		`SELECT count(n.*) FROM notification as n
		LEFT JOIN public."user" as u ON u.id = ?
		WHERE sort_key > u.last_read_noti AND n.updated_at > u.created_at
		AND n.receivers && ?`,
		user.ID, pg.Array(user.NotiReceivers()),
	)
	return count, dbErrWrapf(err, "GetUserUnreadCount(user=%+v)", user)
}

func (r *NotiRepo) GetByKey(ctx context.Context, userID uid.UID, key string) (*entity.Notification, error) {
	var notification NotificationQuery
	err := db(ctx).Model(&notification).
		Column("notification.*").
		ColumnExpr("u.last_read_noti >= sort_key as has_read").
		Join(`LEFT JOIN public."user" as u ON u.id = ?`, userID).
		Where("key = ?", key).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, nil
		}
		return nil, dbErrWrapf(err, "GetNotiByKey(userID=%v, key=%s)", userID, key)
	}
	return notification.ToEntity(), nil
}

func (r *NotiRepo) GetSlice(ctx context.Context, user *entity.User, query entity.SliceQuery) (*entity.NotiSlice, error) {
	var notifications []NotificationQuery
	receivers := user.NotiReceivers()
	q := db(ctx).Model(&notifications).
		Column("notification.*").
		ColumnExpr("u.last_read_noti >= sort_key as has_read").
		Join(`LEFT JOIN public."user" as u on u.id = ?`, user.ID).
		Where("receivers && ?", pg.Array(receivers))
	applySlice := func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error) {
		if cursor != "" {
			c, err := uid.ParseUID(cursor)
			if err != nil {
				return nil, uerr.Wrap(uerr.ParamsError, err, "invalid cursor")
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
		return nil, dbErrWrapf(err, "GetNotiSlice(user=%+v, query=%+v)", user, query)
	}

	sliceInfo := &entity.SliceInfo{HasNext: len(notifications) > query.Limit}
	var entities []*entity.Notification
	dealSlice := func(i int, isFirst bool, isLast bool) {
		entities = append(entities, (&notifications[i]).ToEntity())
		if isFirst {
			sliceInfo.FirstCursor = notifications[i].SortKey.ToBase64String()
		}
		if isLast {
			sliceInfo.LastCursor = notifications[i].SortKey.ToBase64String()
		}
	}
	dealSliceResult(dealSlice, &query, len(notifications), query.Before != nil)
	return &entity.NotiSlice{
		Notifications: entities,
		SliceInfo:     sliceInfo,
	}, nil
}

func (r *NotiRepo) Insert(ctx context.Context, noti *entity.Notification) error {
	n, err := NewNotificaionFromEntity(noti)
	if err != nil {
		return err
	}
	_, err = db(ctx).Model(n).Insert()
	return dbErrWrapf(err, "InsertNoti(noti=%+v)", noti)
}

func (r *NotiRepo) UpdateContent(ctx context.Context, noti *entity.Notification) error {
	content, err := noti.EncodeContent()
	if err != nil {
		return err
	}
	var notification Notification
	_, err = db(ctx).Model(&notification).
		Set("content = ?", content).
		Set("sort_key = ?", noti.SortKey).
		Where("key = ?", noti.Key).Update()
	return dbErrWrapf(err, "UpdateNotiContent(noti=%+v)", noti)
}

func (r *NotiRepo) UpdateReadID(ctx context.Context, user *entity.User, id uid.UID) error {
	u := &User{}
	_, err := db(ctx).Model(u).
		Set("last_read_noti = ?", id).Where("id = ?", user.ID).Update()
	return dbErrWrapf(err, "UpdateReadID(userID=%v, id=%v)", user.ID, id)
}
