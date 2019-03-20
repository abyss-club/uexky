import Joi from 'joi';
import mongo from '~/utils/mongo';
import { Base64 } from '~/uid';
import { AuthError, ParamsError } from '~/utils/error';

const authSchema = Joi.object().keys({
  email: Joi.string().email().required(),
  authCode: Joi.string().required(),
  createdAt: Joi.date().required(),
});

const expireTime = {
  code: 1200, // 20 minutes
  token: 86400 * 7, // one week
};

const AuthModel = () => ({
  col: function col() {
    return mongo.collection('authCode');
  },
  addToAuth: async function addToAuth(email) {
    const authCode = Base64.randomString(36);
    const newAuth = { email, authCode, createdAt: new Date() };
    const { error, value } = authSchema.validate(newAuth);
    if (error) {
      throw new ParamsError(`Invalid email or authCode, ${error}`);
    }
    await this.col().updateOne({ email }, { $set: value }, { upsert: true });
  },
  getEmailByCode: async function getEmailByCode(authCode) {
    const result = await this.col().findOne({ authCode });
    if (!result) {
      throw new AuthError('Corresponding email not found.');
    }
    try {
      await this.col().deleteOne({ authCode });
    } catch (e) {
      throw new Error('Failed to invalidate authCode.');
    }
    return result.email;
  },
});

export default AuthModel;
export { expireTime };
