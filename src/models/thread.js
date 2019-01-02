import mongoose from 'mongoose';

import { encode, decode } from '~/utils/uuid';
import findSlice from '~/models/base';
import TagModel from './tag';
import PostModel from './post';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

const ThreadSchema = mongoose.Schema({
  anonymous: Boolean,
  author: String,
  userId: SchemaObjectId,
  mainTag: String,
  subTags: [String],
  tags: [String],
  title: String,
  locked: Boolean,
  blocked: Boolean,
  createdAt: Date,
  updatedAt: Date,
  content: String,
  catalog: [{
    postId: SchemaObjectId,
    createdAt: Date,
  }],
}, { id: false });

ThreadSchema.statics.pubThread = async function pubThread(ctx, input) {
  const user = { ctx };
  const now = new Date();
  const thread = {
    ...input,
    _id: mongoose.Types.ObjectId(),
    userId: user._id,
    tags: [input.mainTag, ...(input.subTags || [])],
    locked: false,
    blocked: false,
    createdAt: now,
    updatedAt: now,
  };

  if (input.anonymous) {
    thread.author = await user.anonymousId(thread._id);
  } else {
    if ((user.name || '') === '') {
      throw new Error('you must set name first');
    }
    thread.author = user.name;
  }
  const session = await mongoose.startSession();
  await ThreadModel(thread, { session }).save();
  await user.onPubThread(thread, { session });
  await TagModel.onPubThread(thread, { session });
  await session.commitTransaction();
  return thread;
};
ThreadSchema.statics.getThreadSlice = async function getThreadSlice(
  tags = [], sliceQuery,
) {
  const option = {
    query: tags.length > 0 ? { tags: { $in: tags } } : {},
    desc: true,
    field: '_id',
    sliceName: 'threads',
    parse: decode,
    toCursor: encode,
  };
  const result = await findSlice(sliceQuery, ThreadModel, option);
  return result;
};

ThreadSchema.methods.id = function id() {
  return encode(this._id);
};
ThreadSchema.methods.replies = async function replies(query) {
  const option = {
    query: { threadId: this._id },
    field: '_id',
    sliceName: 'posts',
    parse: decode,
    toCursor: encode,
  };
  const result = await findSlice(query, PostModel, option);
  return result;
};
ThreadSchema.methods.replyCount = function replyCount() {
  return this.catalog.length;
};
ThreadSchema.methods.onPubPost = async function onPubPost(post, opt) {
  await ThreadModel.update({ _id: this._id }, {
    $push: { catalog: { postId: post._id, createdAt: post.createdAt } },
  }, opt);
};

const ThreadModel = mongoose.model('Thread', ThreadSchema);
export default ThreadModel;
