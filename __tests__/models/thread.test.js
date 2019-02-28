import mongoose from 'mongoose';
import { startRepl } from '../__utils__/mongoServer';

import UserModel, { UserPostsModel } from '~/models/user';
import ThreadModel from '~/models/thread';
import TagModel from '~/models/tag';

// May require additional time for downloading MongoDB binaries
// Temporary hack for parallel tests
jest.setTimeout(600000);

let replSet;
beforeAll(async () => {
  replSet = await startRepl();
});

afterAll(() => {
  mongoose.disconnect();
  replSet.stop();
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
  it('Setting tags', async () => {
    await TagModel.addMainTag('MainA');
  });
  it('Posting a thread', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const newThread = mockThread;
    const { _id } = await ThreadModel.pubThread({ user }, newThread);
    const threadResult = await ThreadModel.findOne({ _id }).exec();
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
  });
  it('Validating the thread in UserPostsModel', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const result = await UserPostsModel.find({ userId: user._id, threadSuid }).exec();
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating the thread in TagModel', async () => {
    const result = await TagModel.getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockTagTree));
  });
});

describe('Testing posting a thread without subTags', () => {
  it('Creating collection', async () => {
    await mongoose.connection.db.dropDatabase();
    await UserPostsModel.createCollection(); // necessary after db dropped
    await ThreadModel.createCollection(); // necessary after db dropped
    await TagModel.addMainTag('MainA');
  });
  it('Posting a thread', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const newThread = mockThreadNoSubTags;
    const { _id } = await ThreadModel.pubThread({ user }, newThread);
    const threadResult = await ThreadModel.findOne({ _id }).exec();
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

    threadSuid = threadResult.suid;
  });
  it('Validating the thread in UserPostsModel', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const result = await UserPostsModel.find({ userId: user._id, threadSuid }).exec();
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating the thread in TagModel', async () => {
    const result = await TagModel.getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockAltTagTree));
  });
});
