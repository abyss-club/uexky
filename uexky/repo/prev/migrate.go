package repo

import (
	"strconv"

	"github.com/go-pg/pg/v9"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/lib/algo"
	"gitlab.com/abyss.club/uexky/lib/uid"
	"gitlab.com/abyss.club/uexky/uexky/entity"
	"gitlab.com/abyss.club/uexky/uexky/repo"
)

type Migrator struct {
	PrevDB *pg.DB
	NewDB  *pg.DB
}

func (m *Migrator) DoMigrate() error {
	var users []User
	if err := m.PrevDB.Model(&users).Select(); err != nil {
		return errors.Wrap(err, "find all users")
	}
	var threads []Thread
	if err := m.PrevDB.Model(&threads).Select(); err != nil {
		return errors.Wrap(err, "find all users")
	}
	var posts []Post
	if err := m.PrevDB.Model(&posts).Select(); err != nil {
		return errors.Wrap(err, "find all users")
	}
	var notis []Notification
	if err := m.PrevDB.Model(&notis).Select(); err != nil {
		return errors.Wrap(err, "find all users")
	}
	var pMainTags []Tag
	if err := m.PrevDB.Model(&pMainTags).Where("is_main = ?", true).Select(); err != nil {
		return errors.Wrap(err, "find all users")
	}

	userIDMap := map[int]uid.UID{}
	for i := range users {
		userIDMap[users[i].ID] = uid.NewUIDFromTime(users[i].CreatedAt)
	}
	notiIDMap := map[int]uid.UID{}
	for i := range notis {
		notiIDMap[notis[i].ID] = uid.NewUIDFromTime(notis[i].UpdatedAt)
	}

	mainTags, err := m.migrateMainTags(pMainTags)
	if err != nil {
		return errors.Wrap(err, "migrateMainTags")
	}
	for i := range users {
		if err := m.migrateUser(&users[i], userIDMap, notiIDMap); err != nil {
			return err
		}
	}
	for i := range threads {
		if err := m.migrateThread(&threads[i], userIDMap, mainTags); err != nil {
			return err
		}
	}
	for i := range posts {
		if err := m.migratePost(&posts[i], userIDMap); err != nil {
			return err
		}
	}
	for i := range notis {
		if err := m.migrateNotification(&notis[i], userIDMap, notiIDMap); err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) migrateMainTags(pTags []Tag) ([]string, error) {
	var tags []repo.Tag
	tagType := "main"
	for _, pt := range pTags {
		tags = append(tags, repo.Tag{
			Name:    pt.Name,
			TagType: &tagType,
		})
	}
	if _, err := m.NewDB.Model(&tags).Insert(); err != nil {
		return nil, errors.Wrap(err, "insert main tags")
	}
	var rst []string
	for _, tag := range tags {
		rst = append(rst, tag.Name)
	}
	return rst, nil
}

func maxUID(uids ...uid.UID) uid.UID {
	var max uid.UID
	for _, u := range uids {
		if u > max {
			max = u
		}
	}
	return max
}

func (m *Migrator) migrateUser(pUser *User, userIDMap, notiIDMap map[int]uid.UID) error {
	nUser := &repo.User{
		ID:           userIDMap[pUser.ID],
		CreatedAt:    pUser.CreatedAt,
		UpdatedAt:    pUser.UpdatedAt,
		Email:        &pUser.Email,
		Name:         pUser.Name,
		Role:         pUser.Role,
		LastReadNoti: 0,
		Tags:         nil,
	}
	// LastReadNoti
	nUser.LastReadNoti = maxUID(
		notiIDMap[pUser.LastReadSystemNoti],
		notiIDMap[pUser.LastReadRepliedNoti],
		notiIDMap[pUser.LastReadQuotedNoti],
	)
	if nUser.ID == 0 || nUser.LastReadNoti == 0 {
		return errors.New("user's id and last_read_noti can't be zero")
	}
	// Tags
	var tags []UsersTag
	if err := m.PrevDB.Model(&tags).Where("user_id = ?", pUser.ID).Select(); err != nil {
		return errors.Wrap(err, "find user's tag")

	}
	for _, tag := range tags {
		nUser.Tags = append(nUser.Tags, tag.TagName)
	}

	if _, err := m.NewDB.Model(nUser).Returning("*").Insert(); err != nil {
		return errors.Wrapf(err, "insert user %v", pUser.ID)
	}
	return nil
}

func (m *Migrator) migrateThread(pThread *Thread, userIDMap map[int]uid.UID, mainTags []string) error {
	nThread := &repo.Thread{
		ID:         pThread.ID,
		LastPostID: pThread.LastPostID,
		CreatedAt:  pThread.CreatedAt,
		UpdatedAt:  pThread.UpdatedAt,
		UserID:     userIDMap[pThread.User.ID],
		Anonymous:  pThread.Anonymous,
		Guest:      false,
		Author:     "",
		Title:      pThread.Title,
		Content:    pThread.Content,
		Locked:     pThread.Locked,
		Blocked:    pThread.Blocked,
		Tags:       nil,
	}

	// Author
	if pThread.Anonymous {
		nThread.Author = pThread.AnonymousID.ToBase64String()
	} else {
		nThread.Author = *pThread.UserName
	}

	// Tags
	var tags []ThreadsTag
	if err := m.PrevDB.Model(&tags).Where("thread_id = ?", pThread.ID).Select(); err != nil {
		return errors.Wrap(err, "find thread's tags")
	}
	var mainTag string
	var subTags []string
	for _, t := range tags {
		if algo.InStrSlice(mainTags, t.TagName) {
			mainTag = t.TagName
		} else {
			subTags = append(subTags, t.TagName)
		}
	}
	nThread.Tags = append([]string{mainTag}, subTags...)

	if _, err := m.NewDB.Model(nThread).Returning("*").Insert(); err != nil {
		return errors.Wrapf(err, "insert thread %v", pThread.ID)
	}
	return nil
}

func (m *Migrator) migratePost(pPost *Post, userIDMap map[int]uid.UID) error {
	nPost := &repo.Post{
		ID:        pPost.ID,
		CreatedAt: pPost.CreatedAt,
		UpdatedAt: pPost.UpdatedAt,
		ThreadID:  pPost.ThreadID,
		UserID:    userIDMap[pPost.UserID],
		Anonymous: pPost.Anonymous,
		Guest:     false,
		Author:    "",
		Blocked:   false,
		Content:   pPost.Content,
		QuotedIDs: nil,
	}

	// Author
	if pPost.Anonymous {
		nPost.Author = pPost.AnonymousID.ToBase64String()
	} else {
		nPost.Author = *pPost.UserName
	}

	// Blocked
	if pPost.Blocked != nil {
		nPost.Blocked = *pPost.Blocked
	}

	// QuotedIDs
	var quoteds []PostsQuote
	q := m.PrevDB.Model(&quoteds).Where("qouter_id = ?", pPost.ID).Order("quoted_id")
	if err := q.Select(); err != nil {
		return errors.Wrap(err, "find quoted posts")
	}
	for _, quoted := range quoteds {
		nPost.QuotedIDs = append(nPost.QuotedIDs, quoted.Quoted.ID)
	}

	if _, err := m.NewDB.Model(nPost).Returning("*").Insert(); err != nil {
		return errors.Wrap(err, "insert post")
	}
	return nil
}

func (m *Migrator) migrateNotification(pNoti *Notification, userIDMap, notiIDMap map[int]uid.UID) error {
	nNoti := repo.Notification{
		Key:       *pNoti.Key,
		SortKey:   notiIDMap[pNoti.ID],
		CreatedAt: pNoti.CreatedAt,
		UpdatedAt: pNoti.UpdatedAt,
		Type:      pNoti.Type,
		Receivers: nil,
		Content: map[string]interface{}{
			"": nil,
		},
	}
	switch {
	case pNoti.SendTo != nil:
		userID := userIDMap[*pNoti.SendTo]
		nNoti.Receivers = []entity.Receiver{entity.SendToUser(userID)}
	case pNoti.SendToGroup != nil:
		nNoti.Receivers = []entity.Receiver{entity.SendToGroup(entity.AllUser)}
	default:
		return errors.Errorf("no receivers in notification")
	}

	var err error
	switch pNoti.Type {
	case entity.NotiTypeSystem:
		err = fillSystemNotiContent(pNoti, &nNoti)
	case entity.NotiTypeReplied:
		err = m.fillRepliedContent(pNoti, &nNoti)
	case entity.NotiTypeQuoted:
		err = m.fillQuotedContent(pNoti, &nNoti)
	default:
		return errors.Errorf("invalide noti type")
	}
	if err != nil {
		return err
	}
	if _, err := m.NewDB.Model(&nNoti).Insert(); err != nil {
		return errors.Wrap(err, "insert noti")
	}
	return nil
}

func fillSystemNotiContent(pNoti *Notification, nNoti *repo.Notification) error {
	type pContent struct {
		Title   string `mapstructure:"title"`
		Content string `mapstructure:"content"`
	}
	var pc pContent
	if err := mapstructure.Decode(pNoti.Content, &pc); err != nil {
		return errors.Wrapf(err, "decode noti content of %v", pNoti.ID)
	}
	content := entity.SystemNoti{
		Title:   pc.Title,
		Content: pc.Content,
	}
	if err := mapstructure.Decode(&content, &nNoti.Content); err != nil {
		return errors.Wrapf(err, "encode noti content of %v", pNoti.ID)
	}
	return nil
}

func (m *Migrator) fillRepliedContent(pNoti *Notification, nNoti *repo.Notification) error {
	type pContent struct {
		ThreadID string `mapstructure:"threadId"`
	}
	var pc pContent
	if err := mapstructure.Decode(pNoti.Content, &pc); err != nil {
		return errors.Wrapf(err, "decode noti content of %v", pNoti.ID)
	}

	threadID, err := strconv.ParseInt(pc.ThreadID, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalivd thread id: %s", pc.ThreadID)
	}
	thread := repo.Thread{}
	q := m.NewDB.Model(&thread).Where("id = ?", threadID)
	if err := q.Select(); err != nil {
		return errors.Wrapf(err, "GetByID(id=%+v)", threadID)
	}
	content := entity.RepliedNoti{
		Thread: &entity.ThreadOutline{
			ID:      thread.ID,
			Title:   thread.Title,
			Content: thread.Content,
			MainTag: thread.Tags[0],
			SubTags: thread.Tags[1:],
		},
		NewRepliesCount: 0,
		FirstReplyID:    thread.LastPostID,
	}
	if err := mapstructure.Decode(&content, &nNoti.Content); err != nil {
		return errors.Wrapf(err, "encode noti content of %v", pNoti.ID)
	}
	return nil
}

func (m *Migrator) fillQuotedContent(pNoti *Notification, nNoti *repo.Notification) error {
	type pContent struct {
		ThreadID string `mapstructure:"threadId"`
		QuotedID string `mapstructure:"quotedId"`
		PostID   string `mapstructure:"threadId"`
	}
	var pc pContent
	if err := mapstructure.Decode(pNoti.Content, &pc); err != nil {
		return errors.Wrapf(err, "decode noti content of %v", pNoti.ID)
	}
	threadID, err := strconv.ParseInt(pc.ThreadID, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalivd thread id: %s", pc.ThreadID)
	}
	quotedID, err := strconv.ParseInt(pc.QuotedID, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalivd quoted post id: %s", pc.QuotedID)
	}
	postID, err := strconv.ParseInt(pc.PostID, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "invalivd post id: %s", pc.PostID)
	}

	var quoted repo.Post
	if err := m.NewDB.Model(&quoted).Where("id = ?", quotedID).Select(); err != nil {
		return errors.Wrapf(err, "GetPost(id=%v)", quotedID)
	}
	var post repo.Post
	if err := m.NewDB.Model(&post).Where("id = ?", postID).Select(); err != nil {
		return errors.Wrapf(err, "GetPost(id=%v)", postID)
	}

	content := entity.QuotedNoti{
		ThreadID: uid.UID(threadID),
		QuotedPost: &entity.PostOutline{
			ID: quoted.ID,
			Author: &entity.Author{
				UserID:    quoted.UserID,
				Guest:     quoted.Guest,
				Anonymous: quoted.Anonymous,
				Author:    quoted.Author,
			},
			Content: quoted.Content,
		},
		Post: &entity.PostOutline{
			ID: post.ID,
			Author: &entity.Author{
				UserID:    post.UserID,
				Guest:     post.Guest,
				Anonymous: post.Anonymous,
				Author:    post.Author,
			},
			Content: post.Content,
		},
	}
	if err := mapstructure.Decode(&content, &nNoti.Content); err != nil {
		return errors.Wrapf(err, "encode noti content of %v", pNoti.ID)
	}
	return nil
}
