import { NotFoundError } from '~/utils/error';
import { query } from '~/utils/pg';
import { ACTION } from '~/models/user';

// pgm.createTable('thread', {
//   id: { type: 'bigint', primaryKey: true },
//   createdAt: { type: 'timestamp', notNull: true },
//   updatedAt: { type: 'timestamp', notNull: true },
//
//   anonymous: { type: 'boolean', notNull: true },
//   userId: { type: 'integer', notNull: true, references: 'user(id)' },
//   userName: { type: 'varchar(16)', references: 'user(name)' },
//   anonymousId: { type: 'bigint', references: 'anonymous_id(anonymous_id)' },
//
//   title: { type: 'text', notNull: true },
//   locked: { type: 'bool', notNull: true },
//   blocked: { type: 'bool', notNull: true },
//   content: { type: 'text', notNull: true },
// });
// pgm.createIndex('title', 'anonymous', 'userId', 'blocked');

async function findById({ threadId }) {
  // TODO: care uid
  const { rows } = await query('SELECT * FROM thread WHERE id=$1', [threadId]);
  if ((rows || []).length === 0) {
    throw NotFoundError(`cant find thread ${threadId}`);
  }
  return rows[0];
}

async function findUserThreads({ user, query }) {
  // TODO: slice query
  // TODO: need by resolvers/user
}

async function lockThread({ ctx, threadId }) {
  // TODO: care uid
  ctx.auth.ensurePermission(ACTION.LOCK_THREAD);
  await query('UPDATE thread SET lock=$1 WHERE id=$2', [true, threadId]);
}

async function blockThread({ ctx, threadId }) {
  // TODO: care uid
  ctx.auth.ensurePermission(ACTION.BLOCK_THREAD);
  await query('UPDATE thread SET lock=$1 WHERE id=$2', [true, threadId]);
}

async function editTags({
  ctx, threadId, mainTag, subTags,
}) {
  // TODO: need by resolvers/user
}

export default {
  findById,
  findUserThreads,
  lockThread,
  blockThread,
  editTags,
};

/*
import JoiBase from '@hapi/joi';
import JoiObjectId from '~/utils/joiObjectId';
import { query } from '~/utils/pg';
import { ParamsError, InternalError } from '~/utils/error';
import Uid from '~/uid';
import validator from '~/utils/validator';
import findSlice from '~/models/base';
import log from '~/utils/log';

import UserPostsModel from './userPosts';
import UserModel from './user';
import TagModel from './tag';

const Joi = JoiBase.extend(JoiObjectId);
const THREAD = 'thread';
const POST = 'post';
const col = () => mongo.collection(THREAD);

const threadSchema = Joi.object().keys({
  suid: Joi.string().alphanum().length(15).required(),
  anonymous: Joi.boolean().required(),
  author: Joi.string().required(),
  userId: Joi.objectId().required(),
  mainTag: Joi.string().required(),
  subTags: Joi.array().items(Joi.string()).required(),
  tags: Joi.array().items(Joi.string()).required(),
  title: Joi.string().required(),
  locked: Joi.boolean().default(false),
  blocked: Joi.boolean().default(false),
  createdAt: Joi.date().required(),
  updatedAt: Joi.date().required(),
  content: Joi.string().required(),
  catalog: Joi.array().items(Joi.object().keys({
    postSuid: Joi.string().alphanum().length(15).required(),
    createdAt: Joi.date().required(),
  })),
});

const ThreadModel = ctx => ({
  pubThread: async function pubThread(input) {
    const {
      anonymous, content, title, mainTag, subTags = [],
    } = input;
    const { user } = ctx;

    const mainTags = await TagModel().getMainTags();
    if (!mainTags.includes(mainTag)) throw new ParamsError('Invalid mainTag');
    if (!validator.isUnicodeLength(title, { max: 28 })) {
      throw new ParamsError('Max length of title is 28.');
    }

    const now = new Date();
    let thread = {
      suid: await Uid.newSuid(),
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

    // console.log({ user });
    thread.author = await UserModel(ctx).methods(user).author(thread.suid, anonymous);

    const { value, error } = threadSchema.validate(thread);
    if (error) {
      log.error(error);
      throw new ParamsError(`Thread validation failed, ${error}`);
    }
    thread = value;

    const session = await mongo.startSession();
    session.startTransaction();
    try {
      await col().insertOne(thread);
      await UserPostsModel(ctx).methods(user).onPubThread(thread, { session });
      await TagModel().onPubThread(thread, { session });
      await session.commitTransaction();
      session.endSession();

      return thread;
    } catch (e) {
      await session.abortTransaction();
      session.endSession();
      throw new InternalError(`Transaction Failed: ${e}`);
    }
  },

  findByUid: async function findByUid(uid) {
    const thread = await col().findOne({ suid: Uid.encode(uid) });
    return thread;
  },

  getThreadSlice: async function getThreadSlice(
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
    const result = await findSlice(sliceQuery, this, option);
    return result;
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

  getContent: function getContent() {
    return doc.blocked ? '' : doc.content;
  },

  replies: async function replies(query) {
    const option = {
      query: { threadId: doc.suid },
      field: '_id',
      sliceName: 'posts',
      parse: Uid.encode,
      toCursor: Uid.decode,
    };
    const result = await findSlice(query, mongo.collection(POST), option);
    return result;
  },

  replyCount: function replyCount() {
    return doc.catalog.length;
  },

  onPubPost: async function onPubPost(post, opt) {
    await col().updateOne({ suid: doc.suid }, {
      $push: { catalog: { postId: post.suid, createdAt: post.createdAt } },
    }, opt);
  },
});

export default ThreadModel;
*/
