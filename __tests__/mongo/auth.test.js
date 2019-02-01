import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';

import { Base64 } from '~/uid';
import AuthModel from '~/models/auth';

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
  const authCode = Base64.randomString(36);
  const mockEmail = 'test@example.com';
  it('add user', async () => {
    const email = mockEmail;
    // const newAuth = { email, authCode, createdAt: new Date() };
    // await AuthModel.update({ email }, newAuth, { upsert: true });
    await AuthModel.addToAuth(mockEmail, authCode);
    const result = await AuthModel.findOne({ email });
    expect(result.email).toEqual(mockEmail);
    expect(result.authCode).toEqual(authCode);
  });
  it('validate user authCode for only once', async () => {
    const result = await AuthModel.getEmailByCode(authCode);
    expect(result).toEqual(mockEmail);
    const deletedResult = await AuthModel.findOne({ authCode });
    expect(deletedResult).toBeNull();
  });
});
