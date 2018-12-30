import mongoose from 'mongoose';
import config from '~/config';

const TagSchema = new mongoose.Schema({
  subTag: { type: String, required: true, unique: true },
  mainTags: [String],
  updateAt: Date,
});

TagSchema.statics.onPubThread = async function onPubThread(thread, opt) {
  const option = { ...opt, upsert: true };
  thread.subTags.forEach(async (tag) => {
    console.log(tag);
    await TagModel.update({ subTag: tag }, {
      $addToSet: { mainTags: thread.mainTag },
      $set: { updateAt: new Date() },
    }, option);
  });
};

TagSchema.statics.getTree = async function getTree(limit) {
  const defaultLimit = 10;
  const tagTrees = config.mainTags.map(async (mainTag) => {
    const subTags = await TagModel.find({ mainTags: mainTag })
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
