import { query, doTransaction } from '~/utils/pg';
import { AuthError, ParamsError, PermissionError } from '~/utils/error';
import validator from '~/utils/validator';

// pgm.createTable('user', {
//   id: { type: 'serial', primaryKey: true },
//   email: { type: 'text', notNull: true, unique: true },
//   name: { type: 'text', unique: true },
//   role: { type: 'text', notNull: true, default: '' },
//   lastReadSystemNoti: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   lastReadRepliedNoti: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
//   lastReadQuotedNoti: { type: 'timestamp', notNull: true, default: pgm.func('now()') },
// });

const makeUser = function makeUser(raw) {
  return {
    id: raw.id,
    email: raw.email,
    async getTags() {
      const { rows } = query(`SELECT "tagName" 
        FROM users_tags WHERE "userId" = $1`, [raw.id]);
      return rows.map(row => row.tagName);
    },
    role: raw.role,
    readNotiTime: {
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
  PUB_POST: 'PUB_POST',
  PUB_THREAD: 'PUB_THREAD',
};

const ROLE = {
  MOD: 'mod',
  BANNED: 'banned',
};

const UserModel = {

  async findByEmail({ email }) {
    const { rows } = await query('SELECT * FROM user WHERE email = $1', [email]);
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
      const { rows } = await query('INSERT INTO user (email) VALUES ($1) RETURNING *', [email]);
      // TODO: send welcome message
      [user] = rows;
    }
    return {
      signedInUser() {
        return makeUser(user);
      },
      ensurePermission(action) {
        if (!user) {
          throw AuthError('not login');
        }
        if (user.role !== ROLE.MOD) {
          throw PermissionError('you are not mod');
        }
        if (!ACTION[action]) {
          throw ParamsError(`unknown action ${action}`);
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

  async addSubbedTags({ ctx, tags }) {
    const user = ctx.auth.signedInUser();
    await doTransaction(async (client) => {
      await Promise.all(tags.map(tag => client.query(
        'INSERT INTO users_tags ("userId", "tagName") VALUES ($1, $2) ON CONFLICT DO NOTHING',
        [user.id, tag],
      )));
    });
  },

  async delSubbedTags({ ctx, tags }) {
    const user = ctx.auth.signedInUser();
    await query('DELETE FROM users_tags WHERE "userId" = $1 "tagName" IN $2', [user.id, tags]);
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
export { ACTION };
