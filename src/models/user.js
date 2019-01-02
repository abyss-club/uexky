import { AuthenticationError } from 'apollo-server-koa';
import mongoose from 'mongoose';

import AuthFail from '~/error';
import { encode } from '../utils/uuid';

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
    type: String,
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

const UserModel = mongoose.model('User', UserSchema);
export default UserModel;
export { UserAIDModel, getUserByEmail, ensureSignIn };
