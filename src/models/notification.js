import mongoose from 'mongoose';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

const NotificationSchema = mongoose.Schema({
  id: [String],
  type: { type: String, enum: ['system', 'replied', 'quoted'] },
  send_to: SchemaObjectId,
  send_to_group: { type: String, enum: ['all'] },
  event_time: Date,
  system: {
    title: String,
    content: String,
  },
  replied: {
    thread_id: SchemaObjectId,
    repliers: [String],
    repliers_ids: [SchemaObjectId],
  },
  quoted: {
    thread_id: SchemaObjectId,
    post_id: SchemaObjectId,
    quoted_post_id: SchemaObjectId,
    quoter: String,
    quoter_id: SchemaObjectId,
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
const NotificationModel = mongoose.Model('Notification', NotificationSchema);
export default NotificationModel;
