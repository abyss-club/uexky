import mongoose from 'mongoose';

const TagSchema = new mongoose.Schema({
  name: { type: String, required: true, unique: true },
  isMain: { type: Boolean, required: true },
  belongsTo: [String],
  updateAt: Date,
}, { autoCreate: true });

TagSchema.statics.getMainTags = async function getMainTags() {
  const tags = await TagModel.find({ isMain: true }).exec();
  return tags.map(tag => tag.name);
};

TagSchema.statics.addMainTag = async function addMainTag(name) {
  const tag = await TagModel.create({ name, isMain: true });
  return tag;
};

TagSchema.statics.onPubThread = async function onPubThread(thread, opt) {
  const option = { ...opt, upsert: true };
  const updateTags = thread.subTags.map(async tag => this.updateOne({ name: tag }, {
    $addToSet: { belongsTo: thread.mainTag },
    $set: { updateAt: new Date(), isMain: false },
  }, option).exec());
  await Promise.all(updateTags);
};

TagSchema.statics.getTree = async function getTree(limit, query = '') {
  const defaultLimit = 10;
  const selector = (mainTag) => {
    if (query === '') {
      return { belongsTo: mainTag };
    }
    return { name: { $regex: query }, isMain: false, belongsTo: mainTag };
  };
  const mainTags = await TagModel.getMainTags();
  const tagTrees = mainTags.map(async (mainTag) => {
    const subTags = await this.find(selector(mainTag))
      .sort({ updatedAt: -1 })
      .limit(limit || defaultLimit)
      .exec();
    return {
      mainTag,
      subTags: subTags.map(tag => tag.name),
    };
  });
  return Promise.all(tagTrees);
};

const TagModel = mongoose.model('Tag', TagSchema);

export default TagModel;
