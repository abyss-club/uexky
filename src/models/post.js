import mongoose from 'mongoose';

const SchemaObjectId = mongoose.Schema.Types.ObjectId;

const PostSchema = mongoose.Schema({
  user_id: SchemaObjectId,
  thread_id: SchemaObjectId,
  anonymous: Boolean,
  author: String,
  create_time: Date,
  blocked: Boolean,
  index: Number,
  quetes: [SchemaObjectId],
  content: String,
});
const PostModel = mongoose.Model('Post', PostSchema);
export default PostModel;
