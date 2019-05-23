import WorkerIdModel from '~/models/workerId';
import startPg, { migrate } from '../__utils__/pgServer';

let pgPool;
// let db;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

test('get worker id', async () => {
  const ids = await Promise.all([1, 2, 3, 4, 5].map(() => WorkerIdModel().newWorkerId()));
  expect(ids.sort()).toEqual([1, 2, 3, 4, 5]);
});
