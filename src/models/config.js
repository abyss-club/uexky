import mongoose from 'mongoose';
import Joi from 'joi';

import { ParamsError } from '~/utils/error';
import log from '~/utils/log';

const ConfigSchema = new mongoose.Schema({
  rateLimit: {
    httpHeader: String,
    queryLimit: Number,
    queryResetTime: Number,
    mutLimit: Number,
    mutResetTime: Number,
  },
  rateCost: {
    createUser: Number,
    pubThread: Number,
    pubPost: Number,
  },
}, { capped: 1 });

const rateLimitObjSchema = Joi.object().keys({
  httpHeader: Joi.string().regex(/^[0-9a-zA-Z-]*/).default(''),
  queryLimit: Joi.number().integer().min(0).default(300),
  queryResetTime: Joi.number().integer().min(0).default(3600),
  mutLimit: Joi.number().integer().min(0).default(30),
  mutResetTime: Joi.number().integer().min(0).default(3600),
});
const rateCostObjSchema = Joi.object().keys({
  createUser: Joi.number().integer().min(0).default(30),
  pubThread: Joi.number().integer().min(0).default(10),
  pubPost: Joi.number().integer().min(0).default(1),
});
const configObjSchema = Joi.object().keys({
  rateLimit: rateLimitObjSchema.default(rateLimitObjSchema.validate({}).value),
  rateCost: rateCostObjSchema.default(rateCostObjSchema.validate({}).value),
});

ConfigSchema.statics.getConfig = async function getConfig() {
  const results = await ConfigModel.find({}).exec();
  if (results.length === 0) {
    return configObjSchema.validate({}).value;
  }
  return results[0].format();
};

ConfigSchema.statics.setConfig = async function setConfig(input) {
  const config = await ConfigModel.getConfig();
  Object.keys(input).forEach((key) => {
    config[key] = Object.assign(config[key] || {}, input[key] || {});
  });
  log.info('will to set new config', config);
  const { error, value: newConfig } = configObjSchema.validate(config);
  if (error) {
    log.error(error);
    throw new ParamsError(`JSON validation failed, ${error}`);
  }
  await ConfigModel.updateOne({}, newConfig, { upsert: true }).exec();
  log.info('config updated!', newConfig);
  return newConfig;
};

ConfigSchema.methods.format = function format() {
  const obj = this.toObject();
  return (({ rateLimit, rateCost }) => ({ rateLimit, rateCost }))(obj);
};

const ConfigModel = mongoose.model('Config', ConfigSchema);

const genConfigReader = () => {
  let configValue = null;
  return async () => {
    if (!configValue) {
      configValue = await ConfigModel.getConfig();
    }
    return configValue;
  };
};

export default ConfigModel;
export { genConfigReader };
