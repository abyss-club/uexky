import mongoose from 'mongoose';

import { ParamsError } from '~/utils/error';
import Uid from '~/uid';
import validator from '~/utils/validator';
import findSlice from '~/models/base';
import TagModel from './tag';
import PostModel from './post';

const SchemaObjectId = mongoose.ObjectId;

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
}, { autoCreate: true });

ThreadSchema.statics.pubThread = async function pubThread({ user }, input) {
  const {
    anonymous, content, title, mainTag, subTags = [],
  } = input;

  const mainTags = await TagModel.getMainTags();
  if (!mainTags.includes(mainTag)) throw new ParamsError('Invalid mainTag');
  if (!validator.isUnicodeLength(title, { max: 28 })) {
    throw new ParamsError('Max length of title is 28.');
  }

  const now = new Date();
  const thread = {
    anonymous,
    userId: user._id,
    tags: [mainTag, ...(subTags)],
    mainTag,
    subTags,
    locked: false,
    blocked: false,
    createdAt: now,
    updatedAt: now,
    content,
    title,
  };
  thread.suid = await Uid.newSuid();
  thread.author = await user.author(thread.suid, anonymous);

  const session = await mongoose.startSession();
  session.startTransaction();
  const threadDoc = new ThreadModel(thread);
  await threadDoc.save({ session });
  await ThreadModel.create(threadDoc, { session });
  await user.onPubThread(threadDoc, { session });
  await TagModel.onPubThread(threadDoc, { session });
  await session.commitTransaction();
  session.endSession();

  return threadDoc;
};
ThreadSchema.statics.findByUid = async function findByUid(uid) {
  const thread = await ThreadModel.findOne({ suid: Uid.encode(uid) }).exec();
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
  await ThreadModel.updateOne({ suid: this.suid }, {
    $push: { catalog: { postId: post.suid, createdAt: post.createdAt } },
  }, opt).exec();
};

const ThreadModel = mongoose.model('Thread', ThreadSchema);
export default ThreadModel;
