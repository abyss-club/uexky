import mongoose from 'mongoose';
import isEmail from 'validator/lib/isEmail';
import { encode, decode, genRandomStr } from '../utils/uuid';

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
  const authToken = genRandomStr(24);
  const newToken = { email, authToken, createdAt: new Date() };
  await TokenModel.update({ email }, newToken, { upsert: true });
  const result = await TokenModel.findOne({ email }).orFail(() => new Error('AuthToken not found'));
  console.log(result);
  return result;
}

async function getEmailByToken(authToken) {
  const result = await TokenModel.findOne({ authToken }).orFail(() => new Error('Email not found'));
  console.log(result);
  return result.email;
}

export default TokenModel;
export { genNewToken, getEmailByToken };
