import mongoose from 'mongoose';

import Uid from '~/uid';
import { ParamsError } from '~/utils/error';
import ThreadModel from './thread';
import NotificationModel from './notification';

const SchemaObjectId = mongoose.ObjectId;

const PostSchema = new mongoose.Schema({
  suid: { type: String, required: true, unique: true },
  userId: SchemaObjectId,
  threadSuid: String,
  anonymous: Boolean,
  author: String,
  createdAt: Date,
  updatedAt: Date,
  blocked: Boolean,
  quoteSuids: [String],
  content: String,
}, { autoCreate: true });

PostSchema.statics.pubPost = async function pubPost({ user }, input) {
  const {
    threadId: threadUid, anonymous, content, quoteIds: quoteUids = [],
  } = input;

  const threadSuid = Uid.encode(threadUid);
  const threadDoc = await ThreadModel.findOne({ suid: threadSuid });
  if (!threadDoc) {
    throw new ParamsError('Thread not found.');
  }
  if (threadDoc.locked) {
    throw new ParamsError('Thread is locked.');
  }

  const now = new Date();
  const post = {
    suid: await Uid.newSuid(),
    userId: user._id,
    threadSuid,
    anonymous,
    author: await user.author(threadSuid, anonymous),
    createdAt: now,
    updatedAt: now,
    blocked: false,
    quoteSuids: [],
    content,
  };

  let quotedPosts = [];
  if (quoteUids.length > 0) {
    quotedPosts = await PostModel.find({
      suid: {
        $in: quoteUids.map(
          q => Uid.encode(q),
        ),
      },
    });
    post.quoteSuids = quotedPosts.map(qp => qp.suid);
  }

  const session = await mongoose.startSession();
  session.startTransaction();
  const postDoc = new PostModel(post);
  await postDoc.save({ session });
  await threadDoc.onPubPost(postDoc, { session });
  await user.onPubPost(threadDoc, postDoc, { session });
  await NotificationModel.sendRepliedNoti(postDoc, threadDoc, { session });
  await NotificationModel.sendQuotedNoti(postDoc, threadDoc, quotedPosts, { session });
  await session.commitTransaction();
  session.endSession();

  return postDoc;
};

PostSchema.statics.findByUid = async function findByUid(uid) {
  const post = await PostModel.findOne({ suid: Uid.encode(uid) }).exec();
  return post;
};
PostSchema.methods.uid = function uid() {
  if (!this.CACHED_UID) this.CACHED_UID = Uid.decode(this.suid);
  return this.CACHED_UID;
};
PostSchema.methods.getQuotes = async function getQuotes() {
  let qs = [];
  if (this.quoteSuids.length > 0) {
    qs = await PostModel.find(
      { suid: { $in: this.quoteSuids } },
    ).sort({ suid: 1 }).exec();
  }
  return qs;
};
PostSchema.methods.getContent = async function getContent() {
  return this.blocked ? '' : this.content;
};
PostSchema.methods.quoteCount = async function quoteCount() {
  const count = await PostModel.find({ quoteSuids: this.suid }).countDocuments().exec();
  return count;
};

const PostModel = mongoose.model('Post', PostSchema);

export default PostModel;
