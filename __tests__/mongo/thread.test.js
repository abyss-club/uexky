import mongoose from 'mongoose';
import sleep from 'sleep-promise';
import { MongoMemoryReplSet } from 'mongodb-memory-server';
import UserModel, { UserPostsModel } from '~/models/user';
import ThreadModel from '~/models/thread';
import TagModel from '~/models/tag';
import ConfigModel from '~/models/config';
import Uid from '~/uid';

// May require additional time for downloading MongoDB binaries
// jasmine.DEFAULT_TIMEOUT_INTERVAL = 600000;

let replSet;
const opts = { useNewUrlParser: true };

beforeAll(async () => {
  replSet = new MongoMemoryReplSet({
    instanceOpts: [
      { storageEngine: 'wiredTiger' },
      // { storageEngine: 'wiredTiger' },
      // { storageEngine: 'wiredTiger' },
    ],
  });
  await replSet.waitUntilRunning();
  // const uri = await replSet.getUri();
  // console.log({ uri });
  // mongoServer = new MongoMemoryServer();
  // const mongoUri = await mongoServer.getConnectionString();
  const mongoUri = `${await replSet.getConnectionString()}?replicaSet=testset`;
  // console.log(mongoUri);
  // const dbName = await replSet.getDbName();

  await sleep(2000);

  await mongoose.connect(mongoUri, opts, (err) => {
    if (err) console.error(err);
  });
});

afterAll(() => {
  mongoose.disconnect();
  replSet.stop();
  // mongoServer.stop();
});

const mockEmail = 'test@example.com';
const mockThread = {
  anonymous: true,
  content: 'Test Content',
  mainTag: 'MainA',
  subTags: ['SubA'],
  title: 'TestTitle',
};
const mockThreadNoSubTags = {
  anonymous: true,
  content: 'Test Content',
  mainTag: 'MainA',
  subTags: [],
  title: 'TestTitle',
};
const mockTagTree = {
  mainTag: 'MainA', subTags: ['SubA'],
};
const mockAltTagTree = {
  mainTag: 'MainA', subTags: [],
};
let threadSuid;

describe('Testing posting a thread', () => {
  it('Preparations', async () => {
    // threadSuid = await Uid.newSuid();
    await ConfigModel.modifyMainTags(['MainA']);
  });
  it('Posting thread', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const newThread = mockThread;
    const { id } = await ThreadModel.pubThread({ user }, newThread);
    const threadResult = await ThreadModel.findByUid(id);

    threadSuid = threadResult.suid;

    const author = await user.author(threadSuid, true);
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
  });
  it('Validating thread in UserPostsModel', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const result = await UserPostsModel.find({ userId: user._id, threadSuid });
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating thread in TagModel', async () => {
    const result = await TagModel.getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockTagTree));
  });
});
//
describe('Testing posting a thread without subTags', () => {
  it('Preparations', async () => {
    await mongoose.connection.db.dropDatabase();
    await UserPostsModel.createCollection(); // necessary after db dropped
    await ConfigModel.modifyMainTags(['MainA']);
  });
  it('Posting thread', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const newThread = mockThreadNoSubTags;
    const { id } = await ThreadModel.pubThread({ user }, newThread);
    const threadResult = await ThreadModel.findByUid(id);

    threadSuid = threadResult.suid;

    const author = await user.author(threadResult.suid, true);
    expect(JSON.stringify(threadResult.subTags))
      .toEqual(JSON.stringify(mockThreadNoSubTags.subTags));
    expect(JSON.stringify(threadResult.tags)).toEqual(
      JSON.stringify([mockThreadNoSubTags.mainTag, ...mockThreadNoSubTags.subTags]),
    );
    expect(threadResult.anonymous).toEqual(true);
    expect(threadResult.content).toEqual(mockThreadNoSubTags.content);
    expect(threadResult.title).toEqual(mockThreadNoSubTags.title);
    expect(threadResult.userId).toEqual(user._id);
    expect(threadResult.author).toEqual(author);
    expect(threadResult.blocked).toEqual(false);
    expect(threadResult.locked).toEqual(false);
  });
  it('Validating thread in UserPostsModel', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const result = await UserPostsModel.find({ userId: user._id, threadSuid });
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating thread in TagModel', async () => {
    const result = await TagModel.getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockAltTagTree));
  });
});
