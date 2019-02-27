import mongoose from 'mongoose';
import Joi from 'joi';

import { ParamsError } from '~/utils/error';
import log from '~/utils/log';

const ConfigSchema = new mongoose.Schema({
  rateLimit: {
    httpHeader: String,
    queryLimit: String,
    queryResetTime: Number,
    mutLimit: String,
    mutResetTime: Number,
  },
  rateCost: {
    createUser: Number,
    pubThread: Number,
    pubPost: Number,
  },
}, { capped: 1 });

const rateCostObjSchema = Joi.object().keys({
  CreateUser: Joi.number().integer().min(0).default(30),
  PubThread: Joi.number().integer().min(0).default(10),
  PubPost: Joi.number().integer().min(0).default(1),
});
const rateLimitObjSchema = Joi.object().keys({
  HTTPHeader: Joi.string().alphanum().default(''),
  QueryLimit: Joi.number().integer().min(0).default(300),
  QueryResetTime: Joi.number().integer().min(0).default(3600),
  MutLimit: Joi.number().integer().min(0).default(30),
  MutResetTime: Joi.number().integer().min(0).default(3600),
});
const configObjSchema = Joi.object().keys({
  rateLimit: rateLimitObjSchema.default(rateLimitObjSchema.validate({}).value),
  rateCost: rateCostObjSchema.default(rateCostObjSchema.validate({}).value),
});

ConfigSchema.statics.getConfig = async function getConfig() {
  const results = await ConfigModel.find({}).exec();
  if (results.length === 0) {
    return configObjSchema.validate({});
  }
  return results[0];
};

ConfigSchema.statics.setConfig = async function setConfig(input) {
  const config = await ConfigModel.getConfig();
  Object.keys(config).forEach((key) => {
    Object.assign(config[key], input[key]);
  });
  const { error, value: newConfig } = configObjSchema.validate(config);
  if (error) {
    log.error(error);
    throw new ParamsError(`JSON validation failed, ${error}`);
  }
  await ConfigModel.updateOne({}, newConfig).exec();
  return newConfig;
};

const ConfigModel = mongoose.model('Config', ConfigSchema);

const getConfig = () => {
  let configValue = null;
  return async () => {
    if (!configValue) {
      configValue = await ConfigModel.getConfig();
    }
    return configValue;
  };
};

export default ConfigModel;
export { getConfig };
