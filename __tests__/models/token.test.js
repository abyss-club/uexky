import { startRepl } from '../__utils__/mongoServer';
import TokenModel from '~/models/token';

jest.setTimeout(60000);

let replSet;
let mongoClient;
let ctx;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

describe('Testing token', () => {
  const mockEmail = 'test@example.com';
  it('validate token by email', async () => {
    const tokenResult = await TokenModel(ctx).genNewToken(mockEmail);
    const emailResult = await TokenModel(ctx).getEmailByToken(tokenResult.authToken);
    expect(emailResult).toEqual(mockEmail);
  });
});
