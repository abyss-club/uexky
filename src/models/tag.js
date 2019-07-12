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
    if (!isNew) {
      ctx.auth.ensurePermission({ action: ACTION.EDIT_TAG });
    }
    const { rows: mainTags } = await query(
      'SELECT name FROM tag WHERE "isMain" = true AND name = ANY ($1)',
      [[mainTag, ...subTags]], txn,
    );
    if (mainTags.length !== 1) {
      throw new ParamsError('you must specified one and only one main tag');
    } else if (mainTags[0].name !== mainTag) {
      throw new ParamsError(`${mainTag} is not main tag`);
    }
    if (!isNew) {
      await query('DELETE FROM threads_tags WHERE "threadId" = $1', [id.suid], txn);
    }
    await Promise.all([
      ...(subTags || []).map(
        tag => query(`INSERT INTO tag (name) VALUES ($1)
        ON CONFLICT (name) DO UPDATE SET "updatedAt" = now()`, [tag], txn),
      ),
      ...[mainTag, ...(subTags || [])].map(
        tag => query(`INSERT INTO threads_tags ("threadId", "tagName")
        VALUES ($1, $2) ON CONFLICT ("threadId", "tagName") DO UPDATE SET "updatedAt"=now()`,
        [id.suid, tag], txn),
      ),
      ...(subTags || []).map(subTag => query(`
        INSERT INTO tags_main_tags (name, "belongsTo") VALUES ($1, $2)
        ON CONFLICT (name, "belongsTo") DO UPDATE SET "updatedAt"=now()`,
      [subTag, mainTag], txn)),
    ]);
  },

};

export default TagModel;
