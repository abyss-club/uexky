import mongoose from 'mongoose';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

const NotificationSchema = mongoose.Schema({
  id: [String],
  type: { type: String, enum: ['system', 'replied', 'quoted'] },
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
const NotificationModel = mongoose.Model('Notification', NotificationSchema);

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
export default NotificationModel;
