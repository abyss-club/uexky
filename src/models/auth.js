import mongoose from 'mongoose';
import isEmail from 'validator/lib/isEmail';
import { encode, decode, genRandomStr } from '../utils/uuid';
import { getUserByEmail } from './user';

const AuthSchema = new mongoose.Schema({
  email: { type: String, required: true, unique: true },
  authCode: { type: String, required: true },
  createdAt: {
    type: Date,
    index: { expireAfterSeconds: 1200 },
    required: true,
  },
});

const AuthModel = mongoose.model('Auth', AuthSchema);

// function validateEmail(email) {
//   if (!isEmail(email)) throw new Error('Invalid email address');
// }

async function addToAuth(email) {
  const authCode = genRandomStr(36);
  const newAuth = { email, authCode, createdAt: new Date() };
  await AuthModel.update({ email }, newAuth, { upsert: true });
  const res = await AuthModel.findOne({ email }, 'authCode').orFail(() => new Error('AuthCode not found for the email'));
  console.log(res);
  return res;
}

async function getEmailByCode(authCode) {
  const email = await AuthModel.findOne({ authCode });
  if (!email) throw new Error('Auth failed');
  await AuthModel.deleteOne({ authCode }).orFail(() => new Error('Failed to invalidate authCode'));
  return email.email;
}

export default AuthModel;
export { AuthSchema, addToAuth, getEmailByCode };
