import gql from 'graphql-tag';
import dbClient, { db } from '~/dbClient';

import { startRepl } from '../__utils__/mongoServer';
import { mockUser, mutate } from '../__utils__/apolloClient';

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
  it('Create collections', async () => {
    await db.createCollection(USERPOSTS);
  });
  it('Setting mainTag', async () => {
    await TagModel().addMainTag('MainA');
  });
  it('Posting thread', async () => {
    const user = await UserModel().getUserByEmail(mockEmail);
    const { data, errors } = await mutate({
      mutation: PUB_THREAD, variables: { thread: mockThread },
    });
    expect(errors).toBeUndefined();

    const result = data.pubThread;
    const author = await UserModel().methods(user).author(Uid.encode(result.id), true);
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
    const user = await UserModel().getUserByEmail(mockEmail);
    const result = await dbClient.collection(USERPOSTS).find(
      { userId: user._id, threadId },
    ).toArray();
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating thread in TagModel', async () => {
    const result = await TagModel().getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockTagTree));
  });
});

describe('Testing posting a thread without subTags', () => {
  it('Creating collection', async () => {
    await db.dropDatabase();
    await db.createCollection(USERPOSTS);
  });
  it('Setting mainTag', async () => {
    await TagModel().addMainTag('MainA');
  });
  it('Posting thread', async () => {
    const user = await UserModel().getUserByEmail(mockEmail);
    const { data, errors } = await mutate({
      mutation: PUB_THREAD, variables: { thread: mockThreadNoSubTags },
    });
    expect(errors).toBeUndefined();

    const result = data.pubThread;
    const author = await UserModel().methods(user).author(Uid.encode(result.id), true);
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
    const user = await UserModel().getUserByEmail(mockEmail);
    const result = await dbClient.collection(USERPOSTS).find(
      { userId: user._id, threadId },
    ).toArray();
    expect(result.length).toEqual(1);
    expect(result[0].posts.length).toEqual(0);
  });
  it('Validating thread in TagModel', async () => {
    const result = await TagModel().getTree();
    expect(JSON.stringify(result[0])).toEqual(JSON.stringify(mockAltTagTree));
  });
});
