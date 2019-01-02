import mongoose from 'mongoose';

import { encode, decode } from '~/utils/uuid';
import ThreadModel from './thread';
import NotificationModel from './notification';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

const PostSchema = new mongoose.Schema({
  userId: SchemaObjectId,
  threadId: SchemaObjectId,
  anonymous: Boolean,
  author: String,
  createdAt: Date,
  updatedAt: Date,
  blocked: Boolean,
  quotes: [SchemaObjectId],
  content: String,
}, { id: false });

PostSchema.statics.findById = async function findById(id) {
  const post = await PostModel.findOne({ _id: decode(id) }).exec();
  return post;
};
PostSchema.statics.pubPost = async function pubPost(ctx, input) {
  const user = { ctx };
  const { threadId, anonymous, content } = input;
  const now = new Date();
  const post = {
    _id: mongoose.Types.ObjectId(),
    userId: user._id,
    threadId,
    anonymous,
    author: await user.author(threadId, anonymous),
    createdAt: now,
    updatedAt: now,
    blocked: false,
    quoteIds: [],
    content,
  };
  const thread = await ThreadModel.findOne({ _id: encode(threadId) }).exec();
  let quotedPosts = [];
  if (input.quotes.length !== 0) {
    quotedPosts = await PostModel.find({
      _id: {
        $in: input.quotes.map(
          qid => encode(qid),
        ),
      },
    }).all().exec();
  }
  const session = await mongoose.startSession();
  await PostModel.create(post, { session });
  await thread.onPubPost(post, { session });
  await user.onPubPost(thread, post);
  await NotificationModel.sendRepliedNoti(post, thread, { session });
  await NotificationModel.sendQuotedNoti(post, thread, quotedPosts, { session });
  await session.commitTransaction();
  return post;
};

PostSchema.methods.id = function id() {
  return encode(this._id);
};
PostSchema.methods.getQuotes = async function quotes() {
  let qs = [];
  if (this.quotes.length !== 0) {
    qs = await PostModel.find({ _id: { $in: this.quoteIds } }).all().exec();
  }
  return qs;
};
PostSchema.methods.quoteCount = async function quoteCount() {
  const count = await PostModel.find({ quotes: this._id }).count().exec();
  return count;
};

const PostModel = mongoose.model('Post', PostSchema);

export default PostModel;
