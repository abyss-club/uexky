import startRepl from '../__utils__/mongoServer';
import { Base64 } from '~/uid';
import AuthModel from '~/models/auth';

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

describe('Testing auth', () => {
  const authCode = Base64.randomString(36);
  const mockEmail = 'test@example.com';
  it('add user', async () => {
    const model = AuthModel(ctx);
    await model.addToAuth(mockEmail);
    const result = await model.col().findOne({ email: mockEmail });
    expect(result.email).toEqual(mockEmail);
  });
  it('validate user authCode for only once', async () => {
    const model = AuthModel(ctx);
    const doc = await model.col().findOne({ email: mockEmail });
    const result = await model.getEmailByCode(doc.authCode);
    expect(result).toEqual(mockEmail);
    const deletedResult = await model.col().findOne({ authCode });
    expect(deletedResult).toBeNull();
  });
});
