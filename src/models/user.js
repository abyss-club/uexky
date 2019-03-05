import mongoose from 'mongoose';

import Uid from '~/uid';
import validator from '~/utils/validator';
import findSlice from '~/models/base';
import PostModel from '~/models/post';
import ThreadModel from '~/models/thread';
import { ParamsError, AuthError, InternalError } from '~/utils/error';

const SchemaObjectId = mongoose.ObjectId;

// MODEL: User
//        storage user info.
const UserSchema = new mongoose.Schema({
  email: { type: String, required: true, unique: true },
  name: {
    type: String,
    index: {
      unique: true,
      partialFilterExpression: { name: { $type: 'string' } },
    },
  },
  tags: [String],
  readNotiTime: {
    system: Date,
    replied: Date,
    quoted: Date,
  },
  role: {
    role: String,
    range: [String],
  },
});
UserSchema.methods.author = async function author(threadSuid, anonymous) {
  if (anonymous) {
    const aid = await UserAidModel.getAid(this._id, threadSuid);
    return aid;
  }
  if ((this.name || '') === '') {
    throw new Error('Name not yet set.');
  }
  return this.name;
};
UserSchema.methods.posts = async function posts(sliceQuery) {
  const option = {
    query: { userId: this._id },
    desc: true,
    field: 'updatedAt',
    sliceName: 'threads',
    parse: value => new Date(value),
    toCursor: value => value.toISOString(),
  };
  const userPosts = await findSlice(sliceQuery, UserPostsModel, option);
  return userPosts;
};
UserSchema.methods.getRole = async function getRole() {
  const { role, tags } = this.role || {};
  return { role: role || 0, tags: tags || [] };
};
UserSchema.methods.onPubPost = async function onPubPost(
  thread, post, { session },
) {
  await UserPostsModel.updateOne({
    userId: this._id,
    threadSuid: thread.suid,
  }, {
    $setOnInsert: {
      userId: this._id,
      threadSuid: thread.suid,
    },
    $set: { updatedAt: Date() },
    $push: {
      posts: {
        suid: post.suid,
        createdAt: post.createdAt,
        anonymous: post.anonymous,
      },
    },
  }, { session, upsert: true });
};
UserSchema.methods.onPubThread = async function onPubThread(thread, { session }) {
  await UserPostsModel.updateOne({
    userId: this._id,
    threadSuid: thread.suid,
  }, {
    $setOnInsert: {
      userId: this._id,
      threadSuid: thread.suid,
    },
    $set: { updatedAt: Date() },
  }, { session, upsert: true });
};
UserSchema.methods.setName = async function setName(name) {
  if (!validator.isUnicodeLength(name, { max: 15 })) {
    throw new ParamsError('Max length of username is 15.');
  }
  if ((this.name || '') !== '') {
    throw new InternalError('Name can only be set once.');
  }
  await UserModel.updateOne({ _id: this._id }, { $set: { name } }).exec();
  const result = await UserModel.findOne({ _id: this._id }).exec();
  return result;
};
UserSchema.methods.syncTags = async function syncTags(tags) {
  await UserModel.updateOne({ _id: this._id }, { tags: tags || [] }).exec();
  const result = await UserModel.findOne({ _id: this._id }).exec();
  return result;
};
UserSchema.methods.addSubbedTags = async function addSubbedTags(tags) {
  if (!Array.isArray(tags) || !tags.length) {
    throw new ParamsError('Provided tags is not a non-empty array.');
  }
  await UserModel.updateOne(
    { _id: this._id },
    { $addToSet: { tags: { $each: tags } } },
  ).exec();
  const result = await UserModel.findOne({ _id: this._id }).exec();
  return result;
};
UserSchema.methods.delSubbedTags = async function delSubbedTags(tags) {
  if (!Array.isArray(tags) || !tags.length) {
    throw new ParamsError('Provided tags is not a non-empty array.');
  }
  await UserModel.updateOne(
    { _id: this._id },
    { $pull: { tags: { $in: tags } } },
  ).exec();
  const result = await UserModel.findOne({ _id: this._id }).exec();
  return result;
};
UserSchema.methods.ensurePermission = function ensurePermission(
  targetUser, action, tag,
) {
  if (this.getRole().role <= targetUser.getRole().role) {
    throw new Error('Params Error: unknown action');
  }
  const { role, tags } = this.getRole();
  if (ACTIONS[action]) {
    throw new Error('Params Error: unknown action');
  }
  if (role === ROLES.TagAdmin) {
    if (tags.findIndex(t => t === tag) === -1) {
      throw new Error('Premitted Error');
    }
  }
  if (ROLES_ACTIONS[role].findIndex(a => a === action) === -1) {
    throw new Error('Premitted Error');
  }
};
UserSchema.methods.banUser = async function banUser(postUid) {
  const post = await PostModel.findByUid(postUid);
  const thread = await ThreadModel.findOne({ suid: post.threadSuid }).exec();
  const target = await UserModel.findOne({ _id: post.userId }).exec();
  this.ensurePermission(target, ACTIONS.BAN_USER, thread.mainTag);
  await UserModel.updateOne(
    { _id: post.userId },
    { $set: { 'role.role': ROLES.Banned } },
  ).exec();
};
UserSchema.methods.blockPost = async function blockPost(postUid) {
  const post = await PostModel.findByUid(postUid);
  const thread = await ThreadModel.findOne({ suid: post.threadSuid }).exec();
  const target = await UserModel.findOne({ _id: post.userId }).exec();
  this.ensurePermission(target, ACTIONS.BLOCK_POST, thread.mainTag);
  await PostModel.updateOne({ _id: post._id }, { $set: { blocked: true } }).exec();
};
UserSchema.methods.lockThread = async function lockThread(threadUid) {
  const thread = await ThreadModel.findByUid(threadUid);
  const target = await UserModel.findOne({ _id: thread.userId }).exec();
  this.ensurePermission(target, ACTIONS.LOCK_THREAD, thread.mainTag);
  await ThreadModel.updateOne({ _id: thread._id }, { $set: { locked: true } }).exec();
};
UserSchema.methods.blockThread = async function blockThread(threadUid) {
  const thread = await ThreadModel.findByUid(threadUid);
  const target = await UserModel.findOne({ _id: thread.userId }).exec();
  this.ensurePermission(target, ACTIONS.BLOCK_THREAD, thread.mainTag);
  await ThreadModel.updateOne({ _id: thread._id }, { $set: { blocked: true } }).exec();
};
UserSchema.methods.editTags = async function editTags(
  threadUid, mainTag, subTags,
) {
  const thread = await ThreadModel.findByUid(threadUid);
  const target = await UserModel.findOne({ _id: thread.userId }).exec();
  this.ensurePermission(target, ACTIONS.BLOCK_THREAD, thread.mainTag);
  await ThreadModel.updateOne(
    { _id: thread._id },
    { $set: { mainTag, subTags } },
  ).exec();
};
UserSchema.methods.setReadNotiTime = async function setReadNotiTime(type, time) {
  const notiType = `readNotiTime.${type}`;
  await UserModel.updateOne({ _id: this._id }, { $set: { [notiType]: time } });
};
// MODEL: UserAid
//        used for save anonymousId for user in threads.
const UserAidSchema = new mongoose.Schema({
  userId: SchemaObjectId,
  threadSuid: String,
  anonymousId: String, // format: Uid
});

UserAidSchema.statics.getAid = async function getAid(userId, threadSuid) {
  const result = await UserAidModel.findOneAndUpdate({
    userId, threadSuid,
  }, {
    $setOnInsert: {
      userId,
      threadSuid,
      anonymousId: Uid.decode(await Uid.newSuid()),
    },
    $set: { updatedAt: Date() },
  }, { new: true, upsert: true }).exec();
  return result.anonymousId;
};

const UserAidModel = mongoose.model('UserAid', UserAidSchema);

// MODEL: UserPosts
//        used for querying user's posts grouped by thread.
const UserPostsSchema = new mongoose.Schema({
  userId: SchemaObjectId,
  threadSuid: String,
  updatedAt: Date,
  posts: [{
    suid: String,
    createdAt: Date,
    anonymous: Boolean,
  }],
}, { autoCreate: true });
const UserPostsModel = mongoose.model('UserPosts', UserPostsSchema);

UserSchema.statics.getUserByEmail = async function getUserByEmail(email) {
  try {
    const user = await this.findOne({ email }).exec();
    if (user) return user;
    // const newUser = new UserModel({ email });
    const res = await this.create({ email });
    return res;
  } catch (e) {
    throw new AuthError(e);
  }
};

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

const UserModel = mongoose.model('User', UserSchema);
export default UserModel;
export {
  UserAidModel, UserPostsModel, ensureSignIn,
};
