import Joi from '@hapi/joi';

import { query } from '~/utils/pg';
import { ParamsError } from '~/utils/error';
import log from '~/utils/log';
import { ACTION } from '~/models/user';

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

const ConfigModel = {
  async getConfig() {
    const results = await query('SELECT "rateLimit", "rateCost" from config where id = 1');
    if (results.rows.length < 1) {
      return configSchema.validate({}).value;
    }
    return results.rows[0];
  },

  async setConfig(ctx, input) {
    ctx.ensurePermission(ACTION.EDIT_SETTING);
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
    await query(
      `
      INSERT INTO config(id, "rateLimit", "rateCost") VALUES(1, $1, $2)
      ON CONFLICT (id)
        DO UPDATE SET
          "rateLimit" = $1,
          "rateCost" = $2;
      `,
      [newConfig.rateLimit, newConfig.rateCost],
    );
    log.info('config updated!', newConfig);
    return newConfig;
  },
};

export default ConfigModel;
