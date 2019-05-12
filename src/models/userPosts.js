import JoiBase from '@hapi/joi';
import JoiObjectId from '~/utils/joiObjectId';
import mongo from '~/utils/mongo';
import log from '~/utils/log';
import { ParamsError } from '~/utils/error';

const Joi = JoiBase.extend(JoiObjectId);
const USERPOSTS = 'userPosts';
const col = () => mongo.collection(USERPOSTS);

const userPostsSchema = Joi.object().keys({
  userId: Joi.objectId().required(),
  threadSuid: Joi.string().alphanum().length(15),
  updatedAt: Joi.date(),
  posts: Joi.array().items(Joi.object().keys({
    suid: Joi.string().alphanum().length(15),
    createdAt: Joi.date(),
    anonymous: Joi.boolean(),
  })).default([]),
});

// MODEL: UserPosts
//        used for querying user's posts grouped by thread.
// const UserPostsSchema = new mongoose.Schema({
//   userId: SchemaObjectId,
//   threadSuid: String,
//   updatedAt: Date,
//   posts: [{
//     suid: String,
//     createdAt: Date,
//     anonymous: Boolean,
//   }],
// }, { autoCreate: true });
const UserPostsModel = ctx => ({
  methods: function methods(doc) {
    return genDoc(ctx, this, doc);
  },
});

const genDoc = (ctx, model, doc) => ({
  onPubThread: async function onPubThread(thread, { session }) {
    const { value, error } = userPostsSchema.validate({
      userId: doc._id,
      threadSuid: thread.suid,
    });
    if (error) {
      log.error(error);
      throw new ParamsError(`UserPosts validation failed, ${error}`);
    }

    await col().updateOne({ ...value }, {
      $setOnInsert: { ...value },
      $set: { updatedAt: new Date() },
    }, { session, upsert: true });
  },

  onPubPost: async function onPubPost(
    thread, post, { session },
  ) {
    const { error } = userPostsSchema.validate({
      userId: doc._id,
      threadSuid: thread.suid,
      posts: [{ suid: post.suid, anonymous: post.anonymous, createdAt: post.createdAt }],
    });
    if (error) {
      log.error(error);
      throw new ParamsError(`UserPosts validation failed, ${error}`);
    }

    await col().updateOne({
      userId: doc._id,
      threadSuid: thread.suid,
    }, {
      $setOnInsert: {
        userId: doc._id,
        threadSuid: thread.suid,
      },
      $set: { updatedAt: new Date() },
      $push: {
        posts: {
          suid: post.suid,
          createdAt: post.createdAt,
          anonymous: post.anonymous,
        },
      },
    }, { session, upsert: true });
  },
});

export default UserPostsModel;
