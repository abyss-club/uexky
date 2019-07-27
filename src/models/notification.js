import querySlice from '~/models/slice';
import { ParamsError } from '~/utils/error';
import { query } from '~/utils/pg';
import UID from '~/uid';

const USER_GROUPS = { ALL_USER: 'all_user' };
const NOTI_TYPES = {
  SYSTEM: 'system',
  REPLIED: 'replied',
  QUOTED: 'quoted',
};
const isValidType = (type) => {
  const types = Object.keys(NOTI_TYPES);
  for (let i = 0; i < types.length; i += 1) {
    if (NOTI_TYPES[types[i]] === type) {
      return true;
    }
  }
  return false;
};

// content
// system {
//   title string
//   content string
// }
// replied {
//   threadId bigint
// }
// quoted {
//   threadId bigint(string)
//   quotedId bigint(string)
//   postId bigint(string
// }

const makeNoti = (raw, type, user) => {
  const base = {
    id: raw.id,
    key: raw.key,
    eventTime: raw.updated_at,
    hasRead: user.lastReadNoti[type] >= raw.id,
  };
  if (type === NOTI_TYPES.SYSTEM) {
    return {
      ...base,
      type: NOTI_TYPES.SYSTEM,
      title: raw.content.title,
      content: raw.content.content,
    };
  } if (type === NOTI_TYPES.REPLIED) {
    return {
      ...base,
      type: NOTI_TYPES.REPLIED,
      threadId: UID.parse(raw.content.threadId),
    };
  } if (type === NOTI_TYPES.QUOTED) {
    return {
      ...base,
      type: NOTI_TYPES.QUOTED,
      threadId: UID.parse(raw.content.threadId),
      quotedId: UID.parse(raw.content.quotedId),
      postId: UID.parse(raw.content.postId),
    };
  }
  throw new ParamsError(`unknown notification type: ${type}`);
};

const notiSliceOpt = {
  select: 'SELECT * FROM notification',
  before: before => `id > ${parseInt(before, 10)}`,
  after: after => `id < ${parseInt(after, 10)}`,
  order: 'ORDER BY id',
  desc: true,
  toCursor: noti => noti.id.toString(),
};

const NotificationModel = {

  async getUnreadCount({ ctx, type }) {
    if (!isValidType(type)) {
      throw new ParamsError(`unknown notification type: ${type}`);
    }
    const user = ctx.auth.signedInUser();
    const sql = `
      SELECT count(*) FROM notification
      WHERE (send_to=$1 OR send_to_group=$2)
      AND id > $3 AND type=$4 AND updated_at >= $5`;
    const { rows } = await query(sql, [user.id, USER_GROUPS.ALL_USER,
      user.lastReadNoti[type], type, user.createdAt]);
    return parseInt(rows[0].count, 10);
  },

  async findNotiSlice({ ctx, type, query: sq }) {
    if (!isValidType(type)) {
      throw new ParamsError(`unknown notification type: ${type}`);
    }
    const user = ctx.auth.signedInUser();
    const opt = {
      ...notiSliceOpt,
      where: 'WHERE (send_to=$1 OR send_to_group=$2) AND type=$3 AND updated_at >= $4',
      params: [user.id, USER_GROUPS.ALL_USER, type, user.createdAt],
      name: type,
      make: raw => makeNoti(raw, type, user),
    };
    const slice = await querySlice(sq, opt);
    if (slice.sliceInfo.lastCursor !== '') {
      await query(
        `UPDATE public.user SET "last_read_${type}_noti"=$1`,
        [parseInt(slice.sliceInfo.lastCursor, 10)],
      );
    }
    return slice;
  },

  async newSystemNoti({
    sendTo, sendToGroup, title, content,
  }) {
    const key = `system:${(await UID.new()).duid}`;
    if (sendTo) {
      await query(`INSERT INTO notification
      (key, type, send_to, content) VALUES ($1, $2, $3, $4)`,
      [key, NOTI_TYPES.SYSTEM, sendTo, { title, content }]);
    } if (sendToGroup) {
      await query(`INSERT INTO notification
      (key, type, send_to_group, content) VALUES ($1, $2, $3, $4)`,
      [key, NOTI_TYPES.SYSTEM, sendToGroup, { title, content }]);
    }
  },

  async newRepliedNoti({ txn, threadId }) {
    const tid = UID.parse(threadId);
    const key = `replied:${tid.duid}`;
    const content = { threadId: threadId.suid.toString() };
    const { rows } = await query(
      'SELECT * FROM thread WHERE id=$1', [tid.suid], txn,
    );
    await query(`INSERT INTO notification
      (key, type, send_to, content) VALUES ($1, $2, $3, $4)
      ON CONFLICT (key) DO UPDATE SET updated_at=now() RETURNING *`,
    [key, NOTI_TYPES.REPLIED, rows[0].user_id, content], txn);
  },

  async newQuotedNoti({
    txn, threadId, postId, quotedIds,
  }) {
    const { rows } = await query(
      'SELECT id, user_id FROM post WHERE id=ANY($1)',
      [quotedIds.map(qid => qid.suid)], txn,
    );
    await Promise.all(rows.map((row) => {
      const pid = UID.parse(postId);
      const qid = UID.parse(row.id);
      const key = `${NOTI_TYPES.QUOTED}:${qid.duid}:${pid.duid}`;
      const content = {
        threadId: threadId.suid.toString(),
        quotedId: row.id,
        postId: postId.suid.toString(),
      };
      return query(
        `INSERT INTO notification (key, type, send_to, content)
        VALUES ($1, $2, $3, $4)`,
        [key, NOTI_TYPES.QUOTED, row.user_id, content], txn,
      );
    }));
  },
};

export default NotificationModel;
export { NOTI_TYPES, USER_GROUPS, isValidType };
