import ThreadModel from '~/models/thread';
// import { ParamsError } from '~/utils/error';
import { query } from '~/utils/pg';

import startPg, { migrate } from '../__utils__/pgServer';
import mockContext from '../__utils__/context';

// use thread model to test query slice

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

describe('query thread', () => {
  const threadTags = [
    { mainTag: 'MainA', subTags: ['SubA', 'SubB', 'SubC'] }, // <- oldest, 0, end
    { mainTag: 'MainA', subTags: ['SubA', 'SubB', 'SubD'] }, // 1
    { mainTag: 'MainA', subTags: ['SubA', 'SubB', 'SubE'] }, // 2
    { mainTag: 'MainA', subTags: ['SubA', 'SubC', 'SubD'] }, // 3
    { mainTag: 'MainA', subTags: ['SubA', 'SubC', 'SubE'] }, // 4
    { mainTag: 'MainA', subTags: ['SubA', 'SubD', 'SubE'] }, // 5
    { mainTag: 'MainA', subTags: ['SubB', 'SubC', 'SubD'] }, // 6
    { mainTag: 'MainA', subTags: ['SubB', 'SubC', 'SubE'] }, // 7
    { mainTag: 'MainA', subTags: ['SubB', 'SubD', 'SubE'] }, // 8
    { mainTag: 'MainA', subTags: ['SubC', 'SubD', 'SubE'] }, // 9
    { mainTag: 'MainB', subTags: ['SubX'] }, // <- newest, 10, begin
  ];
  const threadIds = [];
  it('parpare data', async () => {
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['MainA', true]);
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['MainB', true]);
    const ctx = await mockContext({ email: 'test@uexky.com' });
    for (let i = 0; i < threadTags.length; i += 1) {
    /* eslint-disable no-await-in-loop */
      const thread = await ThreadModel.new({
        /* eslint-enable no-await-in-loop */
        ctx,
        thread: {
          anonymous: true,
          content: 'test content',
          mainTag: threadTags[i].mainTag,
          subTags: threadTags[i].subTags,
          title: '',
        },
      });
      threadIds.push(thread.id);
    }
  });
  it('after beginning', async () => {
    const { threads, sliceInfo } = await ThreadModel.findSlice({
      tags: ['MainA', 'MainB'],
      query: { after: '', limit: 10 },
    });
    expect(threads.length).toEqual(10);
    expect(threads[0].id.duid).toEqual(threadIds[10].duid);
    expect(threads[9].id.duid).toEqual(threadIds[1].duid);
    expect(sliceInfo.firstCursor).toEqual(threadIds[10].duid);
    expect(sliceInfo.lastCursor).toEqual(threadIds[1].duid);
    expect(sliceInfo.hasNext).toBeTruthy();
  });
  it('before end, filter MainA', async () => {
    const { threads, sliceInfo } = await ThreadModel.findSlice({
      tags: ['MainA'],
      query: { before: '', limit: 10 },
    });
    expect(threads.length).toEqual(10);
    expect(sliceInfo.firstCursor).toEqual(threadIds[9].duid);
    expect(sliceInfo.lastCursor).toEqual(threadIds[0].duid);
    expect(sliceInfo.hasNext).toBeFalsy();
  });
  it('after cursor', async () => {
    const { threads, sliceInfo } = await ThreadModel.findSlice({
      tags: ['SubA', 'SubB'],
      query: { after: threadIds[7].duid, limit: 5 },
    });
    expect(threads.length).toEqual(5);
    expect(sliceInfo.firstCursor).toEqual(threadIds[6].duid);
    expect(threads[1].id.duid).toEqual(threadIds[5].duid);
    expect(threads[2].id.duid).toEqual(threadIds[4].duid);
    expect(threads[3].id.duid).toEqual(threadIds[3].duid);
    expect(sliceInfo.lastCursor).toEqual(threadIds[2].duid);
    expect(sliceInfo.hasNext).toBeTruthy();
  });
  it('before cursor', async () => {
    const { threads, sliceInfo } = await ThreadModel.findSlice({
      tags: ['SubC'],
      query: { before: threadIds[1].duid, limit: 3 },
    });
    expect(threads.length).toEqual(3);
    expect(sliceInfo.firstCursor).toEqual(threadIds[6].duid);
    expect(threads[1].id.duid).toEqual(threadIds[4].duid);
    expect(sliceInfo.lastCursor).toEqual(threadIds[3].duid);
    expect(sliceInfo.hasNext).toBeTruthy();
  });
  it('no result, no after', async () => {
    const { threads, sliceInfo } = await ThreadModel.findSlice({
      tags: ['SubD'],
      query: { after: threadIds[1].duid, limit: 10 },
    });
    expect(threads.length).toEqual(0);
    expect(sliceInfo.firstCursor).toEqual('');
    expect(sliceInfo.lastCursor).toEqual('');
    expect(sliceInfo.hasNext).toBeFalsy();
  });
});
