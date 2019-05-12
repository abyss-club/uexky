import Joi from '@hapi/joi';
import mongo from '~/utils/mongo';
import { Base64 } from '~/uid';
import { AuthError, ParamsError } from '~/utils/error';

const TOKEN = 'token';
const col = () => mongo.collection(TOKEN);

const tokenSchema = Joi.object().keys({
  email: Joi.string().email().required(),
  authToken: Joi.string().required(),
  createdAt: Joi.date().required(),
});

const TokenModel = () => ({
  genNewToken: async function genNewToken(email) {
    const authToken = Base64.randomString(24);
    const newToken = { email, authToken, createdAt: new Date() };
    const { error, value } = tokenSchema.validate(newToken);
    if (error) {
      throw new ParamsError(`Invalid email or authCode, ${error}`);
    }
    await col().updateOne({ email }, { $set: value }, { upsert: true });
    return newToken.authToken;
  },

  getEmailByToken: async function getEmailByToken(authToken, refresh = false) {
    const results = await col().find({ authToken }).toArray();
    if (results.length === 0) {
      throw new AuthError('Email not found');
    }
    if (refresh) {
      await col().updateOne(
        { authToken },
        { $set: { createdAt: new Date() } },
      );
    }
    return results[0].email;
  },
});

export default TokenModel;
