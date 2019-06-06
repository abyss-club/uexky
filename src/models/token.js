import Joi from '@hapi/joi';
import { query } from '~/utils/pg';
import { Base64 } from '~/uid';
import { AuthError, ParamsError } from '~/utils/error';

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
    await query(
      `
      INSERT INTO token(email, "authToken") VALUES($1, $2)
      ON CONFLICT (email)
        DO UPDATE SET
          "authToken" = $2;
      `,
      [value.email, value.authToken],
    );
    return newToken.authToken;
  },

  getEmailByToken: async function getEmailByToken(authToken, refresh = false) {
    const results = await query('SELECT email, "authToken" from token where "authToken" = $1', [authToken]);
    if (results.rows.length < 1) {
      throw new AuthError('Email not found');
    }
    // if (refresh) {
    //   await col().updateOne(
    //     { authToken },
    //     { $set: { createdAt: new Date() } },
    //   );
    // }
    return results.rows[0].email;
  },
});

export default TokenModel;
