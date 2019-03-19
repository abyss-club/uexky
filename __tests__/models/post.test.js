import startRepl from '../__utils__/mongoServer';
import mongo, { db } from '~/utils/mongo';

import Uid from '~/uid';
import UserModel from '~/models/user';
// import UserPostsModel from '~/models/userPosts';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';
import TagModel from '~/models/tag';
import NotificationModel from '~/models/notification';

jest.setTimeout(60000); // for boot replica sets
let replSet;
let mongoClient;
// let db;

beforeAll(async () => {
  ({ replSet, mongoClient } = await startRepl());
});

afterAll(() => {
  mongoClient.close();
  replSet.stop();
});

const THREAD = 'thread';
const POST = 'post';
const USERPOSTS = 'userPosts';
const NOTIFICATION = 'notification';
const col = () => mongo.collection(POST);

const mockUser = {
  email: 'test@example.com',
  name: 'testUser',
};
const mockReplier = {
  email: 'replier@example.com',
  name: 'testReplier',
};

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
const mockTagTree = {
  mainTag: 'MainA', subTags: ['SubA'],
};

let threadSuid;
let threadId;
let postSuid;
let postId;
let replyId;
let replySuid;

describe('Testing posting a thread', () => {
  it('Create collections', async () => {
    await db.createCollection(USERPOSTS);
    await db.createCollection(NOTIFICATION);
  });
  it('Setting tags', async () => {
    await TagModel().addMainTag('MainA');
  });
  it('Posting a thread', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    const newThread = mockThread;
    const { _id } = await ThreadModel({ user }).pubThread({ user }, newThread);
    const threadResult = await mongo.collection(THREAD).findOne({ _id });
    const author = await UserModel({ user }).methods(user).author(threadResult.suid, true);
    expect(JSON.stringify(threadResult.subTags)).toEqual(JSON.stringify(newThread.subTags));
    expect(JSON.stringify(threadResult.tags))
      .toEqual(JSON.stringify([newThread.mainTag, ...newThread.subTags]));
    expect(threadResult.anonymous).toEqual(true);
    expect(threadResult.content).toEqual(newThread.content);
    expect(threadResult.title).toEqual(newThread.title);
    expect(threadResult.userId).toEqual(user._id);
    expect(threadResult.author).toEqual(author);
    expect(threadResult.blocked).toEqual(false);
    expect(threadResult.locked).toEqual(false);

    threadSuid = threadResult.suid;
    threadId = Uid.decode(threadSuid);
    mockPost.threadId = threadId;
    mockReply.threadId = threadId;
  });
  it('Validating the thread in UserPostsModel', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    const result = await mongo.collection(USERPOSTS).find(
      { userId: user._id, threadSuid },
    ).toArray();
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating the thread in TagModel', async () => {
    const result = await TagModel().getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockTagTree));
  });
});

describe('Testing replying a thread', () => {
  it('Create collections', async () => {
    await db.createCollection(NOTIFICATION);
  });
  it('Posting reply', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    const newPost = mockPost;
    const { _id } = await PostModel({ user }).pubPost({ user }, newPost);
    const postResult = await mongo.collection(POST).findOne({ _id });

    postSuid = postResult.suid;
    postId = Uid.decode(postSuid);

    const author = await UserModel({ user }).methods(user).author(threadSuid, true);
    expect(postResult.threadSuid).toEqual(threadSuid);
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(newPost.content);
    expect(postResult.quotes).toBeUndefined();
    expect(postResult.userId).toEqual(user._id);
    expect(postResult.author).toEqual(author);
    expect(postResult.blocked).toEqual(false);

    mockReply.quoteIds.push(postId);
  });
  it('Posting reply with quotes', async () => {
    const user = await UserModel().getUserByEmail(mockReplier.email);
    const newPost = mockReply;
    const { _id } = await PostModel({ user }).pubPost({ user }, newPost);
    const postResult = await mongo.collection(POST).findOne({ _id });

    replySuid = postResult.suid;
    replyId = Uid.decode(replySuid);

    const author = await UserModel({ user }).methods(user).author(threadSuid, true);
    expect(postResult.threadSuid).toEqual(threadSuid);
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toBeUndefined();
    expect(postResult.userId).toEqual(user._id);
    expect(postResult.author).toEqual(author);
    expect(postResult.blocked).toEqual(false);
  });
  it('Validate updated reply', async () => {
    const postResult = await col().findOne({ suid: replySuid });
    expect(postResult.quoteSuids[0]).toEqual(postSuid);
  });
  it('Validating post for mockUser in UserPostsModel', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    const result = await mongo.collection(USERPOSTS).find(
      { userId: user._id, threadSuid },
    ).toArray();
    expect(result.length).toEqual(1);
    expect(result[0].posts[0].suid).toEqual(postSuid);
  });
  it('Validating post for mockReplier in UserPostsModel', async () => {
    const replier = await UserModel().getUserByEmail(mockReplier.email);
    const result = await mongo.collection(USERPOSTS).find(
      { userId: replier._id, threadSuid },
    ).toArray();
    expect(result.length).toEqual(1);
    expect(result[0].posts[0].suid).toEqual(replySuid);
  });
  it('Validating post for mockUser in NotificationModel', async () => {
    const mockUserDoc = await UserModel().getUserByEmail(mockUser.email);
    const mockReplierDoc = await UserModel().getUserByEmail(mockReplier.email);
    const replyAuthor = await UserModel().methods(mockReplierDoc).author(threadSuid, true);

    const repResult = await NotificationModel().getNotiSlice(mockUserDoc, 'replied', { after: '', limit: 10 });
    const { replied } = repResult;
    expect(replied[0].id).toEqual(`replied:${threadId}`);
    expect(replied[0].sendTo).toEqual(mockUserDoc._id);
    expect(replied[0].type).toEqual('replied');

    const quotedResult = await NotificationModel().getNotiSlice(mockUserDoc, 'quoted', { after: '', limit: 10 });
    const { quoted } = quotedResult;
    expect(quoted[0].id).toEqual(`quoted:${replyId}:${postId}`);
    expect(quoted[0].sendTo).toEqual(mockUserDoc._id);
    expect(quoted[0].type).toEqual('quoted');
    expect(quoted[0].quoted.threadId).toEqual(threadId);
    expect(quoted[0].quoted.postId).toEqual(replyId);
    expect(quoted[0].quoted.quotedPostId).toEqual(postId);
    expect(quoted[0].quoted.quoter).toEqual(replyAuthor);
    expect(quoted[0].quoted.quoterId).toEqual(mockReplierDoc._id);
  });
});
