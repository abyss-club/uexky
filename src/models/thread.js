import mongoose from 'mongoose';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

const ThreadSchema = mongoose.Schema({
  anonymouse: Boolean,
  author: String,
  user_id: SchemaObjectId,
  main_tag: String,
  sub_tags: [String],
  tags: [String],
  title: String,
  locked: Boolean,
  blocked: Boolean,
  create_time: Date,
  update_time: Date,
  content: String,
});
const ThreadModel = mongoose.Model('Thread', ThreadSchema);
export default ThreadModel;
