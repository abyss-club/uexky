// import Joi from 'joi';
import mongo from '~/utils/mongo';
import { AuthError } from '~/utils/error';
import { Base64 } from '~/uid';

const TOKEN = 'token';
const col = () => mongo.collection(TOKEN);

// const tagSchema = Joi.object().keys({
//   email: Joi.string().email().required(),
//   authToken: Joi.string().required(),
//   createdAt: Joi.date().required(),
// });

const TokenModel = () => ({
  genNewToken: async function genNewToken(email) {
    const authToken = Base64.randomString(24);
    const newToken = { email, authToken, createdAt: new Date() };
    await col().updateOne({ email }, { $set: { newToken } }, { upsert: true });
    try {
      const result = await col().findOne({ email });
      return result;
    } catch (e) { throw new AuthError('AuthToken not found'); }
  },

  getEmailByToken: async function getEmailByToken(authToken) {
    try {
      const result = await col().findOne({ authToken });
      return result.email;
    } catch (e) { throw new AuthError('Email not found'); }
  },
});

export default TokenModel;
