import { NotFoundError } from '~/utils/error';
import { query, doTransaction } from '~/utils/pg';
import { ACTION } from '~/models/user';
import UID from '~/uid';
import UserAidModel from '~/models/userAid';
import NotificationModel from '~/models/notification';

// pgm.createTable('post', {
//   id: { type: 'bigint', primaryKey: true },
//   createdAt: { type: 'timestamp', notNull: true },
//   updatedAt: { type: 'timestamp', notNull: true },
//   threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },
//
//   anonymous: { type: 'bool', notNull: true },
//   userId: { type: 'integer', notNull: true, references: 'user(id)' },
//   userName: { type: 'varchar(16)', references: 'user(name)' },
//   anonymousId: { type: 'bigint', references: 'anonymous_id(anonymousId)' },
//
//   blocked: { type: 'bool', default: false },
//   content: { type: 'text', notNull: true },
// });
// pgm.createIndex('threadId', 'userId', 'anonymous');
// pgm.createTable('posts_quotes', {
//   id: { type: 'serial', primaryKey: true },
//   quoterId: { type: 'bigint', notNull: true, references: 'post(id)' },
//   quotedId: { type: 'bigint', notNull: true, references: 'post(id)' },
// });
//
// type Post {
//     id: String!
//     createdAt: Time!
//     anonymous: Boolean!
//     author: String!
//     content: String!
//     quotes: [Post!]
//     quoteCount: Int!
//     blocked: Boolean!
// }

const makePost = function makePost(raw) {
  return {
    id: UID.parse(raw.id),
    createdAt: raw.createdAt,
    updatedAt: raw.updatedAt,
    anonymous: raw.anonymous,
    author: raw.anonymous ? UID.parse(raw.anonymousId).duid : raw.userName,
    content: raw.blocked ? '[此内容已被屏蔽]' : raw.content,
    blocked: raw.blocked,

    async getQuotes() {
      const { rows } = await query(`SELECT *
        FROM post INNER JOIN posts_quotes
        ON post.id = posts_quotes."quoterId"
        WHERE posts_quotes."quoterId" = $1 ORDER BY post.id`,
      [this.id.suid]);
      return rows.map(row => makePost(row));
    },

    async getQuotedCount() {
      const { rows } = await query(`SELECT count(*) FROM posts_quotes
        WHERE "quotedId"=$1`, [this.id.suid]);
      return rows[0].count;
    },
  };
};

const PostModel = {

  async findById({ postId }) {
    const { rows } = await query('SELECT * FROM post WHERE id=$1', [UID.parse(postId).suid]);
    if ((rows || []).length === 0) {
      throw NotFoundError(`can't find post ${postId}`);
    }
    return makePost(rows[0]);
  },

  async findThreadPosts({ threadId, query: sq }) {
    const sql = [
      'SELECT * FROM post',
      'WHERE post."threadId" = $1',
      sq.before && 'AND post.id < $2',
      sq.after && 'AND post.id > $3',
      'ORDER BY post.id DESC',
      `LIMIT ${sq.limit || 0}`,
    ].join(' ');
    const { rows } = await query(sql, [UID.parse(threadId).suid, sq.before, sq.after]);
    return rows.map(row => makePost(row));
  },

  async findUserPosts({ user, query: sq }) {
    const sql = [
      'SELECT * FROM post',
      'WHERE post."userId" = $1',
      sq.before && 'AND post.id < $2',
      sq.after && 'AND post.id > $3',
      'ORDER BY post.id DESC',
      `LIMIT ${sq.limit || 0}`,
    ].join(' ');
    const { rows } = await query(sql, [user.id, sq.before, sq.after]);
    return rows.map(row => makePost(row));
  },

  // input PostInput {
  //     threadId: String!
  //     anonymous: Boolean!
  //     content: String!
  //     # Set quoting PostIDs.
  //     quoteIds: [String!]
  // }
  async new({ ctx, post: input }) {
    ctx.auth.ensurePermission(ACTION.PUB_POST);
    const user = ctx.auth.signedInUser();
    const threadId = UID.parse(input.threadId);
    const postId = UID.new();
    const raw = {
      id: postId.suid,
      threadId: threadId.suid,
      anonymous: input.anonymous,
      userId: user.id,
      content: input.content,
    };
    let newPost;
    await doTransaction(async (txn) => {
      if (input.anonymous) {
        raw.anonymousId = await UserAidModel.getAid({
          txn, userId: user.id, threadId: raw.id,
        });
      } else {
        raw.userName = user.name;
      }
      const { rows } = await txn.query(`INSERT INTO post
        (id, anonymous, "userId", "userName", "anonymousId", content)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
      [raw.id.suid, raw.anonymous, raw.userId, raw.userName,
        raw.anonymous || raw.anonymousId.suid, raw.content]);
      newPost = makePost(rows[0]);
      if (raw.quoteIds.length > 0) {
        await Promise.all(raw.quoteIds.map(qid => txn.query(`INSERT
           INTO posts_quotes ("quoterId", "quotedId")
           VALUES ($1, $2)`, [newPost.id.suid, UID.parse(qid).suid])));
        await NotificationModel.newQuotedNoti({
          txn, threadId, postId, quotedIds: raw.quoteIds,
        });
      }
      await NotificationModel.newRepliedNoti({ txn, threadId });
    });
    return newPost;
  },

  async block({ ctx, postId }) {
    ctx.auth.ensurePermission(ACTION.BLOCK_POST);
    await query('UPDATE post SET lock=$1 WHERE id=$2', [true, UID.parse(postId).suid]);
    return true;
  },
};

export default PostModel;
