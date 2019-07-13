import UID from '~/uid';
import { query } from '~/utils/pg';

const AidModel = {
  async getAid({ txn, userId, threadId }) {
    const tid = UID.parse(threadId);
    const aid = await UID.new();
    const { rows } = await query(`INSERT INTO anonymous_id
      (thread_id, user_id, anonymous_id) VALUES ($1, $2, $3)
      ON CONFLICT (thread_id, user_id) DO UPDATE SET updated_at=now() RETURNING *`,
    [tid.suid, userId, aid.suid], txn);
    return UID.parse(rows[0].anonymous_id);
  },
};

export default AidModel;
