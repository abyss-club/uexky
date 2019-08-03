import Joi from '@hapi/joi';
import { Base64 } from '~/uid';
import sendAuthMail from './mail';
import { AuthError, ParamsError } from '~/utils/error';
import getRedis from '~/utils/redis';
import log from '~/utils/log';


const emailSchema = Joi.object().keys({
  email: Joi.string().email().required(),
});

const expireTime = {
  code: 1200, // 20 minutes
  token: 86400 * 7, // one week
};

// Code = {
//   addToAuth(email): send to login mail to email address, with link included code
//   getEmailByCode(authCode): read code, find out which email want to login
// }
// In redis:
//   Key: code
//   Value: email
//   TTL: expireTime.code

const Code = {
  addToAuth: async function addToAuth(email) {
    const { error } = emailSchema.validate({ email });
    if (error) {
      throw new ParamsError(`Invalid email, ${error}`);
    }
    const code = Base64.randomString(36);
    const redis = getRedis();
    try {
      await redis.set(code, email, 'EX', expireTime.code);
    } catch (redisErr) {
      log.error('set redis error', redisErr);
    }
    await sendAuthMail(email, code);
    return code;
  },
  getEmailByCode: async function getEmailByCode(code) {
    const redis = getRedis();
    const email = await redis.get(code);
    if (!email) {
      throw new AuthError('Corresponding email not found.');
    }
    await redis.del(code);
    return email;
  },
};

export default Code;
export { expireTime, emailSchema };
