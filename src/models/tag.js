import mongoose from 'mongoose';

const TagSchema = mongoose.Schema({
  name: String,
  main_tags: [String],
  update_time: Date(),
});
const TagModel = mongoose.Model('Tag', TagSchema);
export default TagModel;
