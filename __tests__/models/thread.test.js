import startRepl from '../__utils__/mongoServer';
import mongo, { db } from '~/utils/mongo';

import UserModel from '~/models/user';
// import UserPostsModel from '~/models/userPosts';
import ThreadModel from '~/models/thread';
import TagModel from '~/models/tag';

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
const USERPOSTS = 'userPosts';
const col = () => mongo.collection(THREAD);

const mockUser = {
  email: 'test@example.com',
  name: 'testUser',
};
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
  it('Create collections', async () => {
    await db.createCollection(USERPOSTS);
  });
  it('Setting tags', async () => {
    await TagModel().addMainTag('MainA');
  });
  it('Posting a thread', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    const newThread = mockThread;
    const { _id } = await ThreadModel({ user }).pubThread(newThread);
    const threadResult = await col().findOne({ _id });
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

describe('Testing posting a thread without subTags', () => {
  it('Creating collection', async () => {
    await db.dropDatabase();
    await db.createCollection(USERPOSTS); // create collection after db droppped
  });
  it('Setting tags', async () => {
    await TagModel().addMainTag('MainA');
  });
  it('Posting a thread', async () => {
    const user = await UserModel().getUserByEmail(mockUser.email);
    const newThread = mockThreadNoSubTags;
    const { _id } = await ThreadModel({ user }).pubThread(newThread);
    const threadResult = await col().findOne({ _id });
    const author = await UserModel({ user }).methods(user).author(threadResult.suid, true);
    expect(JSON.stringify(threadResult.subTags))
      .toEqual(JSON.stringify(newThread.subTags));
    expect(JSON.stringify(threadResult.tags)).toEqual(
      JSON.stringify([newThread.mainTag, ...newThread.subTags]),
    );
    expect(threadResult.anonymous).toEqual(true);
    expect(threadResult.content).toEqual(newThread.content);
    expect(threadResult.title).toEqual(newThread.title);
    expect(threadResult.userId).toEqual(user._id);
    expect(threadResult.author).toEqual(author);
    expect(threadResult.blocked).toEqual(false);
    expect(threadResult.locked).toEqual(false);

    threadSuid = threadResult.suid;
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
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockAltTagTree));
  });
});
