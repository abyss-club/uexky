import mongoose from 'mongoose';
import { startRepl } from '../__utils__/mongoServer';
import WorkerIDModel from '~/models/workId';

jest.setTimeout(60000); // for boot replica sets
let mongoServer;

beforeAll(async () => {
  mongoServer = await startRepl();
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
