import mongoose from 'mongoose';
import findSlice from '~/models/base';

const { ObjectId } = mongoose.Types;

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
  const threadUid = thread.uid();
  await NotificationModel.update({
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
      'replied.repliersIds': post.userId,
    },
  }, option);
};
NotificationSchema.statics.sendQuotedNoti = async function sendQuotedNoti(
  post, thread, quotedPosts, opt,
) {
  quotedPosts.forEach(async (qp) => {
    await NotificationModel.create({
      id: `quoted:${post.uid()}:${qp.uid()}`,
      type: 'quoted',
      sendTo: qp.userId,
      eventTime: post.createdAt,
      quoted: {
        threadId: thread.uid(),
        postId: post.uid(),
        quotedPostId: qp.uid(),
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
