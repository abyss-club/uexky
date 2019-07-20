import NotificationModel, { USER_GROUPS } from '~/models/notification';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';
import { query } from '~/utils/pg';
import { suidRegex, duidRegex } from '~/uid';
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

describe('test notification', () => {
  const mockUser = { email: 'test@uexky.com', name: 'test user' };
  let ctx;
  let thread;
  let post1;
  let post2;
  it('prepare data', async () => {
    await query(
      'INSERT INTO tag (name, is_main) VALUES ($1, $2)',
      ['MainA', true],
    );
    ctx = await mockContext(mockUser);
    const unread = {
      system: await NotificationModel.getUnreadCount({ ctx, type: 'system' }),
      replied: await NotificationModel.getUnreadCount({ ctx, type: 'replied' }),
      quoted: await NotificationModel.getUnreadCount({ ctx, type: 'quoted' }),
    };
    expect(unread).toEqual({ system: 0, replied: 0, quoted: 0 });
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
    post1 = await PostModel.new({
      ctx,
      post: {
        threadId: thread.id,
        anonymous: true,
        content: 'test post1',
        quoteIds: [],
      },
    });
    post2 = await PostModel.new({
      ctx,
      post: {
        threadId: thread.id,
        anonymous: false,
        content: 'test post2',
        quoteIds: [post1.id],
      },
    });
  });
  it('replied notification', async () => {
    const count = await NotificationModel.getUnreadCount({
      ctx, type: 'replied',
    });
    expect(count).toEqual(1);
    const { replied, sliceInfo } = await NotificationModel.findNotiSlice({
      ctx, type: 'replied', query: { after: '', limit: 5 },
    });
    expect(sliceInfo.firstCursor).toMatch(duidRegex);
    expect(sliceInfo.lastCursor).toMatch(duidRegex);
    expect(replied.length).toEqual(1);
    expect(replied[0].id.toString()).toEqual(thread.id.suid.toString());
    expect(replied[0].key).toEqual(`replied:${thread.id.duid}`);
    expect(replied[0].type).toEqual('replied');
    expect(replied[0].hasRead).toBeFalsy();
    expect(replied[0].threadId.duid).toEqual(thread.id.duid);
  });
  it('quoted notification', async () => {
    const count = await NotificationModel.getUnreadCount({
      ctx, type: 'quoted',
    });
    expect(count).toEqual(1);
    const { quoted, sliceInfo } = await NotificationModel.findNotiSlice({
      ctx, type: 'quoted', query: { after: '', limit: 5 },
    });
    expect(sliceInfo.firstCursor).toMatch(duidRegex);
    expect(sliceInfo.lastCursor).toMatch(duidRegex);
    expect(quoted.length).toEqual(1);
    expect(quoted[0].id).toMatch(suidRegex);
    expect(quoted[0].key).toEqual(
      `quoted:${post1.id.duid}:${post2.id.duid}`,
    );
    expect(quoted[0].type).toEqual('quoted');
    expect(quoted[0].hasRead).toBeFalsy();
    expect(quoted[0].threadId.duid).toEqual(thread.id.duid);
    expect(quoted[0].quotedId.duid).toEqual(post1.id.duid);
    expect(quoted[0].postId.duid).toEqual(post2.id.duid);
  });
  it('system notification', async () => {
    ctx = await mockContext(mockUser); // refresh user
    const user = ctx.auth.signedInUser();
    const noti1 = {
      sendToGroup: USER_GROUPS.ALL_USER,
      title: 'to all',
      content: 'to all',
    };
    const noti2 = {
      sendTo: user.id,
      title: 'hello',
      content: 'system notification',
    };
    await NotificationModel.newSystemNoti(noti1);
    await NotificationModel.newSystemNoti(noti2);
    const unread = {
      system: await NotificationModel.getUnreadCount({ ctx, type: 'system' }),
      replied: await NotificationModel.getUnreadCount({ ctx, type: 'replied' }),
      quoted: await NotificationModel.getUnreadCount({ ctx, type: 'quoted' }),
    };
    expect(unread).toEqual({ system: 2, replied: 0, quoted: 0 });

    const { system, sliceInfo } = await NotificationModel.findNotiSlice({
      ctx, type: 'system', query: { after: '', limit: 5 },
    });
    expect(sliceInfo.firstCursor).toMatch(duidRegex);
    expect(sliceInfo.lastCursor).toMatch(duidRegex);
    expect(system.length).toEqual(2);
    expect(system[0].id).toMatch(suidRegex);
    expect(system[0].key).toMatch(/system:[0-9a-zA-Z-_]{6,}/);
    expect(system[0].hasRead).toBeFalsy();
    expect(system[0].title).toEqual(noti2.title);
    expect(system[0].content).toEqual(noti2.content);
    expect(system[1].id).toMatch(suidRegex);
    expect(system[1].key).toMatch(/system:[0-9a-zA-Z-_]{6,}/);
    expect(system[1].hasRead).toBeFalsy();
    expect(system[1].title).toEqual(noti1.title);
    expect(system[1].content).toEqual(noti1.content);
  });
});
