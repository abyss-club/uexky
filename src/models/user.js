import { AuthenticationError } from 'apollo-server-koa';
import mongoose from 'mongoose';

import AuthFail from '~/error';
import { encode } from '~/utils/uuid';
import PostModel from '~/models/post';
import ThreadModel from '~/models/thread';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

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
UserSchema.methods.author = async function author(threadId, anonymous) {
  if (anonymous) {
    const obj = { userId: this.ObjectId, threadId };
    await UserAIDModel.update(obj, obj, { upsert: true });
    const aid = await UserAIDModel.findOne(obj);
    return aid.anonymousId;
  }
  if ((this.name || '') === '') {
    throw new Error('you must set name first');
  }
  return this.name;
};
UserSchema.methods.posts = async function posts() {
  // TODO: use slice query
  const userPosts = await UserPostsModel.find({ userId: this._id })
    .sort({ updatedAt: -1 }).limit(10).exec();
  return userPosts;
};
UserSchema.methods.getRole = async function getRole() {
  const { role, tags } = this.role || {};
  return { role: role || 0, tags: tags || [] };
};
UserSchema.methods.onPubPost = async function onPubPost(thread, post) {
  await UserPostsModel.update({
    userId: this._id,
    threadId: thread._id,
  }, {
    $push: { posts: post._id },
    $set: { updatedAt: Date() },
  });
};
UserSchema.methods.setName = async function setName(name) {
  // TODO: validate length
  if ((this.name || '') !== '') {
    throw new Error('already set name!');
  }
  await UserModel.update({ _id: this._id }, { $set: { name } }).exec();
};
UserSchema.methods.syncTags = async function syncTags(tags) {
  await UserModel.update({ _id: this._id }, { tags: tags || [] });
};
UserSchema.methods.addSubbedTags = async function addSubbedTags(tags) {
  if ((tags || []).length === 0) {
    return;
  }
  await UserModel.update(
    { _id: this._id },
    { $addToSet: { tags: { $each: tags } } },
  );
};
UserSchema.methods.delSubbedTags = async function delSubbedTags(tags) {
  if ((tags || []).length === 0) {
    return;
  }
  await UserModel.update(
    { _id: this._id },
    { $pull: { tags: { $in: tags } } },
  );
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
UserSchema.methods.banUser = async function banUser(postId) {
  const post = await PostModel.findById(postId);
  const thread = await ThreadModel.findById(post.threadId);
  const target = await UserModel.findOne({ _id: post.userId });
  this.ensurePermission(target, ACTIONS.BAN_USER, thread.mainTag);
  await UserModel.update(
    { _id: post.userId },
    { $set: { 'role.role': ROLES.Banned } },
  );
};
UserSchema.methods.blockPost = async function blockPost(postId) {
  const post = await PostModel.findById(postId);
  const thread = await ThreadModel.findById(post.threadId);
  const target = await UserModel.findOne({ _id: post.userId });
  this.ensurePermission(target, ACTIONS.BLOCK_POST, thread.mainTag);
  await PostModel.update({ _id: post._id }, { $set: { blocked: true } });
};
UserSchema.methods.lockThread = async function lockThread(threadId) {
  const thread = await ThreadModel.findById(threadId);
  const target = await UserModel.findOne({ _id: thread.userId });
  this.ensurePermission(target, ACTIONS.LOCK_THREAD, thread.mainTag);
  await ThreadModel.update({ _id: thread._id }, { $set: { locked: true } });
};
UserSchema.methods.blockThread = async function blockThread(threadId) {
  const thread = await ThreadModel.findById(threadId);
  const target = await UserModel.findOne({ _id: thread.userId });
  this.ensurePermission(target, ACTIONS.BLOCK_THREAD, thread.mainTag);
  await ThreadModel.update({ _id: thread._id }, { $set: { blocked: true } });
};
UserSchema.methods.editTags = async function editTags(
  threadId, mainTag, subTags,
) {
  const thread = await ThreadModel.findById(threadId);
  const target = await UserModel.findOne({ _id: thread.userId });
  this.ensurePermission(target, ACTIONS.BLOCK_THREAD, thread.mainTag);
  await ThreadModel.update(
    { _id: thread._id },
    { $set: { mainTag, subTags } },
  );
};

// MODEL: UserAID
//        used for save anonymousId for user in threads.
const UserAIDSchema = new mongoose.Schema({
  userId: SchemaObjectId,
  threadId: SchemaObjectId,
});
const UserAIDModel = mongoose.model('UserAID', UserAIDSchema);

UserAIDSchema.methods.anonymousId = function anonymousId() {
  return encode(this.ObjectId);
};

// MODEL: UserPosts
//        used for querying user's posts grouped by thread.
const UserPostsSchema = mongoose.Schema({
  userId: SchemaObjectId,
  threadId: SchemaObjectId,
  posts: [SchemaObjectId],
  updatedAt: Date,
});
const UserPostsModel = mongoose.model('UserPosts', UserPostsSchema);

async function getUserByEmail(email) {
  try {
    const user = await UserModel.findOne({ email });
    if (user) return user;
    // const newUser = new UserModel({ email });
    const res = await UserModel.create({ email });
    return res;
  } catch (e) {
    throw new AuthFail(e);
  }
}

function ensureSignIn(ctx) {
  if (!ctx.user) throw new AuthenticationError('Authentication needed.');
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
  UserAIDModel, getUserByEmail, ensureSignIn,
};
