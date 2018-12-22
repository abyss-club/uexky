import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';
import { genRandomStr } from '~/utils/uuid';
import { AuthSchema } from '~/models/auth';
import { TokenSchema } from '~/models/token';
import { UserSchema } from '~/models/user';

const AuthModel = mongoose.model('Auth', AuthSchema);
const TokenModel = mongoose.model('Token', TokenSchema);
const UserModel = mongoose.model('User', UserSchema);

// May require additional time for downloading MongoDB binaries
// jasmine.DEFAULT_TIMEOUT_INTERVAL = 600000;

let mongoServer;
const opts = { useNewUrlParser: true };

beforeAll(async () => {
  mongoServer = new MongoMemoryServer();
  const mongoUri = await mongoServer.getConnectionString();
  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) console.error(err);
  });
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});

describe('Testing auth', () => {
  const authCode = genRandomStr(36);
  const mockEmail = 'test@example.com';
  it('add user', async () => {
    const email = mockEmail;
    const newAuth = { email, authCode, createdAt: new Date() };
    await AuthModel.update({ email }, newAuth, { upsert: true });
    const result = await AuthModel.findOne({ email });
    expect(result.email).toEqual(mockEmail);
    expect(result.authCode).toEqual(authCode);
  });
  it('validate user authCode for only once', async () => {
    const result = await AuthModel.findOne({ authCode });
    expect(result.email).toEqual(mockEmail);
    await AuthModel.deleteOne({ authCode });
    const deletedResult = await AuthModel.findOne({ authCode });
    expect(deletedResult).toBeNull();
  });
});

// describe('Testing auth after 20min', () => {
//   const sleep = m => new Promise(r => setTimeout(r, m))
//   const authCode = genRandomStr(36);
//   const mockEmail = 'fail@example.com';
//   it('add user', async () => {
//     const email = mockEmail;
//     const outdatedDate = new Date(Date.now() - 1000 * 1201);
//     const newAuth = { email, authCode, createdAt: outdatedDate };
//     await AuthModel.update({ email }, newAuth, { upsert: true });
//     const result = await AuthModel.findOne({ email });
//     expect(result.email).toEqual(mockEmail);
//     expect(result.authCode).toEqual(authCode);
//   });
//   it('validate user authCode for only once', async () => {
//     await sleep(3000);
//     const result = await AuthModel.findOne({ authCode });
//     expect(result.email).toEqual(mockEmail);
//     await AuthModel.deleteOne({ authCode });
//     const deletedResult = await AuthModel.findOne({ authCode });
//     expect(deletedResult).toBeNull();
//   });
// });
