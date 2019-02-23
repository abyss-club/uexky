import mongoose from 'mongoose';
import { startRepl } from '../__utils__/mongoServer';

import Uid from '~/uid';
import UserModel, { UserPostsModel } from '~/models/user';
import ThreadModel from '~/models/thread';
import PostModel from '~/models/post';
import TagModel from '~/models/tag';
import ConfigModel from '~/models/config';
import NotificationModel from '~/models/notification';

// May require additional time for downloading MongoDB binaries
// Temporary hack for parallel tests
jasmine.DEFAULT_TIMEOUT_INTERVAL = 600000;

let replSet;
beforeAll(async () => {
  replSet = await startRepl();
});

afterAll(() => {
  mongoose.disconnect();
  replSet.stop();
});

const mockEmail = 'test@example.com';
const mockReplierEmail = 'replier@example.com';
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
  it('Setting tags', async () => {
    await ConfigModel.modifyMainTags(['MainA']);
  });
  it('Posting a thread', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const newThread = mockThread;
    const { _id } = await ThreadModel.pubThread({ user }, newThread);
    const threadResult = await ThreadModel.findOne({ _id });
    const author = await user.author(threadResult.suid, true);
    expect(JSON.stringify(threadResult.subTags)).toEqual(JSON.stringify(mockThread.subTags));
    expect(JSON.stringify(threadResult.tags))
      .toEqual(JSON.stringify([mockThread.mainTag, ...mockThread.subTags]));
    expect(threadResult.anonymous).toEqual(true);
    expect(threadResult.content).toEqual(mockThread.content);
    expect(threadResult.title).toEqual(mockThread.title);
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
    const user = await UserModel.getUserByEmail(mockEmail);
    const result = await UserPostsModel.find({ userId: user._id, threadSuid });
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating the thread in TagModel', async () => {
    const result = await TagModel.getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockTagTree));
  });
});

describe('Testing replying a thread', () => {
  it('Posting reply', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const newPost = mockPost;
    const { _id } = await PostModel.pubPost({ user }, newPost);
    const postResult = await PostModel.findOne({ _id });

    postSuid = postResult.suid;
    postId = Uid.decode(postSuid);

    const author = await user.author(threadSuid, true);
    expect(postResult.threadSuid).toEqual(threadSuid);
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toBeUndefined();
    expect(postResult.userId).toEqual(user._id);
    expect(postResult.author).toEqual(author);
    expect(postResult.blocked).toEqual(false);

    mockReply.quoteIds.push(postId);
  });
  it('Posting reply with quotes', async () => {
    const user = await UserModel.getUserByEmail(mockReplierEmail);
    const newPost = mockReply;
    const { _id } = await PostModel.pubPost({ user }, newPost);
    const postResult = await PostModel.findOne({ _id });

    replySuid = postResult.suid;
    replyId = Uid.decode(replySuid);

    const author = await user.author(threadSuid, true);
    expect(postResult.threadSuid).toEqual(threadSuid);
    expect(postResult.anonymous).toEqual(true);
    expect(postResult.content).toEqual(mockPost.content);
    expect(postResult.quotes).toBeUndefined();
    expect(postResult.userId).toEqual(user._id);
    expect(postResult.author).toEqual(author);
    expect(postResult.blocked).toEqual(false);
  });
  it('Validate updated post', async () => {
    const postResult = await PostModel.findOne({ suid: replySuid });
    expect(postResult.quoteSuids[0]).toEqual(postSuid);
  });
  it('Validating post for mockUser in UserPostsModel', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const result = await UserPostsModel.find({ userId: user._id, threadSuid });
    expect(result.length).toEqual(1);
    expect(result[0].posts[0].suid).toEqual(postSuid);
  });
  it('Validating post for mockReplier in UserPostsModel', async () => {
    const replier = await UserModel.getUserByEmail(mockReplierEmail);
    const result = await UserPostsModel.find({ userId: replier._id, threadSuid });
    expect(result.length).toEqual(1);
    expect(result[0].posts[0].suid).toEqual(replySuid);
  });
  it('Validating post for mockUser in NotificationModel', async () => {
    const mockUser = await UserModel.getUserByEmail(mockEmail);
    const mockReplier = await UserModel.getUserByEmail(mockReplierEmail);
    const replyAuthor = await mockReplier.author(threadSuid, true);

    const repResult = await NotificationModel.getNotiSlice(mockUser, 'replied', { after: '', limit: 10 });
    const { replied } = repResult;
    expect(replied[0].id).toEqual(`replied:${threadId}`);
    expect(replied[0].sendTo).toEqual(mockUser._id);
    expect(replied[0].type).toEqual('replied');

    const quotedResult = await NotificationModel.getNotiSlice(mockUser, 'quoted', { after: '', limit: 10 });
    const { quoted } = quotedResult;
    expect(quoted[0].id).toEqual(`quoted:${replyId}:${postId}`);
    expect(quoted[0].sendTo).toEqual(mockUser._id);
    expect(quoted[0].type).toEqual('quoted');
    expect(quoted[0].quoted.threadId).toEqual(threadId);
    expect(quoted[0].quoted.postId).toEqual(replyId);
    expect(quoted[0].quoted.quotedPostId).toEqual(postId);
    expect(quoted[0].quoted.quoter).toEqual(replyAuthor);
    expect(quoted[0].quoted.quoterId).toEqual(mockReplier._id);
  });
});
