import gql from 'graphql-tag';
import { db } from '~/utils/mongo';

import startRepl from '../__utils__/mongoServer';
import {
  mockUser, mockAltUser, query, mutate, altMutate,
} from '../__utils__/apolloClient';

import Uid from '~/uid';
import UserModel from '~/models/user';
import TagModel from '~/models/tag';

jest.setTimeout(60000);

let replSet;
let mongoClient;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

const USERPOSTS = 'userPosts';
const NOTIFICATION = 'notification';

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
      id, anonymous, author, content, createdAt, quotes { id }, quoteCount, blocked
    }
  }
`;

const GET_POST = gql`
  query Post($id: String!) {
    post(id: $id) {
      id, anonymous, author, content, createdAt, quotes { id }, quoteCount, blocked
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
      quoted { id, type, eventTime, hasRead, thread { id }, quotedPost { id }, post { id }, quoter }
      sliceInfo { firstCursor, lastCursor, hasNext }
    }
  }
`;

describe('Testing posting a thread', () => {
  it('Create collections', async () => {
    await db.createCollection(USERPOSTS);
  });
  it('Setting mainTag', async () => {
    await TagModel().addMainTag('MainA');
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
  it('Create collections', async () => {
    await db.createCollection(NOTIFICATION);
  });
  it('Posting reply', async () => {
    const user = await UserModel().getUserByEmail(mockEmail);
    const newPost = mockPost;
    const { data, errors } = await mutate({
      mutation: PUB_POST, variables: { post: newPost },
    });
    expect(errors).toBeUndefined();

    const postResult = data.pubPost;
    const author = await UserModel().methods(user).author(Uid.encode(threadId), true);
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toEqual([]);
    expect(postResult.quoteCount).toEqual(0);
    expect(postResult.author).toEqual(author);
    expect(postResult.blocked).toEqual(false);

    postId = postResult.id;
    mockReply.quoteIds.push(postResult.id);
  });
  it('Posting reply with quotes', async () => {
    const replier = await UserModel().getUserByEmail(mockReplierEmail);
    const newPost = mockReply;
    const { data, errors } = await altMutate({
      mutation: PUB_POST, variables: { post: newPost },
    });
    expect(errors).toBeUndefined();

    const postResult = data.pubPost;
    const author = await UserModel().methods(replier).author(Uid.encode(threadId), true);
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toEqual([{ id: postId }]);
    expect(postResult.author).toEqual(author);
    expect(postResult.blocked).toEqual(false);

    replyId = postResult.id;
  });
  it('Validate updated post', async () => {
    const user = await UserModel().getUserByEmail(mockEmail);
    const { data, errors } = await query({ query: GET_POST, variables: { id: postId } });
    expect(errors).toBeUndefined();

    const postResult = data.post;
    const author = await UserModel().methods(user).author(Uid.encode(threadId), true);
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toEqual([]);
    expect(postResult.quoteCount).toEqual(1);
    expect(postResult.author).toEqual(author);
    expect(postResult.blocked).toEqual(false);
  });
  it('Validate reply', async () => {
    const replier = await UserModel().getUserByEmail(mockReplierEmail);
    const { data, errors } = await query({ query: GET_POST, variables: { id: replyId } });
    expect(errors).toBeUndefined();

    const postResult = data.post;
    const author = await UserModel().methods(replier).author(Uid.encode(threadId), true);
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toEqual([{ id: postId }]);
    expect(postResult.author).toEqual(author);
    expect(postResult.blocked).toEqual(false);
  });
  it('Validating post for mockUser in NotificationModel', async () => {
    const user = await UserModel().getUserByEmail(mockEmail);
    const replier = await UserModel().getUserByEmail(mockReplierEmail);
    const userAuthor = await UserModel().methods(user).author(Uid.encode(threadId), true);
    const replyAuthor = await UserModel().methods(replier).author(Uid.encode(threadId), true);


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
    expect(repData.notification.replied[0].repliers).toEqual([userAuthor, replyAuthor]);

    const quoResult = await query({
      query: GET_NOTI,
      variables: { type: 'quoted', query: { after: '', limit: 10 } },
    });
    const { data: quoData, errors: quoErrors } = quoResult;
    expect(quoErrors).toBeUndefined();
    expect(quoData.unreadNotiCount).toEqual({ system: 0, replied: 0, quoted: 1 });
    expect(quoData.notification.replied).toBeNull();
    expect(quoData.notification.quoted[0].id).toEqual(`quoted:${replyId}:${postId}`);
    expect(quoData.notification.quoted[0].type).toEqual('quoted');
    expect(quoData.notification.quoted[0].hasRead).toEqual(false);
    expect(quoData.notification.quoted[0].thread.id).toEqual(threadId);
    expect(quoData.notification.quoted[0].post.id).toEqual(postId);
    expect(quoData.notification.quoted[0].quoter).toEqual(replyAuthor);
    expect(quoData.notification.quoted[0].quotedPost.id).toEqual(postId);
  });
});
