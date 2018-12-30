import mongoose from 'mongoose';
import isEmail from 'validator/lib/isEmail';
import { encode, decode, genRandomStr } from '../utils/uuid';

const AuthSchema = new mongoose.Schema({
  email: { type: String, required: true, unique: true },
  authCode: { type: String, required: true },
  createdAt: {
    type: Date,
    index: { expireAfterSeconds: 1200 },
    required: true,
  },
});

AuthSchema.statics.addToAuth = async function addToAuth(email, code) {
  const authCode = code || genRandomStr(36);
  const newAuth = { email, authCode, createdAt: new Date() };
  await AuthModel.update({ email }, newAuth, { upsert: true });
};

AuthSchema.statics.getEmailByCode = async function getEmailByCode(authCode) {
  const result = await AuthModel.findOne({ authCode });
  if (!result) throw new Error('Auth failed');
  await AuthModel.deleteOne({ authCode }).orFail(() => new Error('Failed to invalidate authCode'));
  return result.email;
};

const AuthModel = mongoose.model('Auth', AuthSchema);

export default AuthModel;
