package repo

import (
	"context"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/abyss.club/uexky/lib/postgres"
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
	return count, postgres.ErrHandlef(err, "GetUserUnreadCount(user=%+v)", user)
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
		return nil, postgres.ErrHandlef(err, "GetNotiByKey(userID=%v, key=%s)", userID, key)
	}
	return notification.ToEntity(), nil
}

func (r *NotiRepo) GetSlice(ctx context.Context, user *entity.User, query entity.SliceQuery) (*entity.NotiSlice, error) {
	var notifications []NotificationQuery
	var entities []*entity.Notification
	q := db(ctx).Model(&notifications).
		Column("notification.*").
		ColumnExpr("u.last_read_noti >= sort_key as has_read").
		Join(`LEFT JOIN public."user" as u on u.id = ?`, user.ID).
		Where("receivers && ?", pg.Array(user.NotiReceivers()))
	h := sliceHelper{
		Column:      "sort_key",
		Desc:        true,
		TransCursor: func(s string) (interface{}, error) { return uid.ParseUID(s) },
		SQ:          &query,
	}
	if err := h.Select(q); err != nil {
		return nil, postgres.ErrHandlef(err, "GetNotiSlice(user=%+v, query=%+v)", user, query)
	}
	h.DealResults(len(notifications), func(i int) {
		entities = append(entities, (&notifications[i]).ToEntity())
	})
	sliceInfo := &entity.SliceInfo{HasNext: len(notifications) > query.Limit}
	if len(entities) > 0 {
		sliceInfo.FirstCursor = entities[0].SortKey.ToBase64String()
		sliceInfo.LastCursor = entities[len(entities)-1].SortKey.ToBase64String()
	}
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
	return postgres.ErrHandlef(err, "InsertNoti(noti=%+v)", noti)
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
	return postgres.ErrHandlef(err, "UpdateNotiContent(noti=%+v)", noti)
}

func (r *NotiRepo) UpdateReadID(ctx context.Context, user *entity.User, id uid.UID) error {
	u := &User{}
	_, err := db(ctx).Model(u).
		Set("last_read_noti = ?", id).Where("id = ?", user.ID).Update()
	return postgres.ErrHandlef(err, "UpdateReadID(userID=%v, id=%v)", user.ID, id)
}
