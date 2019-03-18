import { startRepl } from '../__utils__/mongoServer';
import { ObjectId } from 'bson-ext';

import Uid from '~/uid';
import NotificationModel from '~/models/notification';
import UserModel from '~/models/user';

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

const mockEmail = 'test@example.com';
const mockReplierEmail = 'replier@example.com';
const threadId = ObjectId();
const post1Id = ObjectId();
const post2Id = ObjectId();
let threadSuid;
let post1Suid;
let post2Suid;
const threadTime = new Date();
const post1Time = new Date();
post1Time.setMinutes(post1Time.getMinutes() + 10);
const post2Time = new Date();
post2Time.setMinutes(post2Time.getMinutes() + 10);

describe('Testing notification', () => {
  it('AIO test suite', async () => {
    threadSuid = await Uid.newSuid();
    post1Suid = await Uid.newSuid();
    post2Suid = await Uid.newSuid();

    const mockUser = await UserModel().getUserByEmail(mockEmail);
    const mockReplier = await UserModel().getUserByEmail(mockReplierEmail);
    const mockThread = {
      _id: threadId,
      suid: threadSuid,
      uid: () => (Uid.decode(threadSuid)),
      anonymous: true,
      userId: mockUser._id,
      author: await UserModel().methods(mockUser).author(threadSuid, true),
      tags: ['MainA'],
      locked: false,
      blocked: false,
      createdAt: threadTime,
      updatedAt: threadTime,
    };

    const mockPost1 = {
      _id: post1Id,
      suid: post1Suid,
      uid: () => (Uid.decode(post1Suid)),
      userId: mockUser._id,
      threadId,
      anonymous: true,
      author: await UserModel().methods(mockUser).author(threadSuid, true),
      createdAt: post1Time,
      updatedAt: post1Time,
      blocked: false,
      quoteIds: [],
      content: 'Post1',
    };

    const mockPost2 = {
      _id: post2Id,
      suid: post2Suid,
      uid: () => (Uid.decode(post2Suid)),
      userId: mockReplier._id,
      threadId,
      anonymous: true,
      author: await UserModel().methods(mockReplier).author(threadSuid, true),
      createdAt: post2Time,
      updatedAt: post2Time,
      blocked: false,
      quoteIds: [post1Id],
      content: 'Post2',
    };

    await NotificationModel().sendQuotedNoti(mockPost2, mockThread, [mockPost1]);
    await NotificationModel().sendRepliedNoti(mockPost2, mockThread);

    const quotedResult = await NotificationModel().getNotiSlice(mockUser, 'quoted', { after: '', limit: 10 });
    const { quoted, sliceInfo: quotedSliceInfo } = quotedResult;
    expect(quoted[0].id).toEqual(`quoted:${mockPost2.uid()}:${mockPost1.uid()}`);
    expect(quoted[0].eventTime).toEqual(post2Time);
    expect(quoted[0].sendTo).toEqual(mockUser._id);
    expect(quoted[0].type).toEqual('quoted');
    expect(quoted[0].quoted.threadId).toEqual(mockThread.uid());
    expect(quoted[0].quoted.postId).toEqual(mockPost2.uid());
    expect(quoted[0].quoted.quotedPostId).toEqual(mockPost1.uid());
    expect(quoted[0].quoted.quoter).toEqual(mockPost2.author);
    expect(quoted[0].quoted.quoterId).toEqual(mockReplier._id);
    expect(quotedSliceInfo.firstCursor).toEqual(quoted[0]._id.toHexString());
    expect(quotedSliceInfo.lastCursor).toEqual(quoted[0]._id.toHexString());
    expect(quotedSliceInfo.hasNext).toBe(false);


    const repResult = await NotificationModel().getNotiSlice(mockUser, 'replied', { after: '', limit: 10 });
    const { replied, sliceInfo: repSliceInfo } = repResult;
    expect(replied[0].id).toEqual(`replied:${mockThread.uid()}`);
    expect(replied[0].eventTime).toEqual(post2Time);
    expect(replied[0].sendTo).toEqual(mockUser._id);
    expect(replied[0].type).toEqual('replied');
    expect(repSliceInfo.firstCursor).toEqual(replied[0]._id.toHexString());
    expect(repSliceInfo.lastCursor).toEqual(replied[0]._id.toHexString());
    expect(repSliceInfo.hasNext).toBe(false);
  });
});
