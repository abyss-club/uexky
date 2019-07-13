import { query, doTransaction } from '~/utils/pg';
import { AuthError, ParamsError, PermissionError } from '~/utils/error';
import validator from '~/utils/validator';

const makeUser = function makeUser(raw) {
  return {
    id: raw.id,
    name: raw.name,
    email: raw.email,
    async getTags() {
      const { rows } = await query(
        'SELECT "tagName" FROM users_tags WHERE "userId"=$1',
        [raw.id],
      );
      return (rows || []).map(row => row.tagName);
    },
    role: raw.role,
    lastReadNoti: {
      system: raw.lastReadSystemNoti,
      replied: raw.lastReadRepliedNoti,
      quoted: raw.lastReadQuotedNoti,
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

const UserModel = {

  async findByEmail({ email }) {
    const { rows } = await query('SELECT * FROM public.user WHERE email=$1', [email]);
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
        signedInUser: () => { throw AuthError('you are not signed in'); },
        ensurePermission: () => { throw PermissionError('permission denyed'); },
      };
    }

    // signed in
    if (!validator.isValidEmail(email)) {
      throw new ParamsError(`Invalid Email: ${email}`);
    }
    let user = await this.findByEmail({ email });
    if (!user) { // new user
      const text = 'INSERT INTO public.user (email) VALUES ($1) RETURNING *';
      const values = [email];
      const { rows } = await query(text, values);
      // TODO: send welcome message
      [user] = rows;
    }
    return {
      signedInUser() {
        return user;
      },
      ensurePermission(action) {
        if (!user) {
          throw AuthError('not login');
        }
        const ar = actionRole[action] || '';
        if (ar === ROLE.ADMIN) {
          if (user.role !== ROLE.ADMIN) {
            throw PermissionError('you are not admin');
          }
          return;
        } if (ar === ROLE.MOD) {
          if ((user.role !== ROLE.ADMIN) && (user.role !== ROLE.MOD)) {
            throw PermissionError('you are not mod');
          }
          return;
        } if (user.role === ROLE.BANNED) {
          throw PermissionError('you are banned');
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
    await query('UPDATE public.user SET name=$1 WHERE email=$2', [name, user.email]);
    user.name = name;
    return user;
  },

  async addSubbedTag({ ctx, tag }) {
    const user = ctx.auth.signedInUser();
    await query(
      'INSERT INTO users_tags ("userId", "tagName") VALUES ($1, $2)',
      [user.id, tag],
    );
  },

  async delSubbedTag({ ctx, tag }) {
    const user = ctx.auth.signedInUser();
    await query(
      'DELETE FROM users_tags WHERE "userId"=$1 AND "tagName"=$2',
      [user.id, tag],
    );
  },

  async syncTags({ ctx, tags }) {
    const user = ctx.auth.signedInUser();
    await doTransaction(async (txn) => {
      await txn.query('DELETE FROM users_tags WHERE "userId"=$1', [user.id]);
      await Promise.all(tags.map(tag => txn.query(`INSERT INTO
        users_tags ("userId", "tagName")
        VALUES ($1, $2) ON CONFLICT DO NOTHING`,
      [user.id, tag])));
    });
  },

  async banUser({ ctx, userId }) {
    ctx.auth.ensurePermission({ action: ACTION.BAN_USER });
    await query('UPDATE public.user SET role=$1 WHERE id=$2', [ROLE.BANNED, userId]);
  },
};

export default UserModel;
export { ACTION, ROLE };
