import startPg, { migrate } from '../__utils__/pgServer';
import newSuid from '~/uid/generator';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

const analyzeSuid = (suid) => {
  let n = suid;
  const random = n % BigInt(2 ** 9);
  n /= BigInt(2 ** 9);
  const sequence = n % BigInt(2 ** 10);
  n /= BigInt(2 ** 10);
  const workerID = n % BigInt(2 ** 9);
  n /= BigInt(2 ** 9);
  const timestamp = n;
  return {
    timestamp, workerID, sequence, random,
  };
};

test('generator.newID', async () => {
  const id1 = await newSuid();
  expect(typeof id1).toEqual('bigint');
  const id2 = await newSuid();
  expect(id2 > id1).toBeTruthy();

  const id1s = analyzeSuid(id1);
  const id2s = analyzeSuid(id2);
  expect(id2s.timestamp > id1s.timestamp || id2s.timestamp === id1s.timestamp).toBeTruthy();
  expect(id2s.workerID === id1s.workerID).toBeTruthy();
  expect(id2s.sequence > id1s.sequence).toBeTruthy();
});
