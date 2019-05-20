import JoiBase from '@hapi/joi';
import JoiObjectId from '~/utils/joiObjectId';
import mongo from '~/utils/mongo';
import Uid from '~/uid';
import { ParamsError, InternalError } from '~/utils/error';
import log from '~/utils/log';

import UserPostsModel from './userPosts';
import UserModel from './user';
import ThreadModel from './thread';
import NotificationModel from './notification';

const Joi = JoiBase.extend(JoiObjectId);
const THREAD = 'thread';
const POST = 'post';
const col = () => mongo.collection(POST);

const postSchema = Joi.object().keys({
  suid: Joi.string().alphanum().length(15).required(),
  threadSuid: Joi.string().alphanum().length(15).required(),
  anonymous: Joi.boolean().required(),
  author: Joi.string().required(),
  userId: Joi.objectId().required(),
  locked: Joi.boolean().default(false),
  blocked: Joi.boolean().default(false),
  createdAt: Joi.date().required(),
  updatedAt: Joi.date().required(),
  content: Joi.string().required(),
  quoteSuids: Joi.array().items(Joi.string().alphanum().length(15)).default([]),
});

// const PostSchema = new mongoose.Schema({
//   suid: { type: String, required: true },
//   userId: SchemaObjectId,
//   threadSuid: String,
//   anonymous: Boolean,
//   author: String,
//   createdAt: Date,
//   updatedAt: Date,
//   blocked: Boolean,
//   quoteSuids: [String],
//   content: String,
// }, { autoCreate: true });

const PostModel = ctx => ({
  pubPost: async function pubPost(input) {
    const {
      threadId: threadUid, anonymous, content, quoteIds: quoteUids = [],
    } = input;
    const { user } = ctx;

    const threadSuid = Uid.encode(threadUid);
    const threadDoc = await mongo.collection(THREAD).findOne({ suid: threadSuid });
    if (!threadDoc) {
      throw new ParamsError('Thread not found.');
    }
    if (threadDoc.locked) {
      throw new ParamsError('Thread is locked.');
    }

    const now = new Date();
    let post = {
      suid: await Uid.newSuid(),
      userId: user._id,
      threadSuid,
      anonymous,
      author: await UserModel(ctx).methods(user).author(threadSuid, anonymous),
      createdAt: now,
      updatedAt: now,
      blocked: false,
      quoteSuids: [],
      content,
    };

    let quotedPosts = [];
    if (quoteUids.length > 0) {
      quotedPosts = await col().find({
        suid: {
          $in: quoteUids.map(
            q => Uid.encode(q),
          ),
        },
      }).toArray();
      post.quoteSuids = quotedPosts.map(qp => qp.suid);
    }

    const { value, error } = postSchema.validate(post);
    if (error) {
      log.error(error);
      throw new ParamsError(`Thread validation failed, ${error}`);
    }
    post = value;

    const session = await mongo.startSession();
    session.startTransaction();

    try {
      await col().insertOne(post);
      await ThreadModel(ctx).methods(threadDoc).onPubPost(post, { session });
      await UserPostsModel(ctx).methods(user).onPubPost(threadDoc, post, { session });
      await NotificationModel(ctx).sendRepliedNoti(post, threadDoc, { session });
      await NotificationModel(ctx).sendQuotedNoti(post, threadDoc, quotedPosts, { session });
      await session.commitTransaction();
      session.endSession();
      return post;
    } catch (e) {
      await session.abortTransaction();
      session.endSession();
      throw new InternalError(`Transaction Failed: ${e}`);
    }
  },

  findByUid: async function findByUid(uid) {
    const post = await col().findOne({ suid: Uid.encode(uid) });
    return post;
  },

  methods: function methods(doc) {
    return genDoc(ctx, this, doc);
  },
});

const genDoc = (ctx, model, doc) => ({
  CACHED_UID: '',

  uid: function uid() {
    if (!this.CACHED_UID) this.CACHED_UID = Uid.decode(doc.suid);
    return this.CACHED_UID;
  },

  getQuotes: async function getQuotes() {
    let qs = [];
    if (doc.quoteSuids.length > 0) {
      qs = await col().find(
        { suid: { $in: doc.quoteSuids } },
      ).sort({ suid: 1 }).toArray();
    }
    return qs;
  },

  getContent: async function getContent() {
    return doc.blocked ? '' : doc.content;
  },

  quoteCount: async function quoteCount() {
    const count = await col().find({ quoteSuids: doc.suid }).count();
    return count;
  },
});

export default PostModel;
