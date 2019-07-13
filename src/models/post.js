import { NotFoundError, ParamsError } from '~/utils/error';
import { query, doTransaction } from '~/utils/pg';
import { ACTION } from '~/models/user';
import UID from '~/uid';
import UserAidModel from '~/models/userAid';
import NotificationModel from '~/models/notification';
import querySlice from '~/models/slice';

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

const blockedContent = '[此内容已被管理员屏蔽]';

const makePost = function makePost(raw) {
  return {
    id: UID.parse(raw.id),
    createdAt: raw.createdAt,
    updatedAt: raw.updatedAt,
    anonymous: raw.anonymous,
    author: raw.anonymous ? UID.parse(raw.anonymousId).duid : raw.userName,
    content: raw.blocked ? blockedContent : raw.content,
    blocked: raw.blocked,

    async getQuotes() {
      const { rows } = await query(`SELECT post.*
        FROM post INNER JOIN posts_quotes
        ON post.id = posts_quotes."quotedId"
        WHERE posts_quotes."quoterId" = $1 ORDER BY post.id`,
      [this.id.suid]);
      return rows.map(row => makePost(row));
    },

    async getQuotedCount() {
      const { rows } = await query(`SELECT count(*) FROM posts_quotes
        WHERE "quotedId"=$1`, [this.id.suid]);
      return parseInt(rows[0].count, 10);
    },
  };
};

const postSliceOpt = {
  select: 'SELECT * FROM post',
  before: before => `post.id < ${UID.parse(before).suid}`,
  after: after => `post.id < ${UID.parse(after).suid}`,
  order: 'ORDER BY post.id',
  desc: false,
  name: 'posts',
  make: makePost,
  toCursor: post => post.id.duid,
};

const PostModel = {

  async findById({ postId }) {
    const { rows } = await query(
      'SELECT * FROM post WHERE id=$1',
      [UID.parse(postId).suid],
    );
    if ((rows || []).length === 0) {
      throw NotFoundError(`can't find post ${postId}`);
    }
    return makePost(rows[0]);
  },

  async findThreadPosts({ threadId, query: sq }) {
    const slice = await querySlice(sq, {
      ...postSliceOpt,
      where: 'WHERE post."threadId" = $1',
      params: [UID.parse(threadId).suid],
    });
    return slice;
  },

  async getThreadReplyCount({ threadId }) {
    const { rows } = await query(
      'SELECT count(*) FROM post WHERE "threadId"=$1',
      [UID.parse(threadId).suid],
    );
    return parseInt(rows[0].count || '0', 10);
  },

  async getThreadCatelog({ threadId }) {
    const { rows } = await query(
      `SELECT id "postId", "createdAt" "createdAt"
      FROM post WHERE "threadId"=$1 ORDER BY id`,
      [UID.parse(threadId).suid],
    );
    return (rows || []).map(row => ({
      postId: UID.parse(row.postId).duid,
      createdAt: row.createdAt,
    }));
  },

  async findUserPosts({ user, query: sq }) {
    const slice = await querySlice(sq, {
      ...postSliceOpt,
      where: 'WHERE post."userId" = $1',
      params: [user.id],
    });
    return slice;
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
    const raw = {
      id: await UID.new(),
      threadId: UID.parse(input.threadId),
      anonymous: input.anonymous,
      userId: user.id,
      userName: null,
      anonymousId: null,
      content: input.content,
      quoteIds: (input.quoteIds || []).map(qid => UID.parse(qid)),
    };
    let newPost;
    await doTransaction(async (txn) => {
      const { rows: ts } = await query(
        'SELECT locked as locked FROM thread WHERE id=$1', [raw.threadId.suid],
      );
      if ((ts || []).length === 0) {
        throw new NotFoundError(`thread ${raw.threadId.duid} not found`);
      }
      if (ts[0].locked) {
        throw new ParamsError(`thread ${raw.threadId.duid} is locked`);
      }

      if (input.anonymous) {
        raw.anonymousId = await UserAidModel.getAid({
          txn, userId: user.id, threadId: raw.threadId,
        });
      } else {
        raw.userName = user.name;
      }
      const { rows } = await txn.query(`INSERT INTO post
        (id, "threadId", anonymous, "userId", "userName", "anonymousId", content)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
      [raw.id.suid, raw.threadId.suid, raw.anonymous, raw.userId, raw.userName,
        raw.anonymousId && raw.anonymousId.suid, raw.content]);
      newPost = makePost(rows[0]);
      if (raw.quoteIds.length > 0) {
        await Promise.all(raw.quoteIds.map(qid => txn.query(`INSERT
           INTO posts_quotes ("quoterId", "quotedId")
           VALUES ($1, $2)`, [newPost.id.suid, qid.suid])));
        await NotificationModel.newQuotedNoti({
          txn,
          threadId: raw.threadId,
          postId: raw.id,
          quotedIds: raw.quoteIds,
        });
      }
      await NotificationModel.newRepliedNoti({ txn, threadId: raw.threadId });
    });
    return newPost;
  },

  async block({ ctx, postId }) {
    ctx.auth.ensurePermission(ACTION.BLOCK_POST);
    await query(
      'UPDATE post SET blocked=$1 WHERE id=$2',
      [true, UID.parse(postId).suid],
    );
    return true;
  },

};

export default PostModel;
export { blockedContent };
