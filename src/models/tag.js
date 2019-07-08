import { ParamsError } from '~/utils/error';
import UID from '~/uid';
import { query } from '~/utils/pg';
import { ACTION } from '~/models/user';

const makeTag = function makeTag(raw) {
  return {
    name: raw.name,
    isMain: raw.isMain,
    async belongsTo() {
      if (raw.isMain) {
        return [];
      }
      const { rows } = query(
        'SELECT belongsTo FROM tags_main_tags WHERE name=$1', [raw],
      );
      return rows.map(row => row.belongsTo);
    },
  };
};

const TagModel = {

  async getMainTags() {
    const { rows } = await query('SELECT name FROM tag WHERE isMain=true');
    return rows.map(row => row.name);
  },

  async findTags({ query: qtext, limit = 10 }) {
    if (qtext) {
      const { rows } = await query(`SELECT (name, isMain) FROM tag
        WHERE name LIKE %$1% ORDER BY "updatedAt" DESC LIMIT $2`, [qtext, limit]);
      return rows.map(row => makeTag(row));
    }
    const { rows } = await query(
      'SELECT (name, isMain) FROM tag ORDER BY "updatedAt" DESC LIMIT $1', [limit],
    );
    return rows.map(row => makeTag(row));
  },

  async setThreadTags({
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
    await Promise.all([
      ...subTags.forEach(
        tag => q(`INSERT INTO tag (name) VALUES ($1)
        ON CONFLICT UPDATE "updatedAt" = now()`, [tag]),
      ),
      ...[mainTag, ...subTags].forEach(
        tag => q(`INSERT INTO threads_tags ("threadId", "tagName")
        VALUES ($1, $2) ON CONFLICT UPDATE SET "updatedAt" = now()`,
        [id.suid, tag]),
      ),
      ...subTags.map(subTag => query(`
        INSERT INTO tags_main_tags (name, belongsTo)
        VALUES ($1, $2) ON CONFLICT DO NOTHING`,
      [subTag, mainTag])),
    ]);
  },

};

export default TagModel;
