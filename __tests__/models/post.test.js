import { ROLE } from '~/models/user';
import ThreadModel from '~/models/thread';
import PostModel, { blockedContent } from '~/models/post';
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

describe('publish post and read', () => {
  const mockUser = { email: 'test@uexky.com', name: 'test user' };
  let ctx;
  let thread;
  const postIds = [];
  it('parpare data', async () => {
    await query('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['MainA', true]);
    ctx = await mockContext(mockUser);
    thread = await ThreadModel.new({
      ctx,
      thread: {
        anonymous: true,
        content: 'test thread content',
        mainTag: 'MainA',
        subTags: [],
        title: 'test thread title',
      },
    });
    const replyCount = await PostModel.getThreadReplyCount({ threadId: thread.id });
    expect(replyCount).toEqual(0);
    const catalog = await PostModel.getThreadCatalog({ threadId: thread.id });
    expect(catalog.length).toEqual(0);
  });
  it('new post', async () => {
    const input = {
      threadId: thread.id,
      anonymous: true,
      content: 'test post1',
      quoteIds: [],
    };
    const before = new Date();
    const post = await PostModel.new({ ctx, post: input });
    const after = new Date();
    expect(post.id.type).toEqual('UID');
    expect(post.createdAt.getTime()).toBeGreaterThan(before.getTime() - 500);
    expect(post.createdAt.getTime()).toBeLessThan(after.getTime() + 500);
    expect(post.updatedAt.getTime()).toBeGreaterThan(before.getTime() - 500);
    expect(post.updatedAt.getTime()).toBeLessThan(after.getTime() + 500);
    expect(post.anonymous).toBeTruthy();
    expect(post.author).toMatch(/[0-9a-zA-Z-_]{6,}/);
    expect(post.content).toEqual(input.content);
    expect(post.blocked).toBeFalsy();
    postIds.push(post.id);
  });
  it('new post with name', async () => {
    const input = {
      threadId: thread.id,
      anonymous: false,
      content: 'test post2',
      quoteIds: [],
    };
    const post = await PostModel.new({ ctx, post: input });
    expect(post.author).toEqual(mockUser.name);
    postIds.push(post.id);
  });
  it('new post with quotes', async () => {
    const input = {
      threadId: thread.id,
      anonymous: false,
      content: 'test post2',
      quoteIds: postIds.map(pid => pid.suid),
    };
    const post = await PostModel.new({ ctx, post: input });
    const quotedPosts = await post.getQuotes();
    expect(quotedPosts.length).toEqual(2);
    const qduids = quotedPosts.map(qp => qp.id.suid).sort();
    expect(qduids[0].toString()).toEqual(postIds[0].suid.toString());
    expect(qduids[1].toString()).toEqual(postIds[1].suid.toString());
    postIds.push(post.id);
  });
  it('quoted count', async () => {
    const post0 = await PostModel.findById({ postId: postIds[0] });
    const post1 = await PostModel.findById({ postId: postIds[1] });
    const post0qc = await post0.getQuotedCount();
    const post1qc = await post1.getQuotedCount();
    expect(post0qc).toEqual(1);
    expect(post1qc).toEqual(1);
  });
  it('find thread posts', async () => {
    const { posts, sliceInfo } = await PostModel.findThreadPosts({
      threadId: thread.id,
      query: { after: '', limit: 10 },
    });
    expect(posts.length).toEqual(3);
    expect(sliceInfo.firstCursor).toEqual(postIds[0].duid);
    expect(sliceInfo.lastCursor).toEqual(postIds[2].duid);
  });
  it('find user posts', async () => {
    const user = ctx.auth.signedInUser();
    const { posts, sliceInfo } = await PostModel.findUserPosts({
      user,
      query: { after: '', limit: 10 },
    });
    expect(posts.length).toEqual(3);
    expect(sliceInfo.firstCursor).toEqual(postIds[0].duid);
    expect(sliceInfo.lastCursor).toEqual(postIds[2].duid);
  });
  it('reply count', async () => {
    const replyCount = await PostModel.getThreadReplyCount({ threadId: thread.id });
    expect(replyCount).toEqual(3);
  });
  it('thread catalog', async () => {
    const catalog = await PostModel.getThreadCatalog({ threadId: thread.id });
    expect(catalog.length).toEqual(3);
    expect(catalog[0].postId).toEqual(postIds[0].duid);
    expect(catalog[1].postId).toEqual(postIds[1].duid);
    expect(catalog[2].postId).toEqual(postIds[2].duid);
  });
  it('block thread', async () => {
    const modContext = await mockContext({
      email: 'mod@uexky.com',
      role: ROLE.MOD,
    });
    await PostModel.block({ ctx: modContext, postId: postIds[1] });
    const post = await PostModel.findById({ postId: postIds[1] });
    expect(post.blocked).toBeTruthy();
    expect(post.content).toEqual(blockedContent);
  });
});
