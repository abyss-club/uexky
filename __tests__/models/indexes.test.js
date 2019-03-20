import startRepl from '../__utils__/mongoServer';
import mongo from '~/utils/mongo';

import { index } from '~/models/indexes';

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

describe('test creating indexes', () => {
  it('creating indexes', async () => {
    const indexes = [
      { key: { send_to: 1, type: 1, eventTime: 1 } },
      {
        key: { send_to_group: 1, type: 1, eventTime: 1 },
        partialFilterExpression: { send_to_group: { $exists: true } },
      },
    ];
    const expects = indexes.map(idx => ({ ...idx, background: true }));
    await index('indexTest', indexes);
    const results = await mongo.collection('indexTest').indexes();
    console.log('indexes', results);
    expect(results[1]).toMatchObject(expects[0]);
    expect(results[2]).toMatchObject(expects[1]);
  });
});
