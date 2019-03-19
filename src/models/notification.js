import JoiBase from 'joi';
import JoiObjectId from '~/utils/joiObjectId';
import mongo from '~/utils/mongo';
import findSlice from '~/models/base';
import { ParamsError } from '~/utils/error';
import { timeZero } from '~/uid/generator';
import { ObjectId } from 'bson-ext';
import log from '~/utils/log';

import ThreadModel from './thread';
import UserModel from './user';
import PostModel from './post';

const Joi = JoiBase.extend(JoiObjectId);
const NOTIFICATION = 'notification';
const col = () => mongo.collection(NOTIFICATION);

const notiTypes = ['system', 'replied', 'quoted'];
const isValidType = type => notiTypes.findIndex(t => t === type) !== -1;
const userGroups = {
  AllUser: 'all_user',
};

const notificationSchema = Joi.object().keys({
  id: Joi.string().required(),
  type: Joi.string().valid(notiTypes).required(),
  sendTo: Joi.objectId(),
  sendToGroup: Joi.string().valid('all'),
  eventTime: Joi.date(),
  system: Joi.object().keys({
    title: Joi.string(),
    content: Joi.string(),
  }),
  replied: Joi.object().keys({
    threadId: Joi.string(),
    repliers: Joi.array().items(Joi.string()),
    replierIds: Joi.array().items(Joi.objectId()),
  }),
  quoted: {
    threadId: Joi.string(),
    postId: Joi.string(),
    quotedPostId: Joi.string(),
    quoter: Joi.string(),
    quoterId: Joi.objectId(),
  },
});

const NotificationModel = ctx => ({
  sendRepliedNoti: async function sendRepliedNoti(
    post, thread, opt,
  ) {
    const option = { ...opt, upsert: true };
    const threadUid = ThreadModel(ctx).methods(thread).uid();

    const { error } = notificationSchema.validate({
      id: `replied:${threadUid}`,
      type: 'replied',
      sendTo: thread.userId,
      eventTime: post.createdAt,
      replied: {
        threadId: threadUid,
        repliers: [post.author],
        replierIds: [post.userId],
      },
    });

    if (error) {
      log.error(error);
      throw new ParamsError(`Notification validation failed, ${error}`);
    }

    await col().updateOne({
      id: `replied:${threadUid}`,
    }, {
      $setOnInsert: {
        id: `replied:${threadUid}`,
        type: 'replied',
        sendTo: thread.userId,
        'replied.threadId': threadUid,
      },
      $set: {
        eventTime: post.createdAt,
      },
      $addToSet: {
        'replied.repliers': post.author,
        'replied.replierIds': post.userId,
      },
    }, option);
  },

  sendQuotedNoti: async function sendQuotedNoti(
    post, thread, quotedPosts, opt,
  ) {
    if (quotedPosts.length < 1) return;

    const postUid = PostModel(ctx).methods(post).uid();
    const threadUid = ThreadModel(ctx).methods(thread).uid();

    const validateQp = (qp) => {
      const qpUid = PostModel(ctx).methods(qp).uid();
      const { value, error } = notificationSchema.validate({
        id: `quoted:${postUid}:${qpUid}`,
        type: 'quoted',
        sendTo: qp.userId,
        eventTime: post.createdAt,
        quoted: {
          threadId: threadUid,
          postId: postUid,
          quotedPostId: qpUid,
          quoter: post.author,
          quoterId: post.userId,
        },
      });
      if (error) {
        log.error(error);
        throw new ParamsError(`Notification validation failed, ${error}`);
      }
      return value;
    };

    const docs = quotedPosts.map(qp => (validateQp(qp)));
    await col().insertMany(docs, opt);
  },

  getUnreadCount: async function getUnreadCount(
    user, type,
  ) {
    if (!isValidType(type)) {
      throw new ParamsError(`Invalid type: ${type}`);
    }

    let userReadNotiTime;
    if (!user.readNotiTime || !user.readNotiTime[type]) {
      userReadNotiTime = timeZero;
    } else {
      userReadNotiTime = user.readNotiTime[type];
    }

    const count = await col().find({
      $or: [
        { sendTo: user._id },
        { sendToGroup: userGroups.AllUser },
      ],
      type,
      eventTime: { $gt: userReadNotiTime },
    }).count();

    return count;
  },

  getNotiSlice: async function getNotiSlice(
    user, type, sliceQuery,
  ) {
    if (!isValidType(type)) {
      throw new ParamsError(`Invalid type: ${type}`);
    }
    if (!sliceQuery) {
      throw new ParamsError('Invalid SliceQuery.');
    }
    const option = {
      query: { type, $or: [{ sendTo: user._id }, { sendToGroup: 'all' }] },
      desc: true,
      field: '_id',
      sliceName: type,
      parse: value => ObjectId(value),
      toCursor: value => value.toHexString(),
    };
    const result = await findSlice(sliceQuery, col(), option);
    await UserModel(ctx).methods(user).setReadNotiTime(type, new Date());
    return result;
  },
});

export default NotificationModel;
export { notiTypes };
