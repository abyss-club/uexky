import UID from '~/uid';
import AidModel from '~/models/aid';

import startPg, { migrate } from '../__utils__/pgServer';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

describe('test aid', () => {
  const userId = 1;
  let threadId;
  let aid;
  it('parpare date', async () => {
    threadId = await UID.new();
  });
  it('new aid', async () => {
    aid = await AidModel.getAid({ userId, threadId });
    expect(aid.type).toEqual('UID');
  });
  it('already have aid', async () => {
    const id = await AidModel.getAid({ userId, threadId });
    expect(id.type).toEqual('UID');
    expect(id.duid).toEqual(aid.duid);
    expect(id.suid.toString()).toEqual(aid.suid.toString());
  });
});
