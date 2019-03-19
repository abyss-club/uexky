import startRepl from '../__utils__/mongoServer';
import mongo from '~/utils/mongo';
import { Base64 } from '~/uid';
import AuthModel from '~/models/auth';

jest.setTimeout(60000);

let replSet;
let mongoClient;
let ctx; // placeholder
// let db;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

const AUTH = 'auth';

describe('Testing auth', () => {
  const authCode = Base64.randomString(36);
  const mockEmail = 'test@example.com';
  it('add user', async () => {
    const email = mockEmail;
    await AuthModel(ctx).addToAuth({ email: mockEmail, authCode });
    const result = await mongo.collection(AUTH).findOne({ email });
    expect(result.email).toEqual(mockEmail);
    expect(result.authCode).toEqual(authCode);
  });
  it('validate user authCode for only once', async () => {
    const result = await AuthModel(ctx).getEmailByCode({ authCode });
    expect(result).toEqual(mockEmail);
    const deletedResult = await mongo.collection(AUTH).findOne({ authCode });
    expect(deletedResult).toBeNull();
  });
});
