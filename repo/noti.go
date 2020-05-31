package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/abyss.club/uexky/lib/postgres"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
)

type NotiRepo struct{}

func (n *NotiRepo) db(ctx context.Context) postgres.Session {
	return postgres.GetSessionFromContext(ctx)
}

func (n *NotiRepo) toEntitySystemNoti(noti *Notification) *entity.SystemNoti {
	en := &entity.SystemNoti{
		ID:        noti.ID,
		Type:      entity.NotiTypeSystem,
		EventTime: noti.CreatedAt,
		SendTo: entity.SendTo{
			UserID:    noti.SendTo,
			SendGroup: (*entity.SendGroup)(noti.SendToGroup),
		},
	}
	_ = json.Unmarshal(noti.Content, &en.Content)
	return en
}

func (n *NotiRepo) toEntityRepliedNoti(noti *Notification) *entity.RepliedNoti {
	en := &entity.RepliedNoti{
		ID:        noti.ID,
		Type:      entity.NotiTypeReplied,
		EventTime: noti.CreatedAt,
		SendTo: entity.SendTo{
			UserID:    noti.SendTo,
			SendGroup: (*entity.SendGroup)(noti.SendToGroup),
		},
	}
	_ = json.Unmarshal(noti.Content, &en.Content)
	return en
}

func (n *NotiRepo) toEntityQuotedNoti(noti *Notification) *entity.QuotedNoti {
	en := &entity.QuotedNoti{
		ID:        noti.ID,
		Type:      entity.NotiTypeQuoted,
		EventTime: noti.CreatedAt,
		SendTo: entity.SendTo{
			UserID:    noti.SendTo,
			SendGroup: (*entity.SendGroup)(noti.SendToGroup),
		},
	}
	_ = json.Unmarshal(noti.Content, &en.Content)
	return en
}

func (n *NotiRepo) GetUserUnreadCount(ctx context.Context, user *entity.User) (*entity.UnreadNotiCount, error) {
	var err error
	nc := &entity.UnreadNotiCount{}
	if nc.System, err = n.db(ctx).Model((*Notification)(nil)).
		Where("type = ?", entity.NotiTypeSystem).
		Where("id > ?", user.LastReadNoti.SystemNoti).Count(); err != nil {
		return nil, err
	}
	if nc.Replied, err = n.db(ctx).Model((*Notification)(nil)).
		Where("type = ?", entity.NotiTypeReplied).
		Where("id > ?", user.LastReadNoti.RepliedNoti).Count(); err != nil {
		return nil, err
	}
	if nc.Quoted, err = n.db(ctx).Model((*Notification)(nil)).
		Where("type = ?", entity.NotiTypeQuoted).
		Where("id > ?", user.LastReadNoti.RepliedNoti).Count(); err != nil {
		return nil, err
	}
	return nc, nil
}

func (n *NotiRepo) GetNotiSlice(
	ctx context.Context, search *entity.NotiSearch, query entity.SliceQuery,
) (*entity.NotiSlice, error) {
	var notis []Notification
	q := n.db(ctx).Model(&notis).Where("type = ?", search.Type).
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.WhereOr("send_to = ?", search.UserID).
				WhereOr("send_to_group = ?", entity.AllUser)
			return q, nil
		})
	applySlice := func(q *orm.Query, isAfter bool, cursor string) (*orm.Query, error) {
		if cursor == "" {
			return q, nil
		}
		id, err := strconv.Atoi(cursor)
		if err != nil {
			return nil, err
		}
		if !isAfter {
			// before
			return q.Where("id > ?", id).Order("id"), nil
		}
		// after
		return q.Where("? < ?", id).Order("id DESC"), nil
	}
	var err error
	q, err = applySliceQuery(applySlice, q, &query)
	if err != nil {
		return nil, err
	}
	if err := q.Select(); err != nil {
		return nil, err
	}
	sliceInfo := &entity.SliceInfo{HasNext: len(notis) > query.Limit}
	slice := &entity.NotiSlice{}
	dealSlice := func(i int, isFirst bool, isLast bool) {
		switch search.Type {
		case entity.NotiTypeSystem:
			slice.System = append(slice.System, n.toEntitySystemNoti(&notis[i]))
		case entity.NotiTypeReplied:
			slice.Replied = append(slice.Replied, n.toEntityRepliedNoti(&notis[i]))
		case entity.NotiTypeQuoted:
			slice.Quoted = append(slice.Quoted, n.toEntityQuotedNoti(&notis[i]))
		}
		if isFirst {
			sliceInfo.FirstCursor = fmt.Sprintf("%v", notis[i].ID)
		}
		if isLast {
			sliceInfo.LastCursor = fmt.Sprintf("%v", notis[i].ID)
		}
	}
	dealSliceResult(dealSlice, &query, len(notis), query.Before != nil)
	return slice, nil
}

func (n *NotiRepo) InsertNoti(ctx context.Context, insert entity.NotiInsert) error {
	noti := &Notification{}
	var err error
	switch {
	case insert.System != nil:
		key := fmt.Sprintf("system:%s", uid.NewUID().ToBase64String())
		noti.Key = &key
		noti.Type = string(entity.NotiTypeSystem)
		noti.SendTo = insert.System.SendTo.UserID
		noti.SendToGroup = (*string)(insert.System.SendTo.SendGroup)
		if noti.Content, err = json.Marshal(insert.System); err != nil {
			return err
		}
	case insert.Replied != nil:
		key := fmt.Sprintf("replied:%s", insert.Replied.Content.ThreadID)
		noti.Key = &key
		noti.Type = string(entity.NotiTypeReplied)
		noti.SendTo = insert.Replied.SendTo.UserID
		noti.SendToGroup = (*string)(insert.Replied.SendTo.SendGroup)
		if noti.Content, err = json.Marshal(insert.Replied); err != nil {
			return err
		}
	case insert.Quoted != nil:
		key := fmt.Sprintf("replied:%s:%s", insert.Quoted.Content.QuotedID, insert.Quoted.Content.PostID)
		noti.Key = &key
		noti.Type = string(entity.NotiTypeQuoted)
		noti.SendTo = insert.Quoted.SendTo.UserID
		noti.SendToGroup = (*string)(insert.Quoted.SendTo.SendGroup)
		if noti.Content, err = json.Marshal(insert.Quoted); err != nil {
			return err
		}
	}
	_, err = n.db(ctx).Model(noti).OnConflict("(key) DO UPDATE").Set("updated_at=now()").Returning("*").Insert()
	return err
}

func (n *NotiRepo) UpdateReadID(ctx context.Context, userID int, nType entity.NotiType, id int) error {
	user := User{}
	column := fmt.Sprintf("last_read_%s_noti", nType)
	_, err := n.db(ctx).Model(user).Set("? = ?", column, id).Where("id = ?", userID).Update()
	return err
}
