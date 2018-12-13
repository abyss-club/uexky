import mongoose from 'mongoose';
import { encode } from './uuid';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

// MODEL: User
//        storage user info.
const UserSchema = new mongoose.Schema({
  email: { type: String, unique: true },
  name: { type: String, unique: true },
  tags: [String],
  read_noti_time: {
    system: Date,
    replied: Date,
    quoted: Date,
  },
  role: {
    type: String,
    range: [String],
  },
});
const UserModel = mongoose.model('User', UserSchema);
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
UserSchema.methods.posts = async function posts() {
  // TODO: use slice query
  const userPosts = await UserPostsModel.find({ userId: this._id })
    .sort({ updatedAt: -1 }).limit(10).exec();
  return userPosts;
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
  updatedAt: Date(),
});
const UserPostsModel = mongoose.model('UserPosts', UserPostsSchema);

export default UserModel;
export { UserAIDModel };
