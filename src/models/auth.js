import Joi from 'joi';
import mongo from '~/utils/mongo';
import { Base64 } from '~/uid';
import { AuthError, ParamsError } from '~/utils/error';

const authSchema = Joi.object().keys({
  email: Joi.string().email().required(),
  authCode: Joi.string().required(),
  createdAt: Joi.date().required(),
});

const AUTH = 'authCode';
const col = () => mongo.collection(AUTH);

const expireTime = {
  code: 1200, // 20 minutes
  token: 86400 * 7, // one week
};

const AuthModel = () => ({
  addToAuth: async function addToAuth({ email, authCode = Base64.randomString(36) }) {
    const newAuth = { email, authCode, createdAt: new Date() };
    const { error, value } = authSchema.validate(newAuth);
    if (error) {
      throw new ParamsError(`Invalid email or authCode, ${error}`);
    }
    await col().updateOne({ email }, { $set: value }, { upsert: true });
  },
  getEmailByCode: async function getEmailByCode({ authCode }) {
    const results = await col().find({ authCode }).toArray();
    if (results.length !== 0) {
      throw new AuthError('Corresponding email not found.');
    }
    try {
      await col().deleteOne({ authCode });
    } catch (e) {
      throw new Error('Failed to invalidate authCode.');
    }
    return results[0].email;
  },
});

export default AuthModel;
export { expireTime };
