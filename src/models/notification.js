import mongoose from 'mongoose';
import findSlice from '~/models/base';
import { ParamsError } from '~/utils/error';

const { ObjectId } = mongoose.Types;
const SchemaObjectId = mongoose.ObjectId;
const notiTypes = ['system', 'replied', 'quoted'];
const isValidType = type => notiTypes.findIndex(t => t === type) !== -1;
const userGroups = {
  AllUser: 'all_user',
};

const NotificationSchema = new mongoose.Schema({
  id: String,
  type: { type: String, enum: notiTypes },
  sendTo: SchemaObjectId,
  sendToGroup: { type: String, enum: ['all'] },
  eventTime: Date,
  system: {
    title: String,
    content: String,
  },
  replied: {
    threadId: String,
    repliers: [String],
    replierIds: [SchemaObjectId],
  },
  quoted: {
    threadId: String,
    postId: String,
    quotedPostId: String,
    quoter: String,
    quoterId: SchemaObjectId,
  },
}, { id: false, autoCreate: true });

NotificationSchema.methods.body = function body() {
  switch (this.type) {
    case 'system':
      return this.system;
    case 'replied':
      return this.replied;
    case 'quoted':
      return this.quoted;
    default:
      return null;
  }
};

NotificationSchema.statics.sendRepliedNoti = async function sendRepliedNoti(
  post, thread, opt,
) {
  const option = { ...opt, upsert: true };
  const threadUid = thread.uid();
  await NotificationModel.updateOne({
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
};

NotificationSchema.statics.sendQuotedNoti = async function sendQuotedNoti(
  post, thread, quotedPosts, opt,
) {
  const postUid = post.uid();
  const threadUid = thread.uid();
  await Promise.all(quotedPosts.map(async (qp) => {
    NotificationModel.create({
      id: `quoted:${postUid}:${qp.uid()}`,
      type: 'quoted',
      sendTo: qp.userId,
      eventTime: post.createdAt,
      quoted: {
        threadId: threadUid,
        postId: postUid,
        quotedPostId: qp.uid(),
        quoter: post.author,
        quoterId: post.userId,
      },
    }, opt);
  }));
};

NotificationSchema.statics.getUnreadCount = async function getUnreadCount(
  user, type,
) {
  if (!isValidType(type)) {
    throw new ParamsError(`Invalid type: ${type}`);
  }
  const count = await NotificationModel.find({
    $or: [
      { send_to: user.ID },
      { send_to_group: userGroups.AllUser },
    ],
    type,
    eventTime: { $gt: user.readNotiTime[type] },
  }).count().exec();
  return count;
};

NotificationSchema.statics.getNotiSlice = async function getNotiSlice(
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
    toCursor: value => value.valueOf(),
  };
  const result = await findSlice(sliceQuery, NotificationModel, option);
  return result;
};

const NotificationModel = mongoose.model('Notification', NotificationSchema);
export default NotificationModel;
export { notiTypes };
