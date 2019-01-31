import mongoose from 'mongoose';

import Uid from '~/uid';
import ThreadModel from './thread';
import NotificationModel from './notification';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

const PostSchema = new mongoose.Schema({
  suid: String, // TODO: unique
  userId: SchemaObjectId,
  threadSuid: String,
  anonymous: Boolean,
  author: String,
  createdAt: Date,
  updatedAt: Date,
  blocked: Boolean,
  quoteSuids: [String],
  content: String,
}, { id: false });

PostSchema.statics.findById = async function findByUid(uid) {
  const post = await PostModel.findOne({ suid: Uid.encode(uid) }).exec();
  return post;
};
PostSchema.statics.pubPost = async function pubPost(ctx, input) {
  const user = { ctx };
  const { threadId: threadUid, anonymous, content } = input;
  const threadSuid = Uid.encode(threadUid);
  const now = new Date();
  const post = {
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
  const thread = await ThreadModel.findOne({ suid: threadSuid }).exec();
  if (thread.locked) {
    throw new Error('this thread is locked');
  }

  let quotedPosts = [];
  if (input.quotes.length !== 0) {
    quotedPosts = await PostModel.find({
      suid: {
        $in: input.quotes.map(
          q => Uid.encode(q),
        ),
      },
    }).all().exec();
    post.quoteIds = quotedPosts.map(qp => qp.suid);
  }

  post.suid = await Uid.newSuid();
  const session = await mongoose.startSession();
  await PostModel.create(post, { session });
  await thread.onPubPost(post, { session });
  await user.onPubPost(thread, post, { session });
  await NotificationModel.sendRepliedNoti(post, thread, { session });
  await NotificationModel.sendQuotedNoti(post, thread, quotedPosts, { session });
  await session.commitTransaction();
  return post;
};

PostSchema.methods.uid = function uid() {
  return Uid.decode(this.suid);
};
PostSchema.methods.getQuotes = async function getQuotes() {
  let qs = [];
  if (this.quotes.length !== 0) {
    qs = await PostModel.find(
      { suid: { $in: this.quoteSuids } },
      { sort: { suid: 1 } },
    ).all().exec();
  }
  return qs;
};
PostSchema.methods.getContent = async function getContent() {
  return this.blocked ? '' : this.content;
};
PostSchema.methods.quoteCount = async function quoteCount() {
  const count = await PostModel.find({ quotes: this.suid }).count().exec();
  return count;
};

const PostModel = mongoose.model('Post', PostSchema);

export default PostModel;
