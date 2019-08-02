import { query, doTransaction } from '~/utils/pg';
import { AuthError, ParamsError, PermissionError } from '~/utils/error';
import validator from '~/utils/validator';
import NotificationModel from '~/models/notification';
import welcome from '~/models/welcome';

const makeUser = function makeUser(raw) {
  const getTagsSql = 'SELECT tag_name FROM users_tags WHERE user_id=$1';
  return {
    id: raw.id,
    createdAt: raw.created_at,
    name: raw.name,
    email: raw.email,
    async getTags() {
      const { rows } = await query(getTagsSql, [raw.id]);
      return (rows || []).map(row => row.tag_name);
    },
    role: raw.role,
    lastReadNoti: {
      system: raw.last_read_system_noti,
      replied: raw.last_read_replied_noti,
      quoted: raw.last_read_quoted_noti,
    },
  };
};

const ACTION = {
  BAN_USER: 'BAN_USER',
  BLOCK_POST: 'BLOCK_POST',
  LOCK_THREAD: 'LOCK_THREAD',
  BLOCK_THREAD: 'BLOCK_THREAD',
  EDIT_TAG: 'EDIT_TAG',
  EDIT_SETTING: 'EDIT_SETTING',
  PUB_POST: 'PUB_POST',
  PUB_THREAD: 'PUB_THREAD',
};

const ROLE = {
  MOD: 'mod',
  BANNED: 'banned',
  ADMIN: 'admin',
};

const actionRole = {
  [ACTION.BAN_USER]: 'mod',
  [ACTION.BLOCK_POST]: 'mod',
  [ACTION.LOCK_THREAD]: 'mod',
  [ACTION.BLOCK_THREAD]: 'mod',
  [ACTION.EDIT_TAG]: 'mod',
  [ACTION.EDIT_SETTING]: 'admin',
};

const newUser = async (email) => {
  const userSql = 'INSERT INTO public.user (email) VALUES ($1) RETURNING *';
  const tagSql = 'INSERT INTO users_tags (user_id, tag_name) VALUES ($1, $2)';
  const mainTagSql = 'SELECT name FROM tag WHERE is_main=true ORDER BY created_at';
  const values = [email];
  const { rows } = await query(userSql, values);
  const [user] = rows;
  // send welcome message
  await NotificationModel.newSystemNoti({ sendTo: user.id, ...welcome });

  // set default tags by main tags
  const { rows: tags } = await query(mainTagSql);
  await Promise.all(tags.map(tag => query(tagSql, [user.id, tag.name])));
  return makeUser(user);
};

const UserModel = {

  async findByEmail({ email }) {
    const sql = 'SELECT * FROM public.user WHERE email=$1';
    const { rows } = await query(sql, [email]);
    if (rows.length !== 0) {
      return makeUser(rows[0]);
    }
    return null;
  },

  // if not specified email, return guest auth helper,
  // otherwise login this email, and return signed in auth helper.
  async authContext({ email }) {
    // guest
    if (!email) {
      return {
        isSignedIn: false,
        signedInUser: () => { throw new AuthError('you are not signed in'); },
        ensurePermission: () => { throw new PermissionError('permission denyed'); },
      };
    }
    // signed in
    if (!validator.isValidEmail(email)) {
      throw new ParamsError(`Invalid Email: ${email}`);
    }
    let user = await this.findByEmail({ email });
    if (!user) { // new user
      user = await newUser(email);
    }
    return {
      isSignedIn: true,
      signedInUser() {
        return user;
      },
      ensurePermission(action) {
        if (!user) {
          throw new AuthError('not login');
        }
        const ar = actionRole[action] || '';
        if (ar === ROLE.ADMIN) {
          if (user.role !== ROLE.ADMIN) {
            throw new PermissionError('you are not admin');
          }
          return;
        } if (ar === ROLE.MOD) {
          if ((user.role !== ROLE.ADMIN) && (user.role !== ROLE.MOD)) {
            throw new PermissionError('you are not mod');
          }
          return;
        } if (user.role === ROLE.BANNED) {
          throw new PermissionError('you are banned');
        }
      },
    };
  },

  async setName({ ctx, name }) {
    if (!validator.isUnicodeLength(name, { max: 15 })) {
      throw new ParamsError('Max length of username is 15.');
    }
    const user = ctx.auth.signedInUser();
    if ((user.name || '') !== '') {
      throw new ParamsError('Name can only be set once.');
    }
    const sql = 'UPDATE public.user SET name=$1 WHERE email=$2';
    try {
      await query(sql, [name, user.email]);
    } catch (error) {
      if (error.message.includes('duplicate key')) {
        throw new ParamsError('duplicate name');
      }
      throw error;
    }
    user.name = name;
    return user;
  },

  async addSubbedTag({ ctx, tag }) {
    const user = ctx.auth.signedInUser();
    const sql = 'INSERT INTO users_tags (user_id, tag_name) VALUES ($1, $2)';
    await query(sql, [user.id, tag]);
  },

  async delSubbedTag({ ctx, tag }) {
    const user = ctx.auth.signedInUser();
    const sql = 'DELETE FROM users_tags WHERE user_id=$1 AND tag_name=$2';
    await query(sql, [user.id, tag]);
  },

  async syncTags({ ctx, tags }) {
    const user = ctx.auth.signedInUser();
    const sql = `INSERT INTO
                 users_tags (user_id, tag_name)
                 VALUES ($1, $2) ON CONFLICT DO NOTHING`;
    await doTransaction(async (txn) => {
      await txn.query('DELETE FROM users_tags WHERE user_id=$1', [user.id]);
      // find valid tags
      const { rows } = await txn.query(
        'SELECT count(*) FROM tag WHERE name=ANY($1)', [tags],
      );
      if (parseInt(rows[0].count, 10) !== tags.length) {
        throw new ParamsError('invalid tags');
      }
      await Promise.all(tags.map(async (tag) => {
        await txn.query(sql, [user.id, tag]);
      }));
    });
  },

  async banUser({ ctx, userId }) {
    ctx.auth.ensurePermission({ action: ACTION.BAN_USER });
    await query('UPDATE public.user SET role=$1 WHERE id=$2', [ROLE.BANNED, userId]);
  },
};

export default UserModel;
export { ACTION, ROLE };
