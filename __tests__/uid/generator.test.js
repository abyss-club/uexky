import mongoose from 'mongoose';
import { startMongo } from '../__utils__/mongoServer';
import generator, { WorkerIDModel } from '~/uid/generator';

let mongoServer;

beforeAll(async () => {
  mongoServer = await startMongo();
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});

test('get worker id', async () => {
  const ids = await Promise.all([1, 2, 3, 4, 5].map(
    WorkerIDModel.newWorkerID,
  ));
  expect(ids.sort()).toEqual([1, 2, 3, 4, 5]);
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
