import { NotFoundError } from '~/utils/error';
import { query, doTransaction } from '~/utils/pg';
import { ACTION } from '~/models/user';
import UserAidModel from '~/models/userAid';
import TagModel from '~/models/tag';
import querySlice from '~/models/slice';
import UID from '~/uid';

// pgm.createTable('thread', {
//   id: { type: 'bigint', primaryKey: true },
//   createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//
//   anonymous: { type: 'boolean', notNull: true },
//   userId: { type: 'integer', notNull: true, references: 'public.user(id)' },
//   userName: { type: 'varchar(16)', references: 'public.user(name)' },
//   anonymousId: { type: 'bigint' },
//
//   title: { type: 'text', default: '' },
//   content: { type: 'text', notNull: true },
//   locked: { type: 'bool', notNull: true, default: false },
//   blocked: { type: 'bool', notNull: true, default: false },
// });
// pgm.createTable('threads_tags', {
//   id: { type: 'serial', primaryKey: true },
//   createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },
//   tagName: { type: 'text', notNull: true, references: 'tag(name)' },
// });

const makeThread = function makeThread(raw) {
  return {
    id: UID.parse(raw.id),
    createdAt: raw.createdAt,
    updatedAt: raw.updatedAt,
    anonymous: raw.anonymous,
    author: raw.anonymous ? UID.parse(raw.anonymousId).duid : raw.userName,
    title: raw.title === '' ? '无题' : raw.title,
    content: raw.blocked ? '[此内容已被屏蔽]' : raw.content,

    async getMainTag() {
      const { rows } = await query(`SELECT *
        FROM threads_tags inner join tag ON threads_tags."tagName" = tag.name
        WHERE threads_tags."threadId" = $1 AND tag."isMain" = true`,
      [this.id.suid]);
      return rows[0].tagName;
    },

    async getSubTags() {
      const { rows } = await query(`SELECT *
        FROM threads_tags INNER JOIN tag ON threads_tags."tagName" = tag.name
        WHERE threads_tags."threadId" = $1 AND tag."isMain" = false`,
      [this.id.suid]);
      return (rows || []).map(row => row.tagName);
    },

    async getReplyCount() {
      const { rows } = await query(`SELECT count(*) FROM post
        WHERE "threadId"=$1`, [this.id.suid]);
      return parseInt(rows[0].count || '0', 10);
    },

    async getCatelog() {
      const { rows } = await query(`SELECT id "postId", "createdAt"
      FROM post WHERE "threadId"=$1 ORDER BY id DESC`, [this.id.suid]);
      return rows || [];
    },

    blocked: raw.blocked,
    locked: raw.locked,
  };
};

const threadSliceOpt = {
  select: 'SELECT * FROM thread',
  before: before => `thread.id > ${UID.parse(before).suid}`,
  after: after => `thread.id < ${UID.parse(after).suid}`,
  order: 'ORDER BY thread.id',
  desc: true,
  name: 'threads',
  make: makeThread,
  toCursor: thread => thread.id.duid,
};

const ThreadModel = {

  async findById({ threadId }) {
    const id = UID.parse(threadId);
    const { rows } = await query('SELECT * FROM thread WHERE id=$1', [id.suid]);
    if ((rows || []).length === 0) {
      throw new NotFoundError(`cant find thread ${threadId}`);
    }
    return makeThread(rows[0]);
  },

  async findSlice({ tags, query: sq }) {
    const slice = await querySlice(sq, {
      ...threadSliceOpt,
      where: `WHERE id IN (
        SELECT "threadId" FROM threads_tags
        WHERE threads_tags."tagName"=ANY($1)
      )`,
      params: [tags],
    });
    return slice;
  },

  async new({ ctx, thread: input }) {
    ctx.auth.ensurePermission(ACTION.PUB_THREAD);
    const user = ctx.auth.signedInUser();
    const threadId = await UID.new();
    const raw = {
      id: threadId,
      anonymous: input.anonymous,
      userId: user.id,
      userName: null,
      anonymousId: null,
      title: input.title,
      content: input.content,
    };
    let newThread;
    await doTransaction(async (txn) => {
      if (input.anonymous) {
        raw.anonymousId = await UserAidModel.getAid({ txn, userId: user.id, threadId: raw.id });
      } else {
        raw.userName = user.name;
      }
      const { rows } = await txn.query(`INSERT INTO thread 
        (id, anonymous, "userId", "userName", "anonymousId", title, content)
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
      [raw.id.suid, raw.anonymous, raw.userId, raw.userName,
        raw.anonymousId && raw.anonymousId.suid, raw.title, raw.content]);
      newThread = makeThread(rows[0]);
      await TagModel.setThreadTags({
        ctx,
        txn,
        isNew: true,
        threadId: newThread.id,
        mainTag: input.mainTag,
        subTags: input.subTags,
      });
    });
    return newThread;
  },

  async findUserThreads({ user, query: sq }) {
    const slice = await querySlice(sq, {
      ...threadSliceOpt,
      where: 'WHERE "userId"=$1',
      params: [user.id],
    });
    return slice;
  },

  async lockThread({ ctx, threadId }) {
    ctx.auth.ensurePermission(ACTION.LOCK_THREAD);
    await query('UPDATE thread SET lock=$1 WHERE id=$2', [true, UID.parse(threadId).suid]);
  },

  async blockThread({ ctx, threadId }) {
    ctx.auth.ensurePermission(ACTION.BLOCK_THREAD);
    await query('UPDATE thread SET block=$1 WHERE id=$2', [true, UID.parse(threadId).suid]);
  },

};

export default ThreadModel;
