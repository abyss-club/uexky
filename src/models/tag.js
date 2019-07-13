import { ParamsError } from '~/utils/error';
import UID from '~/uid';
import { query } from '~/utils/pg';
import { ACTION } from '~/models/user';

const makeTag = function makeTag(raw) {
  return {
    name: raw.name,
    isMain: raw.is_main,
    async getBelongsTo() {
      if (raw.is_main) {
        return [];
      }
      const { rows } = await query(
        'SELECT * FROM tags_main_tags WHERE name=$1', [raw.name],
      );
      return (rows || []).map(row => row.belongs_to);
    },
  };
};

const TagModel = {

  async getMainTags() {
    const { rows } = await query(
      'SELECT name FROM tag WHERE is_main=true ORDER BY created_at',
    );
    return rows.map(row => row.name);
  },

  async findTags({ query: qtext, limit = 10 }) {
    if (qtext) {
      const { rows } = await query(
        `SELECT name, is_main FROM tag
        WHERE name LIKE $1 ORDER BY updated_at DESC LIMIT $2`,
        [`%${qtext}%`, limit],
      );
      return rows.map(row => makeTag(row));
    }
    const { rows } = await query(
      'SELECT name, is_main FROM tag ORDER BY updated_at DESC LIMIT $1',
      [limit],
    );
    return rows.map(row => makeTag(row));
  },

  async setThreadTags({
    ctx, txn, isNew, threadId, mainTag, subTags,
  }) {
    const id = UID.parse(threadId);
    if (!isNew) {
      ctx.auth.ensurePermission({ action: ACTION.EDIT_TAG });
    }
    const { rows: mainTags } = await query(
      'SELECT name FROM tag WHERE is_main = true AND name = ANY ($1)',
      [[mainTag, ...subTags]], txn,
    );
    if (mainTags.length !== 1) {
      throw new ParamsError('you must specified one and only one main tag');
    } else if (mainTags[0].name !== mainTag) {
      throw new ParamsError(`${mainTag} is not main tag`);
    }
    if (!isNew) {
      await query('DELETE FROM threads_tags WHERE thread_id = $1', [id.suid], txn);
    }
    await Promise.all([
      ...[mainTag, ...(subTags || [])].map(
        tag => query(`INSERT INTO tag (name) VALUES ($1)
        ON CONFLICT (name) DO UPDATE SET updated_at = now()`, [tag], txn),
      ),
      ...[mainTag, ...(subTags || [])].map(
        tag => query(`INSERT INTO threads_tags (thread_id, tag_name)
        VALUES ($1, $2) ON CONFLICT (thread_id, tag_name) DO UPDATE SET updated_at=now()`,
        [id.suid, tag], txn),
      ),
      ...(subTags || []).map(subTag => query(`
        INSERT INTO tags_main_tags (name, belongs_to) VALUES ($1, $2)
        ON CONFLICT (name, belongs_to) DO UPDATE SET updated_at=now()`,
      [subTag, mainTag], txn)),
    ]);
  },

};

export default TagModel;
