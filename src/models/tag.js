import mongoose from 'mongoose';
import ConfigModel from './config';

const TagSchema = new mongoose.Schema({
  subTag: { type: String, required: true, unique: true },
  mainTags: [String],
  updateAt: Date,
}, { autoCreate: true });

TagSchema.statics.onPubThread = async function onPubThread(thread, opt) {
  const option = { ...opt, upsert: true };
  const updateTags = thread.subTags.map(async tag => this.updateOne({ subTag: tag }, {
    $addToSet: { mainTags: thread.mainTag },
    $set: { updateAt: new Date() },
  }, option).exec());
  await Promise.all(updateTags);
};

TagSchema.statics.getTree = async function getTree(limit, query = '') {
  const defaultLimit = 10;
  const selector = (mainTag) => {
    if (query === '') {
      return { mainTags: mainTag };
    }
    return { subTag: { $regex: query }, mainTags: mainTag };
  };
  const mainTags = await ConfigModel.getMainTags();
  const tagTrees = mainTags.map(async (mainTag) => {
    const subTags = await this.find(selector(mainTag))
      .sort({ updatedAt: -1 })
      .limit(limit || defaultLimit)
      .exec();
    return {
      mainTag,
      subTags: subTags.map(tag => tag.subTag),
    };
  });
  return Promise.all(tagTrees);
};

const TagModel = mongoose.model('Tag', TagSchema);

export default TagModel;
