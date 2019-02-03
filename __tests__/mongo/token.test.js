import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';
import TokenModel from '~/models/token';

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

describe('Testing token', () => {
  const mockEmail = 'test@example.com';
  it('validate token by email', async () => {
    const tokenResult = await TokenModel.genNewToken(mockEmail);
    const emailResult = await TokenModel.getEmailByToken(tokenResult.authToken);
    expect(emailResult).toEqual(mockEmail);
  });
});
