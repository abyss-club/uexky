import mongoose from 'mongoose';
import MongoMemoryServer from 'mongodb-memory-server';
import Uid from '~/uid';
import NotificationModel from '~/models/notification';
import UserModel from '~/models/user';

// May require additional time for downloading MongoDB binaries
// jasmine.DEFAULT_TIMEOUT_INTERVAL = 600000;

let mongoServer;
const opts = { useNewUrlParser: true };

beforeAll(async () => {
  mongoServer = new MongoMemoryServer();
  const mongoUri = await mongoServer.getConnectionString();
  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) console.error(err);
  });
});

afterAll(() => {
  mongoose.disconnect();
  mongoServer.stop();
});

const mockEmail = 'test@example.com';
const mockReplierEmail = 'replier@example.com';
const threadId = mongoose.Types.ObjectId();
const post1Id = mongoose.Types.ObjectId();
const post2Id = mongoose.Types.ObjectId();
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

    const mockUser = await UserModel.getUserByEmail(mockEmail);
    const mockReplier = await UserModel.getUserByEmail(mockReplierEmail);
    const mockThread = {
      _id: threadId,
      suid: threadSuid,
      uid: () => (Uid.decode(threadSuid)),
      anonymous: true,
      userId: mockUser._id,
      tags: ['MainA'],
      locked: false,
      blocked: false,
      createdAt: threadTime,
      updatedAt: threadTime,
    };
    mockThread.author = await mockUser.author(mockThread.suid, mockThread.anonymous);

    const mockPost1 = {
      _id: post1Id,
      suid: post1Suid,
      uid: () => (Uid.decode(post1Suid)),
      userId: mockUser._id,
      threadId,
      anonymous: true,
      author: await mockUser.author(threadId, true),
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
      author: await mockReplier.author(threadId, true),
      createdAt: post2Time,
      updatedAt: post2Time,
      blocked: false,
      quoteIds: [post1Id],
      content: 'Post2',
    };

    await NotificationModel.sendQuotedNoti(mockPost2, mockThread, [mockPost1]);
    await NotificationModel.sendRepliedNoti(mockPost2, mockThread);

    const quotedResult = await NotificationModel.getNotiSlice(mockUser, 'quoted', { after: '', limit: 10 });
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
    expect(quotedSliceInfo.firstCursor).toEqual(quoted[0]._id.valueOf());
    expect(quotedSliceInfo.lastCursor).toEqual(quoted[0]._id.valueOf());
    expect(quotedSliceInfo.hasNext).toBe(false);


    const repResult = await NotificationModel.getNotiSlice(mockUser, 'replied', { after: '', limit: 10 });
    const { replied, sliceInfo: repSliceInfo } = repResult;
    expect(replied[0].id).toEqual(`replied:${mockThread.uid()}`);
    expect(replied[0].eventTime).toEqual(post2Time);
    expect(replied[0].sendTo).toEqual(mockUser._id);
    expect(replied[0].type).toEqual('replied');
    expect(repSliceInfo.firstCursor).toEqual(replied[0]._id.valueOf());
    expect(repSliceInfo.lastCursor).toEqual(replied[0]._id.valueOf());
    expect(repSliceInfo.hasNext).toBe(false);
  });
});
