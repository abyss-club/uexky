import Joi from 'joi';
import dbClient from '~/dbClient';
import { Base64 } from '~/uid';
import { AuthError, ParamsError } from '~/utils/error';

const authSchema = Joi.object().keys({
  email: Joi.string().email().required(),
  authCode: Joi.string().required(),
  createdAt: Joi.date().required(),
});

const AUTH = 'auth';
const col = () => dbClient.collection(AUTH);

const AuthModel = () => ({
  addToAuth: async function addToAuth({ email, authCode = Base64.randomString(36) }) {
    const newAuth = { email, authCode, createdAt: new Date() };
    const { error, value } = authSchema.validate(newAuth);
    if (error) throw new ParamsError(`Invalid email or authCode, ${error}`);
    await col().updateOne({ email }, { $set: value }, { upsert: true });
  },
  getEmailByCode: async function getEmailByCode({ authCode }) {
    const result = await col().findOne({ authCode });
    if (!result) throw new AuthError('Corresponding email not found.');
    try { await col().deleteOne({ authCode }); } catch (e) { throw new Error('Failed to invalidate authCode.'); }
    return result.email;
  },
});

export default AuthModel;
