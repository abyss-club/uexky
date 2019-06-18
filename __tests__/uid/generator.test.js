import startPg, { migrate } from '../__utils__/pgServer';
import generator from '~/uid/generator';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

const getTimestamp = uid => parseInt(uid.substring(0, 8), 16);
const getWorkerId = (uid) => {
  const slice = parseInt(uid.substring(8, 11), 16);
  return Math.floor(slice / 8);
};
const getSeq = (uid) => {
  const slice = parseInt(uid.substring(10, 13), 16);
  return Math.floor(slice / 2);
};

test('generator.newID', async () => {
  const id1 = await generator.newID();
  expect(id1).toMatch(/[0-9a-f]{15}/);
  const id2 = await generator.newID();
  expect(id2).toMatch(/[0-9a-f]{15}/);
  expect(id2 > id1).toBeTruthy();
  expect(getTimestamp(id2)).toBeGreaterThanOrEqual(getTimestamp(id1));
  expect(getWorkerId(id2)).toEqual(getWorkerId(id1));
  expect(getSeq(id2)).toBeGreaterThan(getSeq(id1));
});
