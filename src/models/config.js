import mongoose from 'mongoose';
import Joi from 'joi';

import { ParamsError } from '~/utils/error';
import log from '~/utils/log';

const ConfigSchema = new mongoose.Schema({
  optionName: { type: String, required: true, unique: true },
  optionValue: { type: String, required: true },
  autoload: { type: Boolean },
});

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

ConfigSchema.statics.getRateLimit = async function getRateLimit() {
  const results = await ConfigModel.find({ optionName: 'rateLimit' }).exec();
  if (results.length === 0) {
    return rateSettingSchema.validate({}).value;
  }
  return JSON.parse(results[0].optionValue);
};

ConfigSchema.statics.modifyRateLimit = async function modifyRateLimit(rateSettings) {
  const { error, value } = Joi.validate(rateSettings, rateSettingSchema);
  if (error) {
    log.error(error);
    throw new ParamsError(`JSON validation failed, ${error}`);
  }
  const newRateLimit = { optionName: 'rateLimit', optionValue: JSON.stringify(value) };
  await ConfigModel.updateOne({ optionName: 'rateLimit' }, newRateLimit, { upsert: true }).exec();
  const result = await this.getRateLimit();
  return result;
};

ConfigSchema.statics.getMainTags = async function getMainTags() {
  const results = await ConfigModel.find({ optionName: 'mainTags' }).exec();
  if (results.length === 0) {
    return [];
  }
  return JSON.parse(results[0].optionValue);
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
  await ConfigModel.updateOne({ optionName: 'mainTags' }, newMainTags, { upsert: true }).exec();
  const result = await this.getMainTags();
  return result;
};

const ConfigModel = mongoose.model('Config', ConfigSchema);

const config = () => {
  const store = {
    mainTags: null,
    rateLimit: null,
  };
  return {
    getMainTags: async function getMainTags() {
      if (!store.mainTags) {
        store.mainTags = await ConfigModel.getMainTags();
      }
      return store.mainTags;
    },
    getRateLimit: async function getRateLimit() {
      if (!store.rateLimit) {
        store.rateLimit = await ConfigModel.getRateLimit();
      }
      return store.rateLimit;
    },
  };
};

export default ConfigModel;
export { config };
