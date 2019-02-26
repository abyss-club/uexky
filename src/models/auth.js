import mongoose from 'mongoose';
import { Base64 } from '~/uid';

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
  const authCode = code || Base64.randomString(36);
  const newAuth = { email, authCode, createdAt: new Date() };
  await AuthModel.updateOne({ email }, newAuth, { upsert: true }).exec();
};

AuthSchema.statics.getEmailByCode = async function getEmailByCode(authCode) {
  const result = await AuthModel.findOne({ authCode }).exec();
  if (!result) throw new Error('Auth failed');
  await AuthModel.deleteOne({ authCode }).orFail(() => new Error('Failed to invalidate authCode')).exec();
  return result.email;
};

const AuthModel = mongoose.model('Auth', AuthSchema);

export default AuthModel;
