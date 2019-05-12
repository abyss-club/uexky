import Joi from '@hapi/joi';
import mongo, { db } from '~/utils/mongo';
import log from '~/utils/log';

import validator from '~/utils/validator';
import findSlice from '~/models/base';
// import PostModel from '~/models/post';
import ThreadModel from '~/models/thread';
import UserAidModel from '~/models/userAid';
// import UserPostsModel from '~/models/userPosts';
import {
  ParamsError, AuthError, InternalError, PermissionError,
} from '~/utils/error';

const USER = 'user';
const USERPOSTS = 'userPosts';
const THREAD = 'thread';
const col = () => mongo.collection(USER);

const tagsSchema = Joi.array().items(Joi.string());

const userSchema = Joi.object().keys({
  email: Joi.string().email().required(),
  name: Joi.string(),
  tags: tagsSchema,
  readNotiTime: Joi.object().keys({
    system: Joi.date(),
    replied: Joi.date(),
    quoted: Joi.date(),
  }),
  role: Joi.object().keys({
    role: Joi.string(),
    range: Joi.array().items(Joi.string()),
  }),
});

const UserModel = ctx => ({
  getUserByEmail: async function getUserByEmail(email) {
    try {
      const user = await col().findOne({ email });
      if (user) return user;
      const { error } = userSchema.validate({ email });
      if (error) {
        log.error(error);
        throw new ParamsError(`Email validation failed, ${error}`);
      }
      const res = await col().insertOne({ email });
      return res.ops[0];
    } catch (e) {
      throw new AuthError(e);
    }
  },
  methods: function methods(doc) {
    return genUserDoc(ctx, this, doc);
  },
});

const genUserDoc = (ctx, model, doc) => ({
  setName: async function setName(name) {
    if (!validator.isUnicodeLength(name, { max: 15 })) {
      throw new ParamsError('Max length of username is 15.');
    }
    if ((doc.name || '') !== '') {
      throw new InternalError('Name can only be set once.');
    }
    await col().updateOne(
      { _id: doc._id }, { $set: { name } },
    );
    return { ...doc, name };
  },

  syncTags: async function syncTags(tags) {
    const { value, error } = tagsSchema.validate(tags);
    if (error) {
      log.error(error);
      throw new ParamsError(`Tags validation failed, ${error}`);
    }
    await col().updateOne(
      { _id: doc._id }, { $set: { tags: value || [] } },
    );
    return { ...doc, tags };
  },

  addSubbedTags: async function addSubbedTags(tags) {
    const { value, error } = tagsSchema.validate(tags);
    if (error) {
      log.error(error);
      throw new ParamsError(`Tags validation failed, ${error}`);
    }
    await col().updateOne(
      { _id: doc._id },
      { $addToSet: { tags: { $each: value } } },
    );
    const newTags = [...new Set([...doc.tags, ...tags])];
    return { ...doc, tags: newTags };
  },

  delSubbedTags: async function delSubbedTags(tags) {
    const { value, error } = tagsSchema.validate(tags);
    if (error) {
      log.error(error);
      throw new ParamsError(`Tags validation failed, ${error}`);
    }
    // console.log({ doc });
    await col().updateOne(
      { _id: doc._id },
      { $pullAll: { tags: value } },
    );
    const newTags = doc.tags.filter(tag => !tags.includes(tag));
    return { ...doc, tags: newTags };
  },

  author: async function author(threadSuid, anonymous) {
    if (anonymous) {
      // console.log({ doc });
      const aid = await UserAidModel(ctx).getAid(doc._id, threadSuid);
      return aid;
    }
    if ((doc.name || '') === '') {
      throw new InternalError('Name not yet set.');
    }
    return doc.name;
  },

  setReadNotiTime: async function setReadNotiTime(type, time) {
    const notiType = `readNotiTime.${type}`;
    await col().updateOne({ _id: doc._id }, { $set: { [notiType]: time } });
  },

  getRole: async function getRole() {
    const { role, tags } = ctx.role || {};
    return { role: role || 0, tags: tags || [] };
  },

  ensurePermission: function ensurePermission(
    targetUser, action, tag,
  ) {
    if (this.getRole().role <= targetUser.getRole().role) {
      throw new ParamsError('Params Error: unknown action');
    }
    const { role, tags } = this.getRole();
    if (ACTIONS[action]) {
      throw new ParamsError('Params Error: unknown action');
    }
    if (role === ROLES.TagAdmin) {
      if (tags.findIndex(t => t === tag) === -1) {
        throw new PermissionError('Premitted Error');
      }
    }
    if (ROLES_ACTIONS[role].findIndex(a => a === action) === -1) {
      throw new PermissionError('Premitted Error');
    }
  },

  editTags: async function editTags(
    threadUid, mainTag, subTags,
  ) {
    const thread = await ThreadModel(ctx).findByUid(threadUid);
    const target = await col().findOne({ _id: thread.userId });
    this.ensurePermission(target, ACTIONS.BLOCK_THREAD, thread.mainTag);
    await db.collection(THREAD).updateOne(
      { _id: thread._id },
      { $set: { mainTag, subTags } },
    );
  },

  posts: async function posts(sliceQuery) {
    const option = {
      query: { userId: doc._id },
      desc: true,
      field: 'updatedAt',
      sliceName: 'threads',
      parse: value => new Date(value),
      toCursor: value => value.toISOString(),
    };
    const userPosts = await findSlice(sliceQuery, db.collection(USERPOSTS), option);
    return userPosts;
  },
});

// MODEL: User
//        storage user info.
// const UserSchema = new mongoose.Schema({
//   email: { type: String, required: true, unique: true },
//   name: {
//     type: String,
//     index: {
//       unique: true,
//       partialFilterExpression: { name: { $type: 'string' } },
//     },
//   },
//   tags: [String],
//   readNotiTime: {
//     system: Date,
//     replied: Date,
//     quoted: Date,
//   },
//   role: {
//     role: String,
//     range: [String],
//   },
// });


//
// const banUser = async function banUser(postUid) {
//   const post = await PostModel.findByUid(postUid);
//   const thread = await ThreadModel.findOne({ suid: post.threadSuid }).exec();
//   const target = await UserModel.findOne({ _id: post.userId }).exec();
//   this.ensurePermission(target, ACTIONS.BAN_USER, thread.mainTag);
//   await UserModel.updateOne(
//     { _id: post.userId },
//     { $set: { 'role.role': ROLES.Banned } },
//   ).exec();
// };
//
// const blockPost = async function blockPost(postUid) {
//   const post = await PostModel.findByUid(postUid);
//   const thread = await ThreadModel.findOne({ suid: post.threadSuid }).exec();
//   const target = await UserModel.findOne({ _id: post.userId }).exec();
//   this.ensurePermission(target, ACTIONS.BLOCK_POST, thread.mainTag);
//   await PostModel.updateOne({ _id: post._id }, { $set: { blocked: true } }).exec();
// };
//
// const lockThread = async function lockThread(threadUid) {
//   const thread = await ThreadModel.findByUid(threadUid);
//   const target = await UserModel.findOne({ _id: thread.userId }).exec();
//   this.ensurePermission(target, ACTIONS.LOCK_THREAD, thread.mainTag);
//   await ThreadModel.updateOne({ _id: thread._id }, { $set: { locked: true } }).exec();
// };
//
// const blockThread = async function blockThread(threadUid) {
//   const thread = await ThreadModel.findByUid(threadUid);
//   const target = await UserModel.findOne({ _id: thread.userId }).exec();
//   this.ensurePermission(target, ACTIONS.BLOCK_THREAD, thread.mainTag);
//   await ThreadModel.updateOne({ _id: thread._id }, { $set: { blocked: true } }).exec();
// };

function ensureSignIn(ctx) {
  if (!ctx.user) throw new AuthError('Authentication needed.');
  return ctx.user;
}

const ACTIONS = {
  BAN_USER: 'BAN_USER',
  BLOCK_POST: 'BLOCK_POST',
  LOCK_THREAD: 'LOCK_THREAD',
  BLOCK_THREAD: 'BLOCK_THREAD',
  EDIT_TAG: 'EDIT_TAG',
  PUB_POST: 'PUB_POST',
  PUB_THREAD: 'PUB_THREAD',
};

const ROLES = {
  SuperAdmin: 1000,
  TagAdmin: 100,
  Normal: 0,
  Banned: -1,
};

const ROLES_ACTIONS = {
  SuperAdmin: ACTIONS,
  TagAdmin: ACTIONS,
  Normal: [ACTIONS.PUB_POST, ACTIONS.PUB_THREAD],
  Banned: [],
};

export default UserModel;

export {
  // getUserByEmail,
  // syncTags,
  ensureSignIn,
  // ACTIONS,
  // ROLES,
  // ROLES_ACTIONS,
};
