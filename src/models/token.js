import mongoose from 'mongoose';
import { AuthError } from '~/error';
import { Base64 } from '~/uid';

const TokenSchema = new mongoose.Schema({
  email: { type: String, required: true, unique: true },
  authToken: { type: String, required: true },
  createdAt: {
    type: Date,
    // 20 days
    index: { expireAfterSeconds: 1728000 },
    required: true,
  },
});

const TokenModel = mongoose.model('Token', TokenSchema);

async function genNewToken(email) {
  const authToken = Base64.randomString(24);
  const newToken = { email, authToken, createdAt: new Date() };
  await TokenModel.update({ email }, newToken, { upsert: true });
  const result = await TokenModel.findOne({ email }).orFail(() => new AuthError('AuthToken not found'));
  return result;
}

async function getEmailByToken(authToken) {
  const result = await TokenModel.findOne({ authToken }).orFail(() => new AuthError('Email not found'));
  return result.email;
}

export default TokenModel;
export { genNewToken, getEmailByToken };
