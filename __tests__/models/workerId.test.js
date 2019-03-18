import { startRepl } from '../__utils__/mongoServer';
import { newWorkerId } from '~/models/workerId';

jest.setTimeout(60000); // for boot replica sets
let replSet;
let mongoClient;
// let db;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

// const WORKER_ID = 'workerid';

test('get worker id', async () => {
  const ids = await Promise.all([1, 2, 3, 4, 5].map(() => newWorkerId()));
  expect(ids.sort()).toEqual([1, 2, 3, 4, 5]);
});
