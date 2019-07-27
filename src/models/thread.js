import { NotFoundError, ParamsError } from '~/utils/error';
import { query, doTransaction } from '~/utils/pg';
import { ACTION } from '~/models/user';
import AidModel from '~/models/aid';
import TagModel from '~/models/tag';
import querySlice from '~/models/slice';
import UID from '~/uid';

const blockedContent = '[此内容已被管理员屏蔽]';

const makeThread = function makeThread(raw) {
  return {
    id: UID.parse(raw.id),
    createdAt: raw.created_at,
    updatedAt: raw.updated_at,
    anonymous: raw.anonymous,
    author: raw.anonymous ? UID.parse(raw.anonymous_id).duid : raw.user_name,
    title: raw.title === '' ? '无题' : raw.title,
    content: raw.blocked ? blockedContent : raw.content,

    async getMainTag() {
      const { rows } = await query(`SELECT *
        FROM threads_tags inner join tag ON threads_tags.tag_name = tag.name
        WHERE threads_tags.thread_id = $1 AND tag.is_main = true`,
      [this.id.suid]);
      return rows[0].tag_name;
    },

    async getSubTags() {
      const { rows } = await query(`SELECT *
        FROM threads_tags INNER JOIN tag ON threads_tags.tag_name = tag.name
        WHERE threads_tags.thread_id = $1 AND tag.is_main = false`,
      [this.id.suid]);
      return (rows || []).map(row => row.tag_name);
    },

    blocked: raw.blocked,
    locked: raw.locked,
  };
};

const threadSliceOpt = {
  select: 'SELECT * FROM thread',
  before: before => `id > ${UID.parse(before).suid}`,
  after: after => `id < ${UID.parse(after).suid}`,
  order: 'ORDER BY id',
  desc: true,
  name: 'threads',
  make: makeThread,
  toCursor: thread => thread.id.duid,
};

const ThreadModel = {

  async findById({ txn, threadId }) {
    const id = UID.parse(threadId);
    const { rows } = await query('SELECT * FROM thread WHERE id=$1', [id.suid], txn);
    if ((rows || []).length === 0) {
      throw new NotFoundError(`cant find thread ${threadId}`);
    }
    return makeThread(rows[0]);
  },

  async findSlice({ tags, query: sq }) {
    const slice = await querySlice(sq, {
      ...threadSliceOpt,
      where: `WHERE id IN (
        SELECT thread_id FROM threads_tags
        WHERE threads_tags.tag_name=ANY($1)
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
        raw.anonymousId = await AidModel.getAid({ txn, userId: user.id, threadId: raw.id });
      } else {
        if (!user.name) {
          throw new ParamsError('you don\'t have a name');
        }
        raw.userName = user.name;
      }
      const { rows } = await txn.query(`INSERT INTO thread 
        (id, anonymous, user_id, user_name, anonymous_id, title, content)
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
        subTags: input.subTags || [],
      });
    });
    return newThread;
  },

  async findUserThreads({ user, query: sq }) {
    const slice = await querySlice(sq, {
      ...threadSliceOpt,
      where: 'WHERE user_id=$1',
      params: [user.id],
    });
    return slice;
  },

  async lock({ ctx, threadId }) {
    ctx.auth.ensurePermission(ACTION.LOCK_THREAD);
    await query(
      'UPDATE thread SET locked=$1 WHERE id=$2',
      [true, UID.parse(threadId).suid],
    );
  },

  async block({ ctx, threadId }) {
    ctx.auth.ensurePermission(ACTION.BLOCK_THREAD);
    await query(
      'UPDATE thread SET blocked=$1 WHERE id=$2',
      [true, UID.parse(threadId).suid],
    );
  },

};

export default ThreadModel;
export { blockedContent };
