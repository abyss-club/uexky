import UID from '~/uid';
import { query } from '~/utils/pg';

// pgm.createTable('anonymous_id', {
//   id: { type: 'serial', primaryKey: true },
//   createdAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   updatedAt: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   threadId: { type: 'bigint', notNull: true, references: 'thread(id)' },
//   userId: { type: 'integer', notNull: true, references: 'public.user(id)' },
//   anonymousId: { type: 'bigint', notNull: true, unique: true },
// });

const UserAidModel = {
  async new({ txn, userId, threadId }) {
    const q = txn ? txn.query : query;
    const tid = UID.parse(threadId);
    const aid = await UID.new();
    const { row } = await q(`INSERT INTO anonymous_id
      ("threadId", "userId", "anonymousId") VALUES ($1, $2, $3)
      ON CONFLICT DO NOTHING RETURNING *`,
    [tid.suid, userId, aid.suid]);
    return UID.parse(row.anonymousId);
  },
};

export default UserAidModel;
