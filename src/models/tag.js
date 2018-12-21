import mongoose from 'mongoose';
import config from '~/config';

const TagSchema = mongoose.Schema({
  name: String, // unique
  mainTags: [String],
  updateAt: Date(),
});
const TagModel = mongoose.Model('Tag', TagSchema);

TagSchema.statics.onPubThread = async function onPubThread(thread, opt) {
  const option = { ...opt, upsert: true };
  thread.subTags.forEach(async (tag) => {
    await TagModel.update({ name: tag }, {
      $addToSet: { mainTags: thread.mainTag },
      $set: { updatedAt: new Date() },
    }, option);
  });
};
TagSchema.statics.getTree = async function getTree(limit) {
  const defaultLimit = 10;
  const tagTrees = config.mainTags.map(async (mainTag) => {
    const subTags = await TagModel.find({ mainTags: mainTag })
      .sort({ updatedAt: -1 })
      .limit(limit || defaultLimit).all.exec();
    return {
      mainTag,
      subTags: subTags.map(tag => tag.name),
    };
  });
  return tagTrees;
};

export default TagModel;
