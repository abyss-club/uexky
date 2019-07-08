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
  throw ParamsError(`unknown notification type: ${type}`);
};

const NotificationModel = {

  async getUnreadCount({ ctx, type }) {
    if (!isValidType(type)) {
      throw ParamsError(`unknown notification type: ${type}`);
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
      throw ParamsError(`unknown notification type: ${type}`);
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
    if (sendTo) {
      await query(`INSERT INTO notification
      (type, sendTo, content) VALUES ($1, $2, $3)`,
      [NOTI_TYPES.SYSTEM, sendTo, { title, content }]);
    } if (sendToGroup) {
      await query(`INSERT INTO notification
      (type, sendToGroup, content) VALUES ($1, $2, $3)`,
      [NOTI_TYPES.SYSTEM, sendToGroup, { title, content }]);
    }
  },

  async newRepliedNoti({ txn, threadId }) {
    const q = txn ? txn.query : query;
    const { rows } = q('SELECT "userId" FROM thread WHERE id=$1',
      [UID.parse(threadId).suid]);
    await q(`INSERT INTO notification
      (type, sendTo, content) VALUES ($1, $2, $3)
      ON CONFLICT UPDATE SET "updatedAt"=now()`,
    [NOTI_TYPES.REPLIED, rows[0].userId, { threadId: threadId.suid }]);
  },

  async newQuotedNoti({
    txn, threadId, postId, quotedIds,
  }) {
    const q = txn ? txn.query : query;
    const { rows } = await q(
      'SELECT (id, "userId") FROM posts WHERE id IN $1', [quotedIds],
    );
    await Promise.all(rows.map(row => q(`INSERT INTO notification
      (type, sendTo, content) VALUES ($1, $2, $3)`,
    [NOTI_TYPES.QUOTED, row.userId, {
      threadId: threadId.suid,
      quotedId: row.id,
      postId: postId.suid,
    },
    ])));
  },
};

export default NotificationModel;
export { NOTI_TYPES };
