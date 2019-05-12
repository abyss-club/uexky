import Joi from '@hapi/joi';
import mongo from '~/utils/mongo';

import { ParamsError } from '~/utils/error';
import log from '~/utils/log';

const CONFIG = 'config';
const col = () => mongo.collection(CONFIG);

const rateLimitObjSchema = Joi.object().keys({
  httpHeader: Joi.string().regex(/^[a-zA-Z0-9-]*$/).default(''),
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

const configSchema = Joi.object().keys({
  rateLimit: rateLimitObjSchema.default(rateLimitObjSchema.validate({}).value),
  rateCost: rateCostObjSchema.default(rateCostObjSchema.validate({}).value),
});

const ConfigModel = () => ({
  getConfig: async function getConfig() {
    const results = await col().find({}, { projection: { _id: 0 } }).toArray();
    if (results.length === 0) {
      return configSchema.validate({}).value;
    }
    return results[0];
  },

  setConfig: async function setConfig(input) {
    const config = await this.getConfig();
    Object.keys(input).forEach((key) => {
      config[key] = Object.assign(config[key] || {}, input[key] || {});
    });
    log.info('will to set new config', config);
    const { error, value: newConfig } = configSchema.validate(config);
    if (error) {
      log.error(error);
      throw new ParamsError(`JSON validation failed, ${error}`);
    }
    await col().findOneAndUpdate({}, { $set: newConfig }, {
      upsert: true, w: 'majority', j: true, wtimeout: 1000,
    });
    log.info('config updated!', newConfig);
    return newConfig;
  },
});

export default ConfigModel;
