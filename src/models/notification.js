import { ParamsError } from '~/utils/error';
import { query } from '~/utils/pg';
import UID from '~/uid';


const USER_GROUPS = { ALL_USER: 'all_user' };
const NOTI_TYPES = {
  SYSTEM: 'system',
  REPLIED: 'replied',
  QUOTED: 'quoted',
};
const isValidType = type => (NOTI_TYPES[type] || '') === '';

const makeNoti = (raw, type, user) => {
  const base = {
    id: raw.id,
    type: NOTI_TYPES.SYSTEM,
    eventTime: raw.createdAt,
    hasRead: user.readNotiTime.system >= raw.createdAt,
  };
  if (type === NOTI_TYPES.SYSTEM) {
    return {
      ...base,
      title: raw.content.title,
      content: raw.content.content,
    };
  } if (type === NOTI_TYPES.REPLIED) {
    return {
      ...base,
      threadId: raw.content.threadId,
      quotedId: raw.content.quotedId,
      postId: raw.content.postId,
    };
  } if (type === NOTI_TYPES.QUOTED) {
    return {
      ...base,
      quotedId: user.content.quotedId,
      postId: user.content.postId,
    };
  }
  throw new ParamsError(`unknown notification type: ${type}`);
};

const NotificationModel = {

  async getUnreadCount({ ctx, type }) {
    if (!isValidType(type)) {
      throw new ParamsError(`unknown notification type: ${type}`);
    }
    const user = ctx.auth.signedInUser();
    const { rows } = await query(
      `SELECT count(*) FROM notification
      WHERE "updatedAt" < $1`,
      [user.readNotiTime[type]],
    );
    return rows[0].count;
  },

  async findNotiSlice({ ctx, type, query: sq }) {
    if (!isValidType(type)) {
      throw new ParamsError(`unknown notification type: ${type}`);
    }
    const user = ctx.auth.signedInUser();
    const { rows } = await query([
      'SELECT * FROM notification',
      `WHERE ("sendTo"=$1 OR "sendToGroup"=${USER_GROUPS.ALL_USER}')`,
      'AND type=$2',
      sq.before && `AND id < ${sq.before}`,
      sq.after && `AND id > ${sq.after}`,
      `ORDER BY id DESC LIMIT ${sq.limit || 0}`,
    ].join(' '), [user.id, type]);
    return rows.map(raw => makeNoti(raw, type, user));
  },

  async newSystemNoti({
    sendTo, sendToGroup, title, content,
  }) {
    const id = `system:${(await UID.new()).duid}`;
    if (sendTo) {
      await query(`INSERT INTO notification
      (id, type, sendTo, content) VALUES ($1, $2, $3, $4)`,
      [id, NOTI_TYPES.SYSTEM, sendTo, { title, content }]);
    } if (sendToGroup) {
      await query(`INSERT INTO notification
      (id, type, sendToGroup, content) VALUES ($1, $2, $3, $4)`,
      [id, NOTI_TYPES.SYSTEM, sendToGroup, { title, content }]);
    }
  },

  async newRepliedNoti({ txn, threadId }) {
    const tid = UID.parse(threadId);
    const id = `replied:${tid.duid}`;
    const { rows } = await query('SELECT * FROM thread WHERE id=$1', [tid.suid], txn);
    await query(`INSERT INTO notification
      (id, type, "sendTo", content) VALUES ($1, $2, $3, $4)
      ON CONFLICT (id) DO UPDATE SET "updatedAt"=now() RETURNING *`,
    [id, NOTI_TYPES.REPLIED, rows[0].userId, { threadId: threadId.suid.toString() }], txn);
  },

  async newQuotedNoti({
    txn, threadId, postId, quotedIds,
  }) {
    const { rows } = await query(
      'SELECT id id, "userId" "userId" FROM post WHERE id=ANY($1)',
      [quotedIds.map(qid => qid.suid)], txn,
    );
    await Promise.all(rows.map((row) => {
      const pid = UID.parse(postId);
      const qid = UID.parse(row.id);
      const id = `${NOTI_TYPES.QUOTED}:${qid.duid}:${pid.duid}`;
      const content = {
        threadId: threadId.suid.toString(),
        quotedId: row.id,
        postId: postId.suid.toString(),
      };
      return query(
        `INSERT INTO notification (id, type, "sendTo", content)
        VALUES ($1, $2, $3, $4)`,
        [id, NOTI_TYPES.QUOTED, row.userId, content], txn,
      );
    }));
  },
};

export default NotificationModel;
export { NOTI_TYPES };
