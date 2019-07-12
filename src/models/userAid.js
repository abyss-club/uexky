import UID from '~/uid';
import { query } from '~/utils/pg';

// pgm.createTable('anonymous_id', {
//   id: { type: 'serial', primaryKey: true },
//   createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   threadId: { type: 'bigint', notNull: true },
//   userId: { type: 'integer', notNull: true, references: 'public.user(id)' },
//   anonymousId: { type: 'bigint', notNull: true, unique: true },
// });

const UserAidModel = {
  async getAid({ txn, userId, threadId }) {
    const tid = UID.parse(threadId);
    const aid = await UID.new();
    const { rows } = await query(`INSERT INTO anonymous_id
      ("threadId", "userId", "anonymousId") VALUES ($1, $2, $3)
      ON CONFLICT ("threadId", "userId") DO UPDATE SET "updatedAt"=now() RETURNING *`,
    [tid.suid, userId, aid.suid], txn);
    return UID.parse(rows[0].anonymousId);
  },
};

export default UserAidModel;
