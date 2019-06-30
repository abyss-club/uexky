import { query, doTransaction } from '~/utils/pg';
import { AuthError, ParamsError, PermissionError } from '~/utils/error';
import validator from '~/utils/validator';

// pgm.createTable('user', {
//   id: { type: 'serial', primaryKey: true },
//   email: { type: 'text', notNull: true, unique: true },
//   name: { type: 'text', unique: true },
//   role: { type: 'text' },
//   lastReadSystemNoti: { type: 'timestamp', default: 'now()' },
//   lastReadRepliedNoti: { type: 'timestamp', default: 'now()' },
//   lastReadQuotedNoti: { type: 'timestamp', default: 'now()' },
// });

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

function ensurePermission({ user, action }) {
  if (user.role !== ROLE.MOD) {
    throw PermissionError('you are not mod');
  }
  if (!ACTION[action]) {
    throw ParamsError(`unknown action ${action}`);
  }
}

async function authContext({ email }) {
  if (!email) {
    return {
      signedInUser: () => { throw AuthError('you are not signed in'); },
      ensurePermission: () => { throw PermissionError('permission denyed'); },
    };
  }
  if (!validator.isValidEmail(email)) {
    throw new ParamsError(`Invalid Email: ${email}`);
  }
  let user = await findByEmail({ email });
  if (!user) {
    const { rows } = await query('INSERT INTO user (email) VALUES ($1) RETURNING *', [email]);
    // TODO: send welcome message
    [user] = rows;
  }
  return {
    signedInUser: () => user,
    ensurePermission: action => ensurePermission({ user, action }),
  };
}

async function findById({ userId }) {
  const { rows } = await query('SELECT * FROM user WHERE id = $1', [userId]);
  if (rows.length !== 0) {
    return rows[0];
  }
  return null;
}

async function findByEmail({ email }) {
  const { rows } = await query('SELECT * FROM user WHERE email = $1', [email]);
  if (rows.length !== 0) {
    return rows[0];
  }
  return null;
}

async function getTags({ user }) {
  const { rows } = query('SELECT "tagName" FROM users_tags WHERE "userId" = $1', [user.id]);
  return rows.map(row => row.tagName);
}

async function setName({ ctx, name }) {
  if (!validator.isUnicodeLength(name, { max: 15 })) {
    throw new ParamsError('Max length of username is 15.');
  }
  const user = ctx.auth.signedInUser();
  if ((user.name || '') !== '') {
    throw new ParamsError('Name can only be set once.');
  }
  await query('UPDATE public.user SET name=$1 WHERE email=$2', [name, user.email]);
  return { ...user, name };
}

async function addSubbedTags({ ctx, tags }) {
  const user = ctx.auth.signedInUser();
  doTransaction(async (client) => {
    await Promise.all(tags.map(tag => client.query(
      'INSERT INTO users_tags ("userId", "tagName") VALUES ($1, $2) ON CONFLICT DO NOTHING',
      [user.id, tag],
    )));
  });
}

async function delSubbedTags({ ctx, tags }) {
  const user = ctx.auth.signedInUser();
  await query('DELETE FROM users_tags WHERE "userId" = $1 "tagName" IN $2', [user.id, tags]);
}

async function syncTags({ ctx, tags }) {
  const user = ctx.auth.signedInUser();
  const userTags = new Set(await this.getUserTags(ctx));
  const newTags = new Set(tags);
  const toDelete = [];
  const toInsert = [];
  tags.forEach((tag) => {
    if (!userTags.has(tag)) {
      toInsert.push(tag);
    }
  });
  userTags.forEach((tag) => {
    if (!newTags.has(tag)) {
      toDelete.push(tag);
    }
  });
  doTransaction(async (client) => {
    await Promise.all(toInsert.map(tag => client.query(
      'INSERT INTO users_tags ("userId", "tagName") VALUES ($1, $2) ON CONFLICT DO NOTHING',
      [user.id, tag],
    )));
    await query('DELETE FROM users_tags WHERE "userId" = $1 "tagName" IN $2', [user.id, tags]);
  });
}

async function banUser({ ctx, userId }) {
  const user = ctx.auth.signedInUser();
  ensurePermission({ user, action: ACTION.BAN_USER });
  await query('UPDATE public.user SET role=$1 WHERE id=$2', [ROLE.BANNED, userId]);
}


export default {
  findById,
  findByEmail,
  getTags,
  authContext,

  setName,
  addSubbedTags,
  delSubbedTags,
  syncTags,
  banUser,
};
export { ACTION };
