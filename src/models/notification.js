import mongoose from 'mongoose';

import findSlice from '~/models/base';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;
const notiTypes = ['system', 'replied', 'quoted'];
const isValidType = type => notiTypes.findIndex(t => t === type) !== -1;

const NotificationSchema = new mongoose.Schema({
  id: [String],
  type: { type: String, enum: notiTypes },
  sendTo: SchemaObjectId,
  sendToGroup: { type: String, enum: ['all'] },
  eventTime: Date,
  system: {
    title: String,
    content: String,
  },
  replied: {
    threadId: SchemaObjectId,
    repliers: [String],
    repliersIds: [SchemaObjectId],
  },
  quoted: {
    threadId: SchemaObjectId,
    postId: SchemaObjectId,
    quotedPostId: SchemaObjectId,
    quoter: String,
    quoterId: SchemaObjectId,
  },
}, { id: false });

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
  await NotificationModel.update({
    id: `replied:${thread.id()}`,
  }, {
    $setOnInsert: {
      id: `replied:${thread.id()}`,
      type: 'replied',
      sendTo: thread.userId,
      'replied.threadId': thread._id,
    },
    $set: {
      eventTime: post.createdAt,
    },
    $addToSet: {
      'replied.repliers': post.author,
      'replied.repliersIds': post.userId,
    },
  }, option);
};
NotificationSchema.statics.sendQuotedNoti = async function sendQuotedNoti(
  post, thread, quotedPosts, opt,
) {
  quotedPosts.forEach(async (qp) => {
    await NotificationModel.create({
      id: `quoted:${post.id()}:${qp.id()}`,
      type: 'quoted',
      sendTo: qp.userId,
      eventTime: post.createdAt,
      quoted: {
        threadId: thread._id,
        postId: post._id,
        quotedPostId: qp._id,
        quoter: post.author,
        quoterId: post.userId,
      },
    }, opt);
  });
};
NotificationSchema.statics.getNotiSlice = async function getNotiSlice(
  user, type, sliceQuery,
) {
  if (!isValidType(type)) {
    throw new Error(`invalid type: ${type}`);
  }
  const option = {
    query: { $or: [{ sendTo: user._id }, { sendToGroup: 'all' }] },
    desc: true,
    field: 'eventTime',
    sliceName: type,
    parse: value => (new Date(value)),
    toCursor: value => (value.toISOString()),
  };
  const result = await findSlice(sliceQuery, NotificationModel, option);
  return result;
};

const NotificationModel = mongoose.model('Notification', NotificationSchema);
export default NotificationModel;
