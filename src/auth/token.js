// import Joi from '@hapi/joi';
import getRedis from '~/utils/redis';
import { Base64 } from '~/uid';
import { ParamsError } from '~/utils/error';

import { expireTime, emailSchema } from './code';

// Token = {
//   newToken(email): give an email a token for auth
//   getEmailByToken(token, refresh=false): check email from token
// }
// In redis:
//   Key: token
//   Value: email
//   TTL: expireTime.token
const Token = {
  genNewToken: async function genNewToken(email) {
    const { error } = emailSchema.validate({ email });
    if (error) {
      throw new ParamsError(`Invalid email, ${error}`);
    }
    const token = Base64.randomString(24);
    const redis = getRedis();
    await redis.set(token, email, 'EX', expireTime.token);
    return token;
  },
  getEmailByToken: async function getEmailByToken(token, refresh = false) {
    const redis = getRedis();
    const email = await redis.get(token);
    if (!email) {
      return null;
    }
    if (refresh) {
      await redis.set(token, email, 'EX', expireTime.token);
    }
    return email;
  },
};

export default Token;
