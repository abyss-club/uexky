import mongoose from 'mongoose';
import Joi from 'joi';
import { InternalError, ParamsError } from '~/error';

const ConfigSchema = new mongoose.Schema({
  optionName: { type: String, required: true, unique: true },
  optionValue: { type: String, required: true },
  autoload: { type: Boolean },
});

ConfigSchema.statics.getRateLimit = async function getRateLimit() {
  const result = await ConfigModel.findOne({ optionName: 'rateLimit' }, 'optionValue').orFail(() => new InternalError('rateLimit not yet set.'));
  return result.optionValue;
};

ConfigSchema.statics.modifyRateLimit = async function modifyRateLimit(rateSettings) {
  const costSchema = Joi.object().keys({
    CreateUser: Joi.number().integer().min(0).default(30),
    PubThread: Joi.number().integer().min(0).default(10),
    PubPost: Joi.number().integer().min(0).default(1),
  });
  const rateSettingSchema = Joi.object().keys({
    HTTPHeader: Joi.string().alphanum().default(''),
    QueryLimit: Joi.number().integer().min(0).default(300),
    QueryResetTime: Joi.number().integer().min(0).default(3600),
    MutLimit: Joi.number().integer().min(0).default(30),
    MutResetTime: Joi.number().integer().min(0).default(3600),
    Cost: costSchema.default(costSchema.validate({}).value),
  });
  const { error, value } = Joi.validate(rateSettings, rateSettingSchema);
  // console.log(error);
  if (error) throw new ParamsError(`JSON validation failed, ${error}`);
  const newRateLimit = { optionName: 'rateLimit', optionValue: JSON.stringify(value) };
  await ConfigModel.updateOne({ optionName: 'rateLimit' }, newRateLimit, { upsert: true });
  const result = await this.getRateLimit();
  return result;
};

ConfigSchema.statics.getMainTags = async function getMainTags() {
  const result = await ConfigModel.findOne({ optionName: 'mainTags' }, 'optionValue').orFail(() => new InternalError('mainTags not yet set.'));
  return JSON.parse(result.optionValue);
};

ConfigSchema.statics.modifyMainTags = async function modifyMainTags(tags) {
  if (!Array.isArray(tags) || !tags.length) {
    throw new ParamsError('Provided tags is not a non-empty array.');
  }
  tags.forEach((tag) => {
    if (Object.prototype.toString.call(tag) !== '[object String]' || tag === '') {
      throw new ParamsError('Invalid tag provided in array.');
    }
  });
  const newMainTags = { optionName: 'mainTags', optionValue: JSON.stringify(tags) };
  await ConfigModel.updateOne({ optionName: 'mainTags' }, newMainTags, { upsert: true });
  const result = await this.getMainTags();
  return result;
};

const ConfigModel = mongoose.model('Config', ConfigSchema);

const config = () => {
  const store = {
    mainTags: undefined,
    rateLimit: undefined,
  };
  return {
    getMainTags: async function getMainTags() {
      if (!store.mainTags) {
        const mainTags = await ConfigModel.getMainTags();
        store.mainTags = mainTags;
      }
      return store.mainTags;
    },
    getRateLimit: async function getRateLimit() {
      if (!store.rateLimit) {
        const rateLimit = await ConfigModel.getRateLimit();
        store.rateLimit = rateLimit;
      }
      return this.rateLimit;
    },
  };
};

export default config;
export { ConfigModel };
