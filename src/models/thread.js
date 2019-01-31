import mongoose from 'mongoose';

import Uid from '~/uid';
import findSlice from '~/models/base';
import TagModel from './tag';
import PostModel from './post';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

const ThreadSchema = mongoose.Schema({
  suid: String,
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
    postSuid: String,
    createdAt: Date,
  }],
}, { id: false });

ThreadSchema.statics.pubThread = async function pubThread(ctx, input) {
  const user = { ctx };
  const now = new Date();
  const thread = {
    ...input,
    userId: user._id,
    tags: [input.mainTag, ...(input.subTags || [])],
    locked: false,
    blocked: false,
    createdAt: now,
    updatedAt: now,
  };
  // TODO: validate main tag
  thread.suid = await Uid.newSuid();
  thread.author = await user.author(thread.suid);

  const session = await mongoose.startSession();
  await ThreadModel.create(thread, { session }).save();
  await TagModel.onPubThread(thread, { session });
  await session.commitTransaction();
  return thread;
};
ThreadSchema.statics.findById = async function findByUID(uid) {
  const thread = await PostModel.findOne({ suid: Uid.encode(uid) }).exec();
  return thread;
};
ThreadSchema.statics.getThreadSlice = async function getThreadSlice(
  tags = [], sliceQuery,
) {
  const option = {
    query: tags.length > 0 ? { tags: { $in: tags } } : {},
    desc: true,
    field: 'suid',
    sliceName: 'threads',
    parse: Uid.encode,
    toCursor: Uid.decode,
  };
  const result = await findSlice(sliceQuery, ThreadModel, option);
  return result;
};

ThreadSchema.methods.uid = function uid() {
  return Uid.decode(this.suid);
};
ThreadSchema.methods.getContent = function getContent() {
  return this.blocked ? '' : this.content;
};
ThreadSchema.methods.replies = async function replies(query) {
  const option = {
    query: { threadId: this.suid },
    field: '_id',
    sliceName: 'posts',
    parse: Uid.encode,
    toCursor: Uid.decode,
  };
  const result = await findSlice(query, PostModel, option);
  return result;
};
ThreadSchema.methods.replyCount = function replyCount() {
  return this.catalog.length;
};
ThreadSchema.methods.onPubPost = async function onPubPost(post, opt) {
  await ThreadModel.update({ suid: this.suid }, {
    $push: { catalog: { postId: post.suid, createdAt: post.createdAt } },
  }, opt);
};

const ThreadModel = mongoose.model('Thread', ThreadSchema);
export default ThreadModel;
