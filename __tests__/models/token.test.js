import mongoose from 'mongoose';
import { startMongo } from '../__utils__/mongoServer';
import TokenModel from '~/models/token';

// May require additional time for downloading MongoDB binaries
// jasmine.DEFAULT_TIMEOUT_INTERVAL = 600000;

let mongoServer;

beforeAll(async () => {
  mongoServer = await startMongo();
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
