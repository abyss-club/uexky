import { NotFoundError, ParamsError } from '~/utils/error';
import { query, doTransaction } from '~/utils/pg';
import { ACTION } from '~/models/user';
import UserAidModel from '~/models/userAid';
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
      return rows.map(row => row.tagName);
    },

    async getReplyCount() {
      const { rows } = await query(`SELECT count(*) FROM post
        WHERE "threadId"=$1`, [this.id.suid]);
      return rows[0].count;
    },

    async getCatelog() {
      const { rows } = await query(`SELECT ("postId, createdAt")
      FROM post WHERE "threadId=$1" ORDER BY id DESC`, [this.id.suid]);
      return rows;
    },

    blocked: raw.blocked,
    locked: raw.locked,
  };
};

const ThreadModel = {

  async findById({ threadId }) {
    const id = UID.parse(threadId);
    const { rows } = await query('SELECT * FROM thread WHERE id=$1', [id]);
    if ((rows || []).length === 0) {
      throw NotFoundError(`cant find thread ${threadId}`);
    }
    return makeThread(rows[0]);
  },

  async findSlice({ tags, query: sq }) {
    const sql = [
      'SELECT *',
      'FROM thread inner join threads_tags',
      'ON thread.id = threads_tags."threadId"',
      'WHERE threads_tags."tagName" IN $1',
      sq.before && 'AND thread.id > $2',
      sq.after && 'AND thread.id < $3',
      'ORDER BY thread.id DESC',
      `LIMIT ${sq.limit || 0}`,
    ].join(' ');
    const { rows } = await query(sql, [tags, sq.before, sq.after]);
    return rows.map(row => makeThread(row));
  },

  async new({ ctx, thread: input }) {
    ctx.auth.ensurePermission(ACTION.PUB_THREAD);
    const user = ctx.auth.signedInUser();
    const raw = {
      id: UID.new(),
      anonymous: input.anonymous,
      userId: input.userId,
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
        raw.anonymousId && raw.anonymous.suid, raw.title, raw.content]);
      newThread = makeThread(rows[0]);
      await this.setTags({
        ctx, txn, isNew: true, mainTag: input.mainTag, subTags: input.subTags,
      });
    });
    return newThread;
  },

  async findUserThreads({ user, query: sq }) {
    const sql = [
      'SELECT *',
      'FROM thread inner join threads_tags',
      'ON thread.id = threads_tags."threadId"',
      'WHERE thread."userId" == $1',
      sq.before && 'AND thread.id > $2',
      sq.after && 'AND thread.id < $3',
      'ORDERED BY thread."updatedAt" DESC',
      `LIMIT ${sq.limit || 0}`,
    ].join(' ');
    const { threads } = await query(sql, [user.id, sq.before, sq.after]);
    return threads;
  },

  async lockThread({ ctx, threadId }) {
    ctx.auth.ensurePermission(ACTION.LOCK_THREAD);
    await query('UPDATE thread SET lock=$1 WHERE id=$2', [true, UID.parse(threadId).suid]);
  },

  async blockThread({ ctx, threadId }) {
    ctx.auth.ensurePermission(ACTION.BLOCK_THREAD);
    await query('UPDATE thread SET block=$1 WHERE id=$2', [true, UID.parse(threadId).suid]);
  },

  async setTags({
    ctx, txn, isNew, threadId, mainTag, subTags,
  }) {
    const id = UID.parse(threadId);
    const q = txn ? txn.query : query;
    if (!isNew) {
      ctx.auth.ensurePermission({ action: ACTION.EDIT_TAG });
    }
    const { rows: mainTags } = q(
      'SELECT name FROM tag WHERE "isMain" = true AND name IN $1',
      [[mainTag, ...subTags]],
    );
    if (mainTags.length !== 1) {
      throw ParamsError('you must specified one and only one main tag');
    } else if (mainTags[0].name !== mainTag) {
      throw ParamsError(`${mainTag} is not main tag`);
    }
    if (!isNew) {
      await q('DELETE FROM threads_tags WHERE "threadId" = $1', [id.suid]);
    }
    await Promise.all(subTags.forEach(
      tag => q(`INSERT INTO tag (name) VALUES ($1)
        ON CONFLICT UPDATE "updatedAt" = now()`, [tag]),
    ));
    await Promise.all([mainTag, ...subTags].forEach(
      tag => q(`INSERT INTO threads_tags ("threadId", "tagName")
        VALUES ($1, $2) ON CONFLICT UPDATE SET "updatedAt" = now()`,
      [id.suid, tag]),
    ));
  },
};

export default ThreadModel;
