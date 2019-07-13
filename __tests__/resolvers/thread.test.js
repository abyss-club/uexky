import gql from 'graphql-tag';

import { mockUser, mutate } from '../__utils__/apolloClient';

import AidModel from '~/models/aid';
import UserModel from '~/models/user';
import { query as pgq } from '~/utils/pg';

import startPg, { migrate } from '../__utils__/pgServer';
import mockContext from '../__utils__/context';

let pgPool;

beforeAll(async () => {
  await migrate();
  pgPool = await startPg();
});

afterAll(async () => {
  await pgPool.query('DROP SCHEMA public CASCADE; CREATE SCHEMA public;');
  pgPool.end();
});

const mockEmail = mockUser.email;
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
  it('parpare data', async () => {
    await mockContext({ email: mockUser.email, name: mockUser.name });
    await pgq('INSERT INTO tag (name, is_main) VALUES ($1, $2)', ['MainA', true]);
  });
  it('Posting thread', async () => {
    const user = await UserModel.findByEmail({ email: mockEmail });
    const { data, errors } = await mutate({
      mutation: PUB_THREAD, variables: { thread: mockThread },
    });
    expect(errors).toBeUndefined();

    const result = data.pubThread;
    const aid = await AidModel.getAid({ userId: user.id, threadId: result.id });
    expect(result.mainTag).toEqual('MainA');
    expect(result.subTags).toEqual(['SubA']);
    expect(result.anonymous).toEqual(true);
    expect(result.author).toEqual(aid.duid);
    expect(result.content).toEqual(mockThread.content);
    expect(result.title).toEqual(mockThread.title);
    expect(result.blocked).toEqual(false);
    expect(result.locked).toEqual(false);
  });
});

describe('Testing posting a thread without subTags', () => {
  it('Posting thread', async () => {
    const user = await UserModel.findByEmail({ email: mockEmail });
    const { data, errors } = await mutate({
      mutation: PUB_THREAD, variables: { thread: mockThreadNoSubTags },
    });
    expect(errors).toBeUndefined();

    const result = data.pubThread;
    const aid = await AidModel.getAid({ userId: user.id, threadId: result.id });
    expect(JSON.stringify(result.subTags)).toEqual(JSON.stringify([]));
    expect(result.anonymous).toEqual(true);
    expect(result.author).toEqual(aid.duid);
    expect(result.content).toEqual(mockThread.content);
    expect(result.title).toEqual(mockThread.title);
    expect(result.blocked).toEqual(false);
    expect(result.locked).toEqual(false);
  });
});
