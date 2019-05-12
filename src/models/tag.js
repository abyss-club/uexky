import Joi from '@hapi/joi';
import mongo from '~/utils/mongo';
import log from '~/utils/log';
import { ParamsError } from '~/utils/error';

const TAG = 'tag';
const col = () => mongo.collection(TAG);

const tagSchema = Joi.object().keys({
  name: Joi.string().required(),
  isMain: Joi.boolean().required(),
  belongsTo: Joi.array().items(Joi.string()),
  updateAt: Joi.date(),
});

const TagModel = () => ({
  getMainTags: async function getMainTags() {
    const tags = await col().find({ isMain: true }).toArray();
    return tags.map(tag => tag.name);
  },

  addMainTag: async function addMainTag(name) {
    const { value, error } = tagSchema.validate({ name, isMain: true });
    if (error) {
      log.error(error);
      throw new ParamsError(`Tag validation failed, ${error}`);
    }
    const tag = await col().insertOne(value);
    return tag;
  },

  onPubThread: async function onPubThread(thread, opt) {
    const option = { ...opt, upsert: true };
    const updateTags = thread.subTags.map(async tag => col().updateOne({ name: tag }, {
      $addToSet: { belongsTo: thread.mainTag },
      $set: { updateAt: new Date(), isMain: false },
    }, option));
    await Promise.all(updateTags);
  },

  getTree: async function getTree(limit, query = '') {
    const defaultLimit = 10;
    const selector = (mainTag) => {
      if (query === '') {
        return { belongsTo: mainTag };
      }
      return { name: { $regex: query }, isMain: false, belongsTo: mainTag };
    };
    const mainTags = await this.getMainTags();
    const tagTrees = mainTags.map(async (mainTag) => {
      const subTags = await col().find(selector(mainTag))
        .sort({ updatedAt: -1 })
        .limit(limit || defaultLimit)
        .toArray();
      return {
        mainTag,
        subTags: subTags.map(tag => tag.name),
      };
    });
    return Promise.all(tagTrees);
  },
});


export default TagModel;
