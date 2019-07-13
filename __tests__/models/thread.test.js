import PostModel from '~/models/post';
import ThreadModel, { blockedContent } from '~/models/thread';
import { ROLE } from '~/models/user';
import { ParamsError } from '~/utils/error';
import { query } from '~/utils/pg';
import startPg, { migrate } from '../__utils__/pgServer';
import mockContext from '../__utils__/context';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

describe('Testing posting a thread', () => {
  const mockUser = {
    email: 'test@example.com',
    name: 'testUser',
  };
  const mockInput = {
    anonymous: true,
    content: 'Test Content',
    mainTag: 'MainA',
    subTags: ['SubA'],
    title: 'TestTitle',
  };
  const threadIds = [];
  it('parpare data', async () => {
    await query('INSERT INTO tag (name, "isMain") VALUES ($1, $2)', ['MainA', true]);
  });
  it('publish a thread anonymous', async () => {
    const ctx = await mockContext({ email: mockUser.email, name: mockUser.name });
    const newThread = mockInput;
    const before = new Date();
    const thread = await ThreadModel.new({ ctx, thread: newThread });
    const after = new Date();

    expect(thread.id.type).toEqual('UID');
    expect(thread.createdAt.getTime()).toBeGreaterThanOrEqual(before.getTime());
    expect(thread.createdAt.getTime()).toBeLessThanOrEqual(after.getTime());
    expect(thread.updatedAt.getTime()).toBeGreaterThanOrEqual(before.getTime());
    expect(thread.updatedAt.getTime()).toBeLessThanOrEqual(after.getTime());
    expect(thread.anonymous).toBeTruthy();
    expect(thread.author).toMatch(/[0-9a-zA-Z-_]{6,}/);
    expect(thread.title).toEqual(mockInput.title);
    expect(thread.content).toEqual(mockInput.content);
    expect(thread.blocked).toBeFalsy();
    expect(thread.locked).toBeFalsy();

    const mainTag = await thread.getMainTag();
    expect(mainTag).toEqual('MainA');
    const subTags = await thread.getSubTags();
    expect(subTags.length).toEqual(1);
    expect(subTags).toContain('SubA');

    threadIds.push(thread.id);
  });
  it('publish a thread not anonymous', async () => {
    const ctx = await mockContext({ email: mockUser.email, name: mockUser.name });
    const input = { ...mockInput, anonymous: false };
    const thread = await ThreadModel.new({ ctx, thread: input });
    expect(thread.anonymous).toBeFalsy();
    expect(thread.author).toEqual(mockUser.name);
    threadIds.push(thread.id);
  });
  it('check tag belongsTo', async () => {
    const { rows } = await query(
      'SELECT * FROM tags_main_tags WHERE name=$1', [mockInput.subTags[0]],
    );
    expect(rows[0].belongsTo).toEqual(mockInput.mainTag);
  });
  it('publish thread without subTags', async () => {
    const ctx = await mockContext({ email: mockUser.email, name: mockUser.name });
    const input = { ...mockInput, subTags: [] };
    const thread = await ThreadModel.new({ ctx, thread: input });
    const subTags = await thread.getSubTags();
    expect((subTags || []).length).toEqual(0);
    threadIds.push(thread.id);
  });
  it('findSlice by tag', async () => {
    const { threads, sliceInfo } = await ThreadModel.findSlice({
      tags: ['SubA'],
      query: { after: '', limit: 5 },
    });
    expect(threads.length).toEqual(2);
    expect(sliceInfo.firstCursor).toEqual(threadIds[1].duid);
    expect(sliceInfo.lastCursor).toEqual(threadIds[0].duid);
    expect(sliceInfo.hasNext).toBeFalsy();
  });
  it('findById', async () => {
    const thread = await ThreadModel.findById({ threadId: threadIds[0] });
    expect(thread.id.duid).toEqual(threadIds[0].duid);
  });
  it('findUserSlice', async () => {
    const ctx = await mockContext({ email: mockUser.email, name: mockUser.name });
    const user = ctx.auth.signedInUser();
    const { threads, sliceInfo } = await ThreadModel.findUserThreads({
      user,
      query: { before: threadIds[1].duid, limit: 10 },
    });
    expect(threads.length).toEqual(1);
    expect(sliceInfo.firstCursor).toEqual(threadIds[2].duid);
  });
  it('lock thread', async () => {
    const modContext = await mockContext({
      email: 'mod@uexky.com',
      role: ROLE.MOD,
    });
    const ctx = await mockContext({ email: mockUser.email, name: mockUser.name });
    await ThreadModel.lock({ ctx: modContext, threadId: threadIds[0] });
    const thread = await ThreadModel.findById({ threadId: threadIds[0] });
    expect(thread.locked).toBeTruthy();
    await expect(PostModel.new({
      ctx,
      post: {
        threadId: threadIds[0],
        anonymous: true,
        content: 'try reply locked thread',
        quoteIds: [],
      },
    })).rejects.toThrow(ParamsError);
  });
  it('block thread', async () => {
    const modContext = await mockContext({
      email: 'mod@uexky.com',
      role: ROLE.MOD,
    });
    await ThreadModel.block({ ctx: modContext, threadId: threadIds[1] });
    const thread = await ThreadModel.findById({ threadId: threadIds[1] });
    expect(thread.blocked).toBeTruthy();
    expect(thread.content).toEqual(blockedContent);
  });
});
