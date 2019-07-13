import gql from 'graphql-tag';
import {
  mockUser, mockAltUser, query, mutate, altMutate,
} from '../__utils__/apolloClient';

import AidModel from '~/models/aid';
import UserModel from '~/models/user';
import { query as pgq } from '~/utils/pg';

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

const mockEmail = mockUser.email;
const mockReplierEmail = mockAltUser.email;

const mockThread = {
  anonymous: true,
  content: 'Test Content',
  mainTag: 'MainA',
  subTags: ['SubA'],
  title: 'TestTitle',
};
const mockPost = {
  threadId: '',
  anonymous: true,
  content: 'Test Reply',
};
const mockReply = {
  threadId: '',
  anonymous: true,
  content: 'Test Reply',
  quoteIds: [],
};
let threadId;
let postId;
let replyId;

const PUB_THREAD = gql`
  mutation PubThread($thread: ThreadInput!) {
    pubThread(thread: $thread) {
      id, anonymous, author, content, createdAt, mainTag, subTags, title, replyCount, locked, blocked
    }
  }
`;

const PUB_POST = gql`
  mutation PubPost($post: PostInput!) {
    pubPost(post: $post) {
      id, anonymous, author, content, createdAt, quotes { id }, quotedCount, blocked
    }
  }
`;

const GET_POST = gql`
  query Post($id: String!) {
    post(id: $id) {
      id, anonymous, author, content, createdAt, quotes { id }, quotedCount, blocked
    }
  }
`;

const GET_NOTI = gql`
  query GetNoti($type: String!, $query: SliceQuery!) {
    unreadNotiCount {
      system, replied, quoted
    }
    notification(type: $type, query: $query) {
      replied { id, type, eventTime, hasRead, thread { id }, repliers  }
      quoted { id, type, eventTime, hasRead, thread { id }, quotedPost { id }, post { id } }
      sliceInfo { firstCursor, lastCursor, hasNext }
    }
  }
`;

describe('Testing posting a thread', () => {
  it('Setting mainTag', async () => {
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['MainA', true]);
    await mockContext({ email: mockEmail });
    await mockContext({ email: mockReplierEmail });
  });
  it('Posting thread', async () => {
    const { data, errors } = await mutate({
      mutation: PUB_THREAD, variables: { thread: mockThread },
    });
    expect(errors).toBeUndefined();

    const result = data.pubThread;
    expect(JSON.stringify(result.subTags)).toEqual(JSON.stringify(mockThread.subTags));
    expect(result.anonymous).toEqual(true);
    expect(result.content).toEqual(mockThread.content);
    expect(result.title).toEqual(mockThread.title);
    expect(result.blocked).toEqual(false);
    expect(result.locked).toEqual(false);

    threadId = result.id;
    mockPost.threadId = result.id;
    mockReply.threadId = result.id;
  });
});

describe('Testing replying a thread', () => {
  it('Posting reply', async () => {
    const user = await UserModel.findByEmail({ email: mockEmail });
    const newPost = mockPost;
    const { data, errors } = await mutate({
      mutation: PUB_POST, variables: { post: newPost },
    });
    expect(errors).toBeUndefined();

    const postResult = data.pubPost;
    const aid = await AidModel.getAid({ userId: user.id, threadId });
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toEqual([]);
    expect(postResult.quotedCount).toEqual(0);
    expect(postResult.author).toEqual(aid.duid);
    expect(postResult.blocked).toEqual(false);

    postId = postResult.id;
    mockReply.quoteIds.push(postResult.id);
  });
  it('Posting reply with quotes', async () => {
    const replier = await UserModel.findByEmail({ email: mockReplierEmail });
    const newPost = mockReply;
    const { data, errors } = await altMutate({
      mutation: PUB_POST, variables: { post: newPost },
    });
    expect(errors).toBeUndefined();

    const postResult = data.pubPost;
    const aid = await AidModel.getAid({ userId: replier.id, threadId });
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toEqual([{ id: postId }]);
    expect(postResult.author).toEqual(aid.duid);
    expect(postResult.blocked).toEqual(false);

    replyId = postResult.id;
  });
  it('Validate updated post', async () => {
    const user = await UserModel.findByEmail({ email: mockEmail });
    const { data, errors } = await query({ query: GET_POST, variables: { id: postId } });
    expect(errors).toBeUndefined();

    const postResult = data.post;
    const aid = await AidModel.getAid({ userId: user.id, threadId });
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toEqual([]);
    expect(postResult.quotedCount).toEqual(1);
    expect(postResult.author).toEqual(aid.duid);
    expect(postResult.blocked).toEqual(false);
  });
  it('Validate reply', async () => {
    const replier = await UserModel.findByEmail({ email: mockReplierEmail });
    const { data, errors } = await query({ query: GET_POST, variables: { id: replyId } });
    expect(errors).toBeUndefined();

    const postResult = data.post;
    const aid = await AidModel.getAid({ userId: replier.id, threadId });
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toEqual([{ id: postId }]);
    expect(postResult.author).toEqual(aid.duid);
    expect(postResult.blocked).toEqual(false);
  });
  it('Validating post for mockUser in NotificationModel', async () => {
    const user = await UserModel.findByEmail({ email: mockEmail });
    const replier = await UserModel.findByEmail({ email: mockReplierEmail });
    const aid = await AidModel.getAid({ userId: user.id, threadId });
    const replyAid = await AidModel.getAid({ userId: replier.id, threadId });

    const repResult = await query({
      query: GET_NOTI,
      variables: { type: 'replied', query: { after: '', limit: 10 } },
    });
    const { data: repData, errors: repErrors } = repResult;
    expect(repErrors).toBeUndefined();
    expect(repData.unreadNotiCount).toEqual({ system: 0, replied: 1, quoted: 1 });
    expect(repData.notification.quoted).toBeNull();
    expect(repData.notification.replied[0].id).toEqual(`replied:${threadId}`);
    expect(repData.notification.replied[0].type).toEqual('replied');
    expect(repData.notification.replied[0].hasRead).toEqual(false);
    expect(repData.notification.replied[0].thread.id).toEqual(threadId);
    expect(repData.notification.replied[0].repliers).toEqual([aid.duid, replyAid.duid]);

    const quoResult = await query({
      query: GET_NOTI,
      variables: { type: 'quoted', query: { after: '', limit: 10 } },
    });
    const { data: quoData, errors: quoErrors } = quoResult;
    expect(quoErrors).toBeUndefined();
    expect(quoData.unreadNotiCount).toEqual({ system: 0, replied: 0, quoted: 1 });
    expect(quoData.notification.replied).toBeNull();
    expect(quoData.notification.quoted[0].id).toEqual(`quoted:${postId}:${replyId}`);
    expect(quoData.notification.quoted[0].type).toEqual('quoted');
    expect(quoData.notification.quoted[0].hasRead).toEqual(false);
    expect(quoData.notification.quoted[0].thread.id).toEqual(threadId);
    expect(quoData.notification.quoted[0].post.id).toEqual(replyId);
    expect(quoData.notification.quoted[0].quotedPost.id).toEqual(postId);
  });
});
