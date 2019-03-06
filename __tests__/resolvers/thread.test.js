import gql from 'graphql-tag';
import mongoose from 'mongoose';

import { startRepl } from '../__utils__/mongoServer';
import { mockUser, mutate } from '../__utils__/apolloClient';

import Uid from '~/uid';
import UserModel, { UserPostsModel } from '~/models/user';
import TagModel from '~/models/tag';
import ThreadModel from '~/models/thread';

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

// const mockTags = { mainTag: 'MainA', subTags: ['SubA', 'SubB'] };
const mockEmail = mockUser.email;
const mockTagTree = {
  mainTag: 'MainA', subTags: ['SubA'],
};
const mockAltTagTree = {
  mainTag: 'MainA', subTags: [],
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
  title: 'TestTitle',
};
let threadId;

// const ADD_TAGS = gql`
//   mutation AddTags($tags: [String!]) {
//     editConfig(config: {
//       mainTags: $tags,
//     }) { mainTags }
//   }
// `;
//
// const GET_TAGS = gql`
//   query {
//     tags {
//       mainTags,
//       tree {
//         mainTag, subTags
//       }
//     }
//   }
// `;

const PUB_THREAD = gql`
  mutation PubThread($thread: ThreadInput!) {
    pubThread(thread: $thread) {
      id, anonymous, author, content, createdAt, mainTag, subTags, title, replyCount, locked, blocked
    }
  }
`;

describe('Testing posting a thread', () => {
  it('Setting mainTag', async () => {
    await TagModel.addMainTag('MainA');
  });
  it('Posting thread', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const { data, errors } = await mutate({
      mutation: PUB_THREAD, variables: { thread: mockThread },
    });
    expect(errors).toBeUndefined();

    const result = data.pubThread;
    const author = await user.author(Uid.encode(result.id), true);
    expect(JSON.stringify(result.subTags)).toEqual(JSON.stringify(mockThread.subTags));
    expect(result.anonymous).toEqual(true);
    expect(result.author).toEqual(author);
    expect(result.content).toEqual(mockThread.content);
    expect(result.title).toEqual(mockThread.title);
    expect(result.blocked).toEqual(false);
    expect(result.locked).toEqual(false);

    threadId = result._id;
  });
  it('Validating thread in UserPostsModel', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const result = await UserPostsModel.find({ userId: user._id, threadId });
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating thread in TagModel', async () => {
    const result = await TagModel.getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockTagTree));
  });
});

describe('Testing posting a thread without subTags', () => {
  it('Creating collection', async () => {
    await mongoose.connection.db.dropDatabase();
    await TagModel.createCollection();
    await ThreadModel.createCollection();
    await UserPostsModel.createCollection();
    await TagModel.addMainTag('MainA');
  });
  it('Posting thread', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const { data, errors } = await mutate({
      mutation: PUB_THREAD, variables: { thread: mockThreadNoSubTags },
    });
    expect(errors).toBeUndefined();

    const result = data.pubThread;
    const author = await user.author(Uid.encode(result.id), true);
    expect(JSON.stringify(result.subTags)).toEqual(JSON.stringify([]));
    expect(result.anonymous).toEqual(true);
    expect(result.author).toEqual(author);
    expect(result.content).toEqual(mockThread.content);
    expect(result.title).toEqual(mockThread.title);
    expect(result.blocked).toEqual(false);
    expect(result.locked).toEqual(false);

    threadId = result._id;
  });
  it('Validating thread in UserPostsModel', async () => {
    const user = await UserModel.getUserByEmail(mockEmail);
    const result = await UserPostsModel.find({ userId: user._id, threadId });
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating thread in TagModel', async () => {
    const result = await TagModel.getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockAltTagTree));
  });
});
