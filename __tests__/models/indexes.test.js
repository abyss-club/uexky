import startRepl from '../__utils__/mongoServer';
import mongo from '~/utils/mongo';

import createIndexes from '~/models/indexes';

jest.setTimeout(60000);

let replSet;
let mongoClient;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

const expectedIndexes = {
  auth: [
    { key: { email: 1 }, unique: true },
    { key: { authCode: 1 } },
    { key: { createdAt: 1 }, expireAfterSeconds: 1200 },
  ],
  notification: [
    { key: { send_to: 1, type: 1, eventTime: 1 } },
    {
      key: { send_to_group: 1, type: 1, eventTime: 1 },
      partialFilterExpression: { send_to_group: { $exists: true } },
    },
  ],
  post: [
    { key: { suid: 1 }, unique: true },
    { key: { quoteSuids: 1 } },
  ],
  thread: [
    { key: { suid: 1 }, unique: true },
    { key: { tags: 1, suid: -1 } },
  ],
  token: [
    { key: { email: 1 }, unique: true },
    { key: { authToken: 1 }, unique: true },
    { key: { createdAt: 1 }, expireAfterSeconds: 172800 },
  ],
  user: [
    { key: { email: 1 }, unique: true },
    { key: { name: 1 }, unique: true, partialFilterExpression: { name: { $type: 'string' } } },
  ],
  userAid: [
    { key: { userId: 1, threadSuid: 1 }, unique: true },
  ],
  userPosts: [
    { key: { userId: 1, threadSuid: 1, updatedAt: -1 } },
  ],
};

describe('test creating indexes', () => {
  it('Creating indexes', async () => {
    await createIndexes();
  });
  it('Testing auth indexes', async () => {
    const indexes = await mongo.collection('auth').indexes();
    expect(indexes[1]).toMatchObject(expectedIndexes.auth[0]);
    expect(indexes[2]).toMatchObject(expectedIndexes.auth[1]);
    expect(indexes[3]).toMatchObject(expectedIndexes.auth[2]);
  });
  it('Testing notification indexes', async () => {
    const indexes = await mongo.collection('notification').indexes();
    expect(indexes[1]).toMatchObject(expectedIndexes.notification[0]);
    expect(indexes[2]).toMatchObject(expectedIndexes.notification[1]);
  });
  it('Testing post indexes', async () => {
    const indexes = await mongo.collection('post').indexes();
    expect(indexes[1]).toMatchObject(expectedIndexes.post[0]);
    expect(indexes[2]).toMatchObject(expectedIndexes.post[1]);
  });
  it('Testing thread indexes', async () => {
    const indexes = await mongo.collection('thread').indexes();
    expect(indexes[1]).toMatchObject(expectedIndexes.thread[0]);
    expect(indexes[2]).toMatchObject(expectedIndexes.thread[1]);
  });
  it('Testing token indexes', async () => {
    const indexes = await mongo.collection('token').indexes();
    expect(indexes[1]).toMatchObject(expectedIndexes.token[0]);
    expect(indexes[2]).toMatchObject(expectedIndexes.token[1]);
    expect(indexes[3]).toMatchObject(expectedIndexes.token[2]);
  });
  it('Testing user indexes', async () => {
    const indexes = await mongo.collection('user').indexes();
    expect(indexes[1]).toMatchObject(expectedIndexes.user[0]);
    expect(indexes[2]).toMatchObject(expectedIndexes.user[1]);
  });
  it('Testing userAid indexes', async () => {
    const indexes = await mongo.collection('userAid').indexes();
    expect(indexes[1]).toMatchObject(expectedIndexes.userAid[0]);
  });
  it('Testing userPosts indexes', async () => {
    const indexes = await mongo.collection('userPosts').indexes();
    expect(indexes[1]).toMatchObject(expectedIndexes.userPosts[0]);
  });
});
